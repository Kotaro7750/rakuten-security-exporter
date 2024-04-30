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

	performance, err := investmentReport.ConstructPerformanceReport(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local))
	if err != nil {
		log.Fatalf("error %v", err)
	}

	log.Printf("performance %v", performance)
}

type InvestmentReport struct {
	asset                Asset
	depositAndWithdrawal WithdrawalSummary
	dividendHistory      DividendHistory
	rateManager          RateManager
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

	depositAndWithdrawal := constructWithdrawalStatistics(withdrawalHistoryResponse)

	dividendHistory, err := ConstructDividendHistory(dividendHistoryResponse)
	if err != nil {
		return InvestmentReport{}, err
	}

	rateManager := NewRateManager()

	err = rateManager.RegisterRate("USD", "JPY", "156.45")
	if err != nil {
		return InvestmentReport{}, err
	}

	return InvestmentReport{asset, depositAndWithdrawal, dividendHistory, rateManager}, nil
}

func (ir *InvestmentReport) ConstructPerformanceReport(startDate time.Time) (PerformanceReport, error) {
	performance, err := ir.Performance()
	if err != nil {
		return PerformanceReport{}, err
	}

	assetSummary, err := ir.asset.Summarize()
	if err != nil {
		return PerformanceReport{}, err
	}

	performanceExcludingCurrencyImpact, err := assetSummary.PerformanceExcludingCurrencyImpact()
	if err != nil {
		return PerformanceReport{}, err
	}

	return PerformanceReport{
		PerformanceExcludingCurrencyImpact: NewPerformance(performanceExcludingCurrencyImpact, startDate),
		Performance:                        NewPerformance(performance, startDate),
	}, nil
}

func (ir *InvestmentReport) Performance() (float64, error) {
	totalInvestmentAmount := ir.depositAndWithdrawal.TotalInvestmentAmount

	assetSummary, err := ir.asset.Summarize()
	if err != nil {
		return 1, err
	}

	assetTotalPrice := assetSummary.totalPrice

	rate, err := ir.rateManager.GetRate(assetTotalPrice.CurrencyCode(), totalInvestmentAmount.CurrencyCode())
	if err != nil {
		return 1, err
	}

	convertedAssetTotalPrice, err := assetTotalPrice.Convert(totalInvestmentAmount.CurrencyCode(), rate)
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
