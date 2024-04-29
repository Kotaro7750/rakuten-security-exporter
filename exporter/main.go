package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var opt []grpc.DialOption = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial("localhost:50051", opt...)
	if err != nil {
		log.Fatalf("dial error %v", err)
	}
	defer conn.Close()

	client := NewRakutenSecurityScraperClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	total_asset, err := client.TotalAssets(ctx, &TotalAssetRequest{})
	if err != nil {
		log.Fatalf("error %v", err)
	}

	log.Printf("response %v", total_asset)

	withdrawal_history, err := client.ListWithdrawalHistories(ctx, &ListWithdrawalHistoriesRequest{})
	if err != nil {
		log.Fatalf("error %v", err)
	}

  log.Printf("withdrawalStat: %v", constructWithdrawalStatistics(withdrawal_history))

	dividend_history, err := client.ListDividendHistories(ctx, &ListDividendHistoriesRequest{})
	if err != nil {
		log.Fatalf("error %v", err)
	}

	log.Printf("response %v", dividend_history)
}

type WithdrawalSummary struct {
	TotalInvestmentAmount float64
	Currency              string
}

func constructWithdrawalStatistics(withdrawalHistories *ListWithdrawalHistoriesResponse) WithdrawalSummary {
	var totalInvestmentAmount float64 = 0

	for _, withdrawalHistory := range withdrawalHistories.GetHistory() {
		withdrawalType := withdrawalHistory.GetType()
		amount := float64(withdrawalHistory.GetAmount())

		if withdrawalType == "in" {
			totalInvestmentAmount += amount
		} else if withdrawalType == "out" {
			totalInvestmentAmount -= amount
		}
	}

	return WithdrawalSummary{
		TotalInvestmentAmount: totalInvestmentAmount,
		Currency:              "YEN",
	}
}
