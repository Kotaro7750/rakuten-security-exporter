package main

import (
	"math"
	"strconv"
	"time"

	"github.com/bojanz/currency"
)

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
