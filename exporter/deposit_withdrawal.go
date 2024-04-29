package main

import (
	"log"
	"strconv"

	"github.com/Kotaro7750/rakuten-security-exporter/proto"
	"github.com/bojanz/currency"
)

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
