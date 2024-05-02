package main

import (
	"context"
	"log"
	"strconv"
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
	PerformanceExcludingCurrencyImpact Performance
	Performance                        Performance
}

type Performance struct {
	totalReturn  float64
	annualReturn float64
}

func NewPerformance(totalReturn float64, startDate time.Time) Performance {
	return Performance{
		totalReturn:  totalReturn,
		annualReturn: calcAnnualReturn(totalReturn, startDate, time.Now()),
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
	performance, err := ir.Performance(targetCurrencyCode)
	if err != nil {
		return PerformanceReport{}, err
	}

	assetSummary, err := ir.asset.Summarize(targetCurrencyCode, &ir.rateManager)
	if err != nil {
		return PerformanceReport{}, err
	}

	performanceExcludingCurrencyImpact, err := assetSummary.PerformanceExcludingCurrencyImpact(&ir.rateManager)
	if err != nil {
		return PerformanceReport{}, err
	}

	return PerformanceReport{
		PerformanceExcludingCurrencyImpact: NewPerformance(performanceExcludingCurrencyImpact, startDate),
		Performance:                        NewPerformance(performance, startDate),
	}, nil
}

func (ir *InvestmentReport) Performance(targetCurrencyCode string) (float64, error) {
	totalInvestmentAmount, err := ir.depositAndWithdrawalHistory.totalInvestmentAmount(targetCurrencyCode, &ir.rateManager)
	if err != nil {
		return 1, err
	}

	assetSummary, err := ir.asset.Summarize(targetCurrencyCode, &ir.rateManager)
	if err != nil {
		return 1, err
	}

	assetTotalPrice := assetSummary.totalPrice
	convertedAssetTotalPrice, err := ir.rateManager.Convert(assetTotalPrice, totalInvestmentAmount.CurrencyCode())
	if err != nil {
		return 1, err
	}

	totalInvestmentAmount.Number()
	totalInvestmentAmountFloat, err := strconv.ParseFloat(totalInvestmentAmount.Number(), 64)
	if err != nil {
		return 1, err
	}

	convertedAssetTotalPrice.Number()
	convertedAssetTotalPriceFloat, err := strconv.ParseFloat(convertedAssetTotalPrice.Number(), 64)
	if err != nil {
		return 1, err
	}

	log.Printf("%f %f", convertedAssetTotalPriceFloat, totalInvestmentAmountFloat)

	return convertedAssetTotalPriceFloat / totalInvestmentAmountFloat, nil
}
