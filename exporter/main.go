package main

import (
	"context"
	"log"
	"time"

	"github.com/Kotaro7750/rakuten-security-exporter/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var opt []grpc.DialOption = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial("localhost:50051", opt...)
	if err != nil {
		log.Fatalf("dial error %v", err)
	}
	defer conn.Close()

	client := proto.NewRakutenSecurityScraperClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	total_asset, err := client.TotalAssets(ctx, &proto.TotalAssetRequest{})
	if err != nil {
		log.Fatalf("error %v", err)
	}

	withdrawal_history, err := client.ListWithdrawalHistories(ctx, &proto.ListWithdrawalHistoriesRequest{})
	if err != nil {
		log.Fatalf("error %v", err)
	}

	dividend_history, err := client.ListDividendHistories(ctx, &proto.ListDividendHistoriesRequest{})
	if err != nil {
		log.Fatalf("error %v", err)
	}

	investmentReport, err := ConstructInvestmentReport(total_asset, withdrawal_history, dividend_history)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	performance, err := investmentReport.ConstructPerformanceReport(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local), "USD")
	if err != nil {
		log.Fatalf("error %v", err)
	}

	log.Printf("performance %v", performance)

	stat, err := investmentReport.dividendHistory.constructDividendStatistics("USD", &investmentReport.rateManager)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	for year, stat := range stat.total {
		total, err := stat.dividend.calcTotal()
		if err != nil {
			log.Fatalf("error %v", err)
		}

		log.Printf("total dividend: %d %s %f", year, total.String(), stat.totalGrowth)
	}

	for _, s := range stat.security {
		for year, a := range s.statistics {
			total, err := a.dividend.total.calcTotal()
			if err != nil {
				log.Fatalf("error %v", err)
			}

			unitPrice, err := a.dividend.unitPrice.calcTotal()
			if err != nil {
				log.Fatalf("error %v", err)
			}

			log.Printf("%s %d total %s %f unit %s %f", s.ticker, year, total.String(), a.totalGrowth, unitPrice.String(), a.unitPriceGrowth)
		}
	}
}

type InvestmentReport struct {
	asset                       Asset
	depositAndWithdrawalHistory DepositWithdrawalHistory
	dividendHistory             DividendHistory
	rateManager                 RateManager
}

type PerformanceReport struct {
	PerformanceExcludingCurrencyImpact Return
	TotalReturn                        Return
}

type Return struct {
	total  float64
	annual float64
}

func NewReturn(totalReturn float64, startDate time.Time) Return {
	return Return{
		total:  totalReturn,
		annual: calcAnnualReturn(totalReturn, startDate, time.Now()),
	}
}

func ConstructInvestmentReport(
	assetResponse *proto.TotalAssetResponse,
	withdrawalHistoryResponse *proto.ListWithdrawalHistoriesResponse,
	dividendHistoryResponse *proto.ListDividendHistoriesResponse,
) (InvestmentReport, error) {

	asset, err := constructAsset(assetResponse)
	if err != nil {
		return InvestmentReport{}, err
	}

	depositAndWithdrawalHistory, err := constructDepositWithdrawalHistory(withdrawalHistoryResponse)
	if err != nil {
		return InvestmentReport{}, err
	}

	dividendHistory, err := ConstructDividendHistory(dividendHistoryResponse)
	if err != nil {
		return InvestmentReport{}, err
	}

	rateManager := NewRateManager()

	err = rateManager.RegisterRate("USD", "JPY", 156.45)
	if err != nil {
		return InvestmentReport{}, err
	}

	return InvestmentReport{asset, depositAndWithdrawalHistory, dividendHistory, rateManager}, nil
}

func (ir *InvestmentReport) ConstructPerformanceReport(startDate time.Time, targetCurrencyCode string) (PerformanceReport, error) {
	assetSummary, err := ir.asset.Summarize(targetCurrencyCode, &ir.rateManager)
	if err != nil {
		return PerformanceReport{}, err
	}

	performanceExcludingCurrencyImpact, err := assetSummary.PerformanceExcludingCurrencyImpact(&ir.rateManager)
	if err != nil {
		return PerformanceReport{}, err
	}

	totalInvestmentAmount, err := ir.depositAndWithdrawalHistory.totalInvestmentAmount(targetCurrencyCode, &ir.rateManager)
	if err != nil {
		return PerformanceReport{}, err
	}

	totalReturn, err := amountRatio(assetSummary.totalPrice, totalInvestmentAmount, &ir.rateManager)
	if err != nil {
		return PerformanceReport{}, err
	}

	return PerformanceReport{
		PerformanceExcludingCurrencyImpact: NewReturn(performanceExcludingCurrencyImpact, startDate),
		TotalReturn:                        NewReturn(totalReturn, startDate),
	}, nil
}
