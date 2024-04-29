package main

import (
	"strconv"

	"github.com/Kotaro7750/rakuten-security-exporter/proto"
	"github.com/bojanz/currency"
)

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
