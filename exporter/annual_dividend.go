package main

import (
	"github.com/bojanz/currency"
	"strconv"
)

type AnnualDividend [12]currency.Amount

func constructConstantAnnualDividend(dividendMonth [12]bool, amount currency.Amount, targetCurrencyCode string) AnnualDividend {
	annualDividend := AnnualDividend{}

	for monthIndex, isDividendMonth := range dividendMonth {
		if isDividendMonth {
			annualDividend[monthIndex] = amount
		} else {
			currency.NewAmount("0", targetCurrencyCode)
		}
	}

	return annualDividend
}

func calcAnnualDividendGrowth(thisYear *AnnualDividend, prevYear *AnnualDividend, rateManager *RateManager) (float64, error) {
	thisYearTotal, err := thisYear.calcTotal()
	if err != nil {
		return 1, err
	}

	prevYearTotal, err := prevYear.calcTotal()
	if err != nil {
		return 1, err
	}

	return amountRatio(thisYearTotal, prevYearTotal, rateManager)
}

func addAnnualDividend(a AnnualDividend, b AnnualDividend, targetCurrencyCode string, rateManager *RateManager) (AnnualDividend, error) {
	added := AnnualDividend{}

	for i := 0; i < 12; i++ {
		addedAmount, err := addAmount(a[i], b[i], targetCurrencyCode, rateManager)
		if err != nil {
			return AnnualDividend{}, err
		}
		added[i] = addedAmount
	}

	return added, nil
}

func (ad *AnnualDividend) calcTotal() (currency.Amount, error) {
	total := currency.Amount{}

	for _, dividend := range ad {
		newTotal, err := total.Add(dividend)
		if err != nil {
			return currency.Amount{}, err
		}

		total = newTotal
	}

	return total, nil
}

func (ad *AnnualDividend) calcAverageUnitPrice() (currency.Amount, error) {
	count := 0

	for _, amount := range ad {
		if !amount.IsZero() {
			count++
		}
	}

	if count == 0 {
		return ad.calcTotal()
	}

	total, err := ad.calcTotal()
	if err != nil {
		return currency.Amount{}, err
	}

	return total.Div(strconv.FormatInt(int64(count), 10))
}
