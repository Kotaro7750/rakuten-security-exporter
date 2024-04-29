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

	assets, err := constructAsset(total_asset)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	assetSummary, err := assets.Summarize()
	if err != nil {
		log.Fatalf("error %v", err)
	}

	performance, err := assetSummary.PerformanceExcludingCurrencyImpact()
	if err != nil {
		log.Fatalf("error %v", err)
	}

	log.Printf("response %v %v %v", assetSummary.totalPrice.Number(), assetSummary.totalAcquisitionPrice.Number(), performance)

	withdrawal_history, err := client.ListWithdrawalHistories(ctx, &proto.ListWithdrawalHistoriesRequest{})
	if err != nil {
		log.Fatalf("error %v", err)
	}

	log.Printf("withdrawalStat: %v", constructWithdrawalStatistics(withdrawal_history))

	dividend_history, err := client.ListDividendHistories(ctx, &proto.ListDividendHistoriesRequest{})
	if err != nil {
		log.Fatalf("error %v", err)
	}

	log.Printf("response %v", dividend_history)
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
