package main

import (
	"fmt"
	"strconv"

	"github.com/Kotaro7750/rakuten-security-exporter/proto"
	"github.com/bojanz/currency"
)

type Security struct {
	assetType string
	ticker    string
	name      string
}

func newSecurity(assetType, ticker, name string) Security {
	return Security{assetType, ticker, name}
}

func (s *Security) identifier() string {
	// Name must not be used because it varies per file
	return fmt.Sprintf("%s %s", s.assetType, s.ticker)
}

type Asset []*IndividualAsset

func constructAsset(assets *proto.TotalAssetResponse) (Asset, RateManager, error) {
	constructedAssets := make(Asset, 0)

	for _, asset := range assets.GetAsset() {
		convertedAsset, err := NewIndividualAsset(asset)
		if err != nil {
			return nil, RateManager{}, err
		}

		constructedAssets = append(constructedAssets, convertedAsset)
	}

	rateManager := NewRateManager()

	for _, currencyRateToJPY := range assets.GetCurrencyRate() {
		err := rateManager.RegisterRate(currencyRateToJPY.GetCurrencyCode(), "JPY", currencyRateToJPY.GetRate())
		if err != nil {
			return nil, RateManager{}, err
		}
	}

	return constructedAssets, rateManager, nil
}

func (asset *Asset) Summarize(targetCurrencyCode string, rateManager *RateManager) (*AssetSummary, error) {
	var totalAcquisitionPrice currency.Amount
	var totalPrice currency.Amount

	for _, individualAsset := range *asset {
		var err error

		acquisitionPrice, err := individualAsset.averageAcquisitionPrice.Mul(strconv.FormatFloat(individualAsset.count, 'f', -1, 64))
		totalPrice, err = addAmount(totalPrice, individualAsset.currentPrice, targetCurrencyCode, rateManager)
		if err != nil {
			return nil, err
		}

		totalAcquisitionPrice, err = addAmount(totalAcquisitionPrice, acquisitionPrice, targetCurrencyCode, rateManager)
		if err != nil {
			return nil, err
		}
	}

	return &AssetSummary{totalAcquisitionPrice, totalPrice}, nil
}

func (asset *Asset) construceAssetCount() map[Security]float64 {
	assetCount := make(map[Security]float64, 0)

	for _, individualAsset := range *asset {
		assetCount[individualAsset.security] = individualAsset.count
	}

	return assetCount
}

type AssetSummary struct {
	totalAcquisitionPrice currency.Amount
	totalPrice            currency.Amount
}

func (assetSummary *AssetSummary) PerformanceExcludingCurrencyImpact(rateManager *RateManager) (float64, error) {
	return amountRatio(assetSummary.totalPrice, assetSummary.totalAcquisitionPrice, rateManager)
}

type IndividualAsset struct {
	security                Security
	account                 string
	count                   float64
	averageAcquisitionPrice currency.Amount
	currentUnitPrice        currency.Amount
	currentPrice            currency.Amount
}

func NewIndividualAsset(asset *proto.Asset) (*IndividualAsset, error) {
	averageAcquisitionPrice, err := currency.NewAmount(strconv.FormatFloat(asset.GetAverageAcquisitionPrice(), 'f', -1, 64), asset.GetCurrency())
	if err != nil {
		return nil, err
	}

	currentUnitPrice, err := currency.NewAmount(strconv.FormatFloat(asset.GetCurrentUnitPrice(), 'f', -1, 64), asset.GetCurrency())
	if err != nil {
		return nil, err
	}

	currentPrice, err := currency.NewAmount(strconv.FormatFloat(asset.GetCurrentPrice(), 'f', -1, 64), asset.GetCurrency())
	if err != nil {
		return nil, err
	}

	return &IndividualAsset{
		security:                newSecurity(asset.GetType(), asset.GetTicker(), asset.GetName()),
		account:                 asset.GetAccount(),
		count:                   asset.GetCount(),
		averageAcquisitionPrice: averageAcquisitionPrice,
		currentUnitPrice:        currentUnitPrice,
		currentPrice:            currentPrice,
	}, nil
}

func (individualAsset *IndividualAsset) PerformanceExcludingCurrencyImpact(rateManager *RateManager) (float64, error) {
	return amountRatio(individualAsset.currentUnitPrice, individualAsset.averageAcquisitionPrice, rateManager)
}
