package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/Kotaro7750/rakuten-security-exporter/proto"
	"github.com/bojanz/currency"
)

type DepositWithdrawalHistory []DepositWithdrawal

func constructDepositWithdrawalHistory(protoWithdrawalHistory *proto.ListWithdrawalHistoriesResponse) (DepositWithdrawalHistory, error) {
	depositWithdrawalHistory := make(DepositWithdrawalHistory, 0)

	for _, protoWithdrawal := range protoWithdrawalHistory.GetHistory() {
		depositWithdrawal, err := newDepositWithdrawal(protoWithdrawal)
		if err != nil {
			return DepositWithdrawalHistory{}, err
		}

		depositWithdrawalHistory = append(depositWithdrawalHistory, depositWithdrawal)
	}

	return depositWithdrawalHistory, nil
}

func (dwh *DepositWithdrawalHistory) totalInvestmentAmount(targetCurrencyCode string, rateManager *RateManager) (currency.Amount, error) {
	totalInvestmentAmount, err := currency.NewAmount("0", targetCurrencyCode)
	if err != nil {
		return currency.Amount{}, err
	}

	for _, depositWithdrawal := range *dwh {
		totalInvestmentAmount, err = addAmount(totalInvestmentAmount, depositWithdrawal.amount, targetCurrencyCode, rateManager)
		if err != nil {
			return currency.Amount{}, err
		}
	}

	return totalInvestmentAmount, nil
}

type DepositWithdrawal struct {
	date   time.Time
	amount currency.Amount
}

func newDepositWithdrawal(protoWithdrawalHistory *proto.WithdrawalHistory) (DepositWithdrawal, error) {
	date, err := time.Parse("2006/01/02", protoWithdrawalHistory.GetDate())
	if err != nil {
		return DepositWithdrawal{}, nil
	}

	inOrOut := protoWithdrawalHistory.GetType()
	amountStr := strconv.FormatUint(protoWithdrawalHistory.GetAmount(), 10)
	if inOrOut == "out" {
		amountStr = strings.Join([]string{"-", amountStr}, "")
	}

	// TODO
	amount, err := currency.NewAmount(amountStr, "JPY")
	if err != nil {
		return DepositWithdrawal{}, err
	}

	return DepositWithdrawal{date, amount}, nil
}
