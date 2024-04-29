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

	investmentReport, err := ConstructInvestmentReport(total_asset, withdrawal_history)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	performance, err := investmentReport.Performance()
	if err != nil {
		log.Fatalf("error %v", err)
	}

	log.Printf("performance %v", performance)

	dividend_history, err := client.ListDividendHistories(ctx, &proto.ListDividendHistoriesRequest{})
	if err != nil {
		log.Fatalf("error %v", err)
	}

	log.Printf("response %v", dividend_history)
}

type InvestmentReport struct {
	asset                Asset
	depositAndWithdrawal WithdrawalSummary
	rateManager          RateManager
}

func ConstructInvestmentReport(assetResponse *proto.TotalAssetResponse, withdrawalHistoryResponse *proto.ListWithdrawalHistoriesResponse) (InvestmentReport, error) {
	asset, err := constructAsset(assetResponse)
	if err != nil {
		return InvestmentReport{}, err
	}

	depositAndWithdrawal := constructWithdrawalStatistics(withdrawalHistoryResponse)

	rateManager := NewRateManager()

	err = rateManager.RegisterRate("USD", "JPY", "156.45")
	if err != nil {
		return InvestmentReport{}, err
	}

	return InvestmentReport{asset, depositAndWithdrawal, rateManager}, nil
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
