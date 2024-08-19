package main

import (
	"math"
	"strconv"
	"time"

	"github.com/bojanz/currency"
)

func addAmount(a currency.Amount, b currency.Amount, targetCurrencyCode string, rateManager *RateManager) (currency.Amount, error) {
	if a.IsZero() && b.IsZero() {
		return currency.NewAmount("0", targetCurrencyCode)
	} else if a.IsZero() {
		return rateManager.Convert(b, targetCurrencyCode)
	} else if b.IsZero() {
		return rateManager.Convert(a, targetCurrencyCode)
	}

	convertedA, err := rateManager.Convert(a, targetCurrencyCode)
	if err != nil {
		return currency.Amount{}, err
	}

	convertedB, err := rateManager.Convert(b, targetCurrencyCode)
	if err != nil {
		return currency.Amount{}, err
	}

	return convertedA.Add(convertedB)
}

func amountRatio(dividend currency.Amount, divisor currency.Amount, rateManager *RateManager) (float64, error) {
	// Align currency to divisor
	rate, err := rateManager.GetRate(dividend.CurrencyCode(), divisor.CurrencyCode())
	if err != nil {
		return 1, err
	}

	convertedDividend, err := dividend.Convert(dividend.CurrencyCode(), rate)
	if err != nil {
		return 1, err
	}

	dividendFloat, err := strconv.ParseFloat(convertedDividend.Number(), 64)
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

func calcAnnualReturn(totalReturn float64, startDate time.Time, endDate time.Time) float64 {
	passedYears := float64(diffMonth(startDate, endDate)) / 12

	return math.Pow(totalReturn, 1/passedYears)
}

func diffMonth(startDate time.Time, endDate time.Time) int64 {
	diffYear := endDate.Year() - startDate.Year()

	if diffYear == 0 {
		return int64(endDate.Month() - startDate.Month() + 1)
	} else {
		return int64((diffYear-1)*12) + int64(12-startDate.Month()+endDate.Month()+1)
	}
}

func monthsOfPastYear(containJustAnYearAgo bool) []struct {
	year  int
	month int
} {
	now := time.Now()
	thisYear := now.Year()
	thisMonth := now.Month()

	ret := make([]struct {
		year  int
		month int
	}, 0)

	iMax := 12
	if containJustAnYearAgo {
		iMax = 13
	}

	for i := 0; i < iMax; i++ {
		startOffset := 0
		if containJustAnYearAgo {
			startOffset = -1
		}

		year := (thisYear*12 + int(thisMonth) - 12 + i + startOffset) / 12
		month := (thisYear*12+int(thisMonth)-12+i+startOffset)%12 + 1

		ret = append(ret, struct {
			year  int
			month int
		}{year, month})
	}

	return ret
}

func OrMonthBoolArray(a [12]bool, b[12]bool) [12]bool {
  var or [12]bool
  for i := 0; i < 12; i++ {
    or[i] = a[i] || b[i]
  }

  return or
}
