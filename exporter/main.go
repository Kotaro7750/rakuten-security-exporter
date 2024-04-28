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

	withdrawal_history, err := client.ListWithdrawalHistories(ctx, &ListWithdrawalHistoriesRequest{})
	if err != nil {
		log.Fatalf("error %v", err)
	}

	log.Printf("response %v", withdrawal_history)

	dividend_history, err := client.ListDividendHistories(ctx, &ListDividendHistoriesRequest{})
	if err != nil {
		log.Fatalf("error %v", err)
	}

	log.Printf("response %v", dividend_history)
}
