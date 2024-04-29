package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/Kotaro7750/rakuten-security-exporter/proto"
	"github.com/bojanz/currency"
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

type Asset []*IndividualAsset

func constructAsset(assets *proto.TotalAssetResponse) (Asset, error) {
	constructedAssets := make(Asset, 0)

	for _, asset := range assets.GetAsset() {
		convertedAsset, err := NewIndividualAsset(asset)
		if err != nil {
			return nil, err
		}

		constructedAssets = append(constructedAssets, convertedAsset)
	}

	return constructedAssets, nil
}

func (asset *Asset) Summarize() (*AssetSummary, error) {
	var totalAcquisitionPrice currency.Amount
	var totalPrice currency.Amount

	for _, individualAsset := range *asset {
		var err error

		acquisitionPrice, err := individualAsset.averageAcquisitionPrice.Mul(strconv.FormatFloat(individualAsset.count, 'f', -1, 64))
		totalPrice, err = totalPrice.Add(individualAsset.currentPrice)
		if err != nil {
			return nil, err
		}

		totalAcquisitionPrice, err = totalAcquisitionPrice.Add(acquisitionPrice)
		if err != nil {
			return nil, err
		}

	}

	return &AssetSummary{totalAcquisitionPrice, totalPrice}, nil
}

type AssetSummary struct {
	totalAcquisitionPrice currency.Amount
	totalPrice            currency.Amount
}

func (assetSummary *AssetSummary) PerformanceExcludingCurrencyImpact() (float64, error) {
	return amountRatio(assetSummary.totalPrice, assetSummary.totalAcquisitionPrice)
}

type IndividualAsset struct {
	assetType               string
	ticker                  string
	name                    string
	account                 string
	count                   float64
	averageAcquisitionPrice currency.Amount
	currentUnitPrice        currency.Amount
	currentPrice            currency.Amount
}

func NewIndividualAsset(asset *proto.Asset) (*IndividualAsset, error) {
	averageAcquisitionPrice, err := currency.NewAmount(strconv.FormatFloat(asset.GetAverageAcquisitionPrice(), 'f', -1, 64), "USD")
	if err != nil {
		return nil, err
	}

	currentUnitPrice, err := currency.NewAmount(strconv.FormatFloat(asset.GetCurrentUnitPrice(), 'f', -1, 64), "USD")
	if err != nil {
		return nil, err
	}

	currentPrice, err := currency.NewAmount(strconv.FormatFloat(asset.GetCurrentPrice(), 'f', -1, 64), "USD")
	if err != nil {
		return nil, err
	}

	return &IndividualAsset{
		assetType:               asset.GetType(),
		ticker:                  asset.GetTicker(),
		name:                    asset.GetName(),
		account:                 asset.GetAccount(),
		count:                   asset.GetCount(),
		averageAcquisitionPrice: averageAcquisitionPrice,
		currentUnitPrice:        currentUnitPrice,
		currentPrice:            currentPrice,
	}, nil
}

func (individualAsset *IndividualAsset) PerformanceExcludingCurrencyImpact() (float64, error) {
	return amountRatio(individualAsset.currentUnitPrice, individualAsset.averageAcquisitionPrice)
}

type WithdrawalSummary struct {
	TotalInvestmentAmount currency.Amount
}

func constructWithdrawalStatistics(withdrawalHistories *proto.ListWithdrawalHistoriesResponse) WithdrawalSummary {
	var total float64 = 0

	for _, withdrawalHistory := range withdrawalHistories.GetHistory() {
		withdrawalType := withdrawalHistory.GetType()
		amount := float64(withdrawalHistory.GetAmount())

		if withdrawalType == "in" {
			total += amount
		} else if withdrawalType == "out" {
			total -= amount
		}
	}

	totalAmount, err := currency.NewAmount(strconv.FormatFloat(total, 'f', -1, 64), "JPY")
	if err != nil {
		log.Fatalf("error %f", err)
	}

	return WithdrawalSummary{
		TotalInvestmentAmount: totalAmount,
	}
}

func amountRatio(dividend currency.Amount, divisor currency.Amount) (float64, error) {
	dividendFloat, err := strconv.ParseFloat(dividend.Number(), 64)
	if err != nil {
		return 1, err
	}

	divisorFloat, err := strconv.ParseFloat(divisor.Number(), 64)
	if err != nil {
		return 1, err
	}

	if divisorFloat == 0 {
		return 1, nil
	}

	return dividendFloat / divisorFloat, nil
}
