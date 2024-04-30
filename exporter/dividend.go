package main

import (
	"strconv"
	"time"

	"github.com/Kotaro7750/rakuten-security-exporter/proto"
	"github.com/bojanz/currency"
)

type DividendHistory []Dividend

func ConstructDividendHistory(protoDividendHistoryResponse *proto.ListDividendHistoriesResponse) (DividendHistory, error) {
	dividendHistory := make(DividendHistory, 0)

	for _, protoDividendHistory := range protoDividendHistoryResponse.GetHistory() {
		dividend, err := ConstructDividend(protoDividendHistory)
		if err != nil {
			return DividendHistory{}, err
		}

		dividendHistory = append(dividendHistory, dividend)
	}

	return dividendHistory, nil
}

type Dividend struct {
	date             time.Time
	account          string
	assetType        string
	ticker           string
	name             string
	count            float64
	unitPrice        currency.Amount
	totalBeforeTaxed currency.Amount
	totalTaxes       currency.Amount
	total            currency.Amount
}

func ConstructDividend(dividendHistory *proto.DividendHistory) (Dividend, error) {
	date, err := time.Parse("2006/01/02", dividendHistory.GetDate())
	if err != nil {
		return Dividend{}, nil
	}

	dividendCurrency := dividendHistory.GetCurrency()

	unitPrice, err := currency.NewAmount(strconv.FormatFloat(dividendHistory.GetDividendUnitprice(), 'f', -1, 64), dividendCurrency)
	if err != nil {
		return Dividend{}, err
	}

	totalBeforeTaxed, err := currency.NewAmount(strconv.FormatFloat(dividendHistory.GetDividendTotalBeforeTaxes(), 'f', -1, 64), dividendCurrency)
	if err != nil {
		return Dividend{}, err
	}

	totalTaxes, err := currency.NewAmount(strconv.FormatFloat(dividendHistory.GetTotalTaxes(), 'f', -1, 64), dividendCurrency)
	if err != nil {
		return Dividend{}, err
	}

	total, err := currency.NewAmount(strconv.FormatFloat(dividendHistory.GetDividendTotal(), 'f', -1, 64), dividendCurrency)
	if err != nil {
		return Dividend{}, err
	}

	return Dividend{date: date,
		account:          dividendHistory.GetAccount(),
		assetType:        dividendHistory.GetType(),
		ticker:           dividendHistory.GetTicker(),
		name:             dividendHistory.GetName(),
		count:            dividendHistory.GetCount(),
		unitPrice:        unitPrice,
		totalBeforeTaxed: totalBeforeTaxed,
		totalTaxes:       totalTaxes,
		total:            total,
	}, nil
}
