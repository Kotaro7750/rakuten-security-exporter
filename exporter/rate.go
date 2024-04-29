package main

import (
	"fmt"
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

func (nr *RateManager) RegisterRate(originalCurrencyCode string, targetCurrencyCode string, rate string) error {
	// Check argument validity by actually converting
	original, err := currency.NewAmount("1", originalCurrencyCode)
	if err != nil {
		return err
	}

	_, err = original.Convert(targetCurrencyCode, rate)
	if err != nil {
		return err
	}

	nr.rate[fmt.Sprintf("%s/%s", targetCurrencyCode, originalCurrencyCode)] = rate
	return nil
}

func (nr *RateManager) GetRate(originalCurrencyCode string, targetCurrencyCode string) (string, error) {
	original, err := currency.NewAmount("1", originalCurrencyCode)
	if err != nil {
		return "", err
	}

	_, err = original.Convert(targetCurrencyCode, "2.0")
	if err != nil {
		return "", err
	}

	rate, exists := nr.rate[fmt.Sprintf("%s/%s", targetCurrencyCode, originalCurrencyCode)]

	if !exists {
		return "", &RateNotFoundError{originalCurrencyCode, targetCurrencyCode}
	}

	return rate, nil
}
