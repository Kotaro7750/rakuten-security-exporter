package main

import (
	"fmt"
	"strconv"

	"github.com/bojanz/currency"
)

type RateManager struct {
	rate map[string]string
}

type RateNotFoundError struct {
	originalCurrencyCode string
	targetCurrencyCode   string
}

func (e *RateNotFoundError) Error() string {
	return fmt.Sprintf("Exchange rate from %s to %s not found", e.originalCurrencyCode, e.targetCurrencyCode)
}

func NewRateManager() RateManager {
	return RateManager{rate: make(map[string]string)}
}

func (nr *RateManager) RegisterRate(originalCurrencyCode string, targetCurrencyCode string, rate float64) error {
	if err := checkExchangeValidity(originalCurrencyCode, targetCurrencyCode, rate); err != nil {
		return err
	}
	if err := checkExchangeValidity(targetCurrencyCode, originalCurrencyCode, 1/rate); err != nil {
		return err
	}

	rateStr := strconv.FormatFloat(rate, 'f', -1, 64)
	inverseRateStr := strconv.FormatFloat(1/rate, 'f', -1, 64)

	nr.rate[fmt.Sprintf("%s/%s", targetCurrencyCode, originalCurrencyCode)] = rateStr
	nr.rate[fmt.Sprintf("%s/%s", originalCurrencyCode, targetCurrencyCode)] = inverseRateStr
	return nil
}

func (nr *RateManager) GetRate(originalCurrencyCode string, targetCurrencyCode string) (string, error) {
	err := checkExchangeValidity(originalCurrencyCode, targetCurrencyCode, 1)
	if err != nil {
		return "", err
	}

	if originalCurrencyCode == targetCurrencyCode {
		return "1", nil
	}

	rate, exists := nr.rate[fmt.Sprintf("%s/%s", targetCurrencyCode, originalCurrencyCode)]

	if !exists {
		return "", &RateNotFoundError{originalCurrencyCode, targetCurrencyCode}
	}

	return rate, nil
}

func checkExchangeValidity(originalCurrencyCode string, targetCurrencyCode string, rate float64) error {
	original, err := currency.NewAmount("1", originalCurrencyCode)
	if err != nil {
		return err
	}

	_, err = original.Convert(targetCurrencyCode, strconv.FormatFloat(rate, 'f', -1, 64))
	if err != nil {
		return err
	}

	return nil
}
