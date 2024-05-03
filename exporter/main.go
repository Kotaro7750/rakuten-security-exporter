package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Kotaro7750/rakuten-security-exporter/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Metrics struct {
	totalReturn                              prometheus.Gauge
	totalReturnAnnual                        prometheus.Gauge
	performanceExcludingCurrencyImpact       prometheus.Gauge
	performanceExcludingCurrencyImpactAnnual prometheus.Gauge
}

var registry prometheus.Registry

func main() {
	threadSafeInvestmentReport := ThreadSafeInvestmentReport{}

	registry = *prometheus.NewRegistry()
	metrics := Metrics{
		totalReturn: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "rakutensecurity",
			Name:      "total_return",
		}),
		totalReturnAnnual: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "rakutensecurity",
			Name:      "total_return_annual",
		}),
		performanceExcludingCurrencyImpact: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "rakutensecurity",
			Name:      "performance_excluding_currency_impact",
		}),
		performanceExcludingCurrencyImpactAnnual: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "rakutensecurity",
			Name:      "performance_excluding_currency_impact_annual",
		}),
	}

	err := scrapeAndSetMetrics(&threadSafeInvestmentReport, &metrics)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	err = registry.Register(metrics.totalReturn)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	err = registry.Register(metrics.totalReturnAnnual)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	err = registry.Register(metrics.performanceExcludingCurrencyImpact)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	err = registry.Register(metrics.performanceExcludingCurrencyImpactAnnual)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	c := cron.New()
	err = c.AddFunc("*/30 * * * * *", func() {
		scrapeAndSetMetrics(&threadSafeInvestmentReport, &metrics)
		if err != nil {
			log.Fatalf("error %v", err)
		}
	})

	if err != nil {
		log.Fatalf("error %v", err)
	}

	c.Start()

	http.Handle("/metrics", promhttp.HandlerFor(&registry, promhttp.HandlerOpts{Registry: &registry}))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error %v", err)
	}
}

func scrape() (InvestmentReport, error) {
	var opt []grpc.DialOption = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial("localhost:50051", opt...)
	if err != nil {
		return InvestmentReport{}, err
	}
	defer conn.Close()

	client := proto.NewRakutenSecurityScraperClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	total_asset, err := client.TotalAssets(ctx, &proto.TotalAssetRequest{})
	if err != nil {
		return InvestmentReport{}, err
	}

	withdrawal_history, err := client.ListWithdrawalHistories(ctx, &proto.ListWithdrawalHistoriesRequest{})
	if err != nil {
		return InvestmentReport{}, err
	}

	dividend_history, err := client.ListDividendHistories(ctx, &proto.ListDividendHistoriesRequest{})
	if err != nil {
		return InvestmentReport{}, err
	}

	return ConstructInvestmentReport(total_asset, withdrawal_history, dividend_history)
}

func scrapeAndSetMetrics(threadSafeInvestmentReport *ThreadSafeInvestmentReport, metrics *Metrics) error {
	investmentReport, err := scrape()
	if err != nil {
		return nil
	}

	threadSafeInvestmentReport.mu.Lock()
	defer threadSafeInvestmentReport.mu.Unlock()

	threadSafeInvestmentReport.InvestmentReport = investmentReport

	performance, err := threadSafeInvestmentReport.ConstructPerformanceReport(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local), "USD")
	if err != nil {
		log.Fatalf("error %v", err)
	}

	metrics.totalReturn.Set(performance.TotalReturn.total)
	metrics.totalReturnAnnual.Set(performance.TotalReturn.annual)
	metrics.performanceExcludingCurrencyImpact.Set(performance.PerformanceExcludingCurrencyImpact.total)
	metrics.performanceExcludingCurrencyImpactAnnual.Set(performance.PerformanceExcludingCurrencyImpact.annual)

	return nil
}

type ThreadSafeInvestmentReport struct {
	InvestmentReport
	mu sync.Mutex
}

type InvestmentReport struct {
	asset                       Asset
	depositAndWithdrawalHistory DepositWithdrawalHistory
	dividendHistory             DividendHistory
	rateManager                 RateManager
}

type PerformanceReport struct {
	PerformanceExcludingCurrencyImpact Return
	TotalReturn                        Return
}

type Return struct {
	total  float64
	annual float64
}

func NewReturn(totalReturn float64, startDate time.Time) Return {
	return Return{
		total:  totalReturn,
		annual: calcAnnualReturn(totalReturn, startDate, time.Now()),
	}
}

func ConstructInvestmentReport(
	assetResponse *proto.TotalAssetResponse,
	withdrawalHistoryResponse *proto.ListWithdrawalHistoriesResponse,
	dividendHistoryResponse *proto.ListDividendHistoriesResponse,
) (InvestmentReport, error) {

	asset, err := constructAsset(assetResponse)
	if err != nil {
		return InvestmentReport{}, err
	}

	depositAndWithdrawalHistory, err := constructDepositWithdrawalHistory(withdrawalHistoryResponse)
	if err != nil {
		return InvestmentReport{}, err
	}

	dividendHistory, err := ConstructDividendHistory(dividendHistoryResponse)
	if err != nil {
		return InvestmentReport{}, err
	}

	rateManager := NewRateManager()

	err = rateManager.RegisterRate("USD", "JPY", 156.45)
	if err != nil {
		return InvestmentReport{}, err
	}

	return InvestmentReport{asset, depositAndWithdrawalHistory, dividendHistory, rateManager}, nil
}

func (ir *InvestmentReport) ConstructPerformanceReport(startDate time.Time, targetCurrencyCode string) (PerformanceReport, error) {
	assetSummary, err := ir.asset.Summarize(targetCurrencyCode, &ir.rateManager)
	if err != nil {
		return PerformanceReport{}, err
	}

	performanceExcludingCurrencyImpact, err := assetSummary.PerformanceExcludingCurrencyImpact(&ir.rateManager)
	if err != nil {
		return PerformanceReport{}, err
	}

	totalInvestmentAmount, err := ir.depositAndWithdrawalHistory.totalInvestmentAmount(targetCurrencyCode, &ir.rateManager)
	if err != nil {
		return PerformanceReport{}, err
	}

	totalReturn, err := amountRatio(assetSummary.totalPrice, totalInvestmentAmount, &ir.rateManager)
	if err != nil {
		return PerformanceReport{}, err
	}

	return PerformanceReport{
		PerformanceExcludingCurrencyImpact: NewReturn(performanceExcludingCurrencyImpact, startDate),
		TotalReturn:                        NewReturn(totalReturn, startDate),
	}, nil
}
