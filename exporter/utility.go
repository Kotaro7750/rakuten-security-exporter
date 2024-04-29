package main

import (
	"strconv"

	"github.com/bojanz/currency"
)

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
