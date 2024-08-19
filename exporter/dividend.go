package main

import (
	"sort"
	"strconv"
	"time"

	"github.com/Kotaro7750/rakuten-security-exporter/proto"
	"github.com/bojanz/currency"
)

type DividendHistory []Dividend

func ConstructDividendHistory(protoDividendHistoryResponse *proto.ListDividendHistoriesResponse) (DividendHistory, error) {
	dividendHistory := make(DividendHistory, 0)

	for _, protoDividendHistory := range protoDividendHistoryResponse.GetHistory() {
		dividend, err := newDividend(protoDividendHistory)
		if err != nil {
			return DividendHistory{}, err
		}

		dividendHistory = append(dividendHistory, dividend)
	}

	sort.SliceStable(dividendHistory, func(i, j int) bool {
		return dividendHistory[i].date.Unix() < dividendHistory[j].date.Unix()
	})

	return dividendHistory, nil
}

type Dividend struct {
	date             time.Time
	security         Security
	count            float64
	unitPrice        currency.Amount
	totalBeforeTaxed currency.Amount
	totalTaxes       currency.Amount
	total            currency.Amount
}

func newDividend(dividendHistory *proto.DividendHistory) (Dividend, error) {
	date, err := time.Parse("2006/01/02", dividendHistory.GetDate())
	if err != nil {
		return Dividend{}, nil
	}

	dividendCurrency := dividendHistory.GetCurrency()

	unitPrice, err := currency.NewAmount(strconv.FormatFloat(dividendHistory.GetDividendUnitprice(), 'f', -1, 64), dividendCurrency)
	if err != nil {
		return Dividend{}, err
	}

	totalBeforeTaxed, err := currency.NewAmount(strconv.FormatFloat(dividendHistory.GetDividendTotalBeforeTaxes(), 'f', -1, 64), dividendCurrency)
	if err != nil {
		return Dividend{}, err
	}

	totalTaxes, err := currency.NewAmount(strconv.FormatFloat(dividendHistory.GetTotalTaxes(), 'f', -1, 64), dividendCurrency)
	if err != nil {
		return Dividend{}, err
	}

	total, err := currency.NewAmount(strconv.FormatFloat(dividendHistory.GetDividendTotal(), 'f', -1, 64), dividendCurrency)
	if err != nil {
		return Dividend{}, err
	}

	return Dividend{date: date,
		security:         newSecurity(dividendHistory.GetAccount(), dividendHistory.GetType(), dividendHistory.GetTicker(), dividendHistory.GetName()),
		count:            dividendHistory.GetCount(),
		unitPrice:        unitPrice,
		totalBeforeTaxed: totalBeforeTaxed,
		totalTaxes:       totalTaxes,
		total:            total,
	}, nil
}

type DividendResultAnnual struct {
	total           AnnualDividend
	totalGrowthRate float64
	security        map[string]struct {
		total                  AnnualDividend
		totalGrowthRate        float64
		unitPrice              AnnualDividend
		averageUnitPrice       currency.Amount
		averageUnitPriceGrowth float64
	}
}

func newDividendResultAnnual() DividendResultAnnual {
	return DividendResultAnnual{
		security: make(map[string]struct {
			total                  AnnualDividend
			totalGrowthRate        float64
			unitPrice              AnnualDividend
			averageUnitPrice       currency.Amount
			averageUnitPriceGrowth float64
		}),
	}
}

type DividendEstimation struct {
	total    AnnualDividend
	security map[string]struct {
		total     AnnualDividend
		unitPrice AnnualDividend
	}
}

func estimateDividend(securityEstimatedUnitPrice map[string]currency.Amount, securityDividendMonth map[string][12]bool, assetCount map[Security]float64, targetCurrencyCode string, rateManager *RateManager) (DividendEstimation, error) {

	totalDividend := AnnualDividend{}
	securityDividend := make(map[string]struct {
		total     AnnualDividend
		unitPrice AnnualDividend
	}, 0)

	for security, count := range assetCount {
		estimatedUnitPrice, exists := securityEstimatedUnitPrice[security.identifier()]

		if exists {
			securityEstimation := securityDividend[security.identifierWithAccount()]

			estimatedDividendPerMonth, err := estimatedUnitPrice.Mul(strconv.FormatFloat(count, 'f', -1, 64))
			if err != nil {
				return DividendEstimation{}, err
			}

			dividendMonth, _ := securityDividendMonth[security.identifier()]

			securityEstimation.total = constructConstantAnnualDividend(dividendMonth, estimatedDividendPerMonth, targetCurrencyCode)
			securityEstimation.unitPrice = constructConstantAnnualDividend(dividendMonth, estimatedUnitPrice, targetCurrencyCode)

			totalDividend, err = addAnnualDividend(totalDividend, securityEstimation.total, targetCurrencyCode, rateManager)
			if err != nil {
				return DividendEstimation{}, err
			}

			securityDividend[security.identifierWithAccount()] = securityEstimation
		}
	}

	return DividendEstimation{totalDividend, securityDividend}, nil
}

func estimateDividendMonth(annualDividend map[int]AnnualDividend) [12]bool {
	dividendMonth := [12]bool{}

	estimationSourceDateRange := monthsOfPastYear(true)
	for _, monthAndYear := range estimationSourceDateRange {
		year := monthAndYear.year
		monthIndex := monthAndYear.month - 1

		annualDividendOfYear, exists := annualDividend[year]
		if exists && !annualDividendOfYear[monthIndex].IsZero() {
			dividendMonth[monthIndex] = true
		}
	}
	return dividendMonth
}

type DividendReport struct {
	result   map[int]DividendResultAnnual
	estimate DividendEstimation
}

func newDividendReport() DividendReport {
	return DividendReport{result: make(map[int]DividendResultAnnual)}
}

func (dh *DividendHistory) constructDividendReport(assetCount map[Security]float64, targetCurrencyCode string, rateManager *RateManager) (DividendReport, error) {
	dividendReport := newDividendReport()

	latestAverageUnitPrice := make(map[string]currency.Amount, 0)

	dividendMonth := make(map[string][12]bool, 0)

	for _, dividend := range *dh {
		year := dividend.date.Year()
		monthIndex := dividend.date.Month() - 1

		dividendResultAnnual, exists := dividendReport.result[year]
		if !exists {
			dividendResultAnnual = newDividendResultAnnual()
		}

		// First, update total result
		annualTotal, err := addAmount(dividend.total, dividendResultAnnual.total[monthIndex], targetCurrencyCode, rateManager)
		if err != nil {
			return DividendReport{}, err
		}
		dividendResultAnnual.total[monthIndex] = annualTotal

		// Next, update result per security
		// Total, annual unit price and unit price per single pay
		dividendResultAnnualOfSecurity := dividendResultAnnual.security[dividend.security.identifierWithAccount()]

		securityAnnualTotal, err := addAmount(dividend.total, dividendResultAnnualOfSecurity.total[monthIndex], targetCurrencyCode, rateManager)
		if err != nil {
			return DividendReport{}, err
		}
		dividendResultAnnualOfSecurity.total[monthIndex] = securityAnnualTotal

		securityAnnualUnitPrice, err := addAmount(dividend.unitPrice, dividendResultAnnualOfSecurity.unitPrice[monthIndex], targetCurrencyCode, rateManager)
		if err != nil {
			return DividendReport{}, err
		}
		dividendResultAnnualOfSecurity.unitPrice[monthIndex] = securityAnnualUnitPrice

		securityAnnualAverageUnitPrice, err := dividendResultAnnualOfSecurity.unitPrice.calcAverageUnitPrice()
		if err != nil {
			return DividendReport{}, err
		}
		dividendResultAnnualOfSecurity.averageUnitPrice = securityAnnualAverageUnitPrice
    // Because dividend history is sorted, just overwriting is equal to updating latest
		latestAverageUnitPrice[dividend.security.identifier()] = securityAnnualAverageUnitPrice

		dividendResultAnnual.security[dividend.security.identifierWithAccount()] = dividendResultAnnualOfSecurity
		dividendReport.result[year] = dividendResultAnnual
	}

	for security, _ := range assetCount {
		securityAnnualDividend := make(map[int]AnnualDividend, 0)

		for year, annualDividend := range dividendReport.result {
			if _, exists := annualDividend.security[security.identifierWithAccount()]; exists {
				securityAnnualDividend[year] = annualDividend.security[security.identifierWithAccount()].total
			}
		}

    // Not distinguish between different account with same ticker
		dividendMonth[security.identifier()] = OrMonthBoolArray(dividendMonth[security.identifier()], estimateDividendMonth(securityAnnualDividend))
	}

	dividendEstimation, err := estimateDividend(latestAverageUnitPrice, dividendMonth, assetCount, targetCurrencyCode, rateManager)
	if err != nil {
		return DividendReport{}, err
	}
	dividendReport.estimate = dividendEstimation

	// Calculate growth
	for year, dividendReportAnnual := range dividendReport.result {
		prevYearDividendReportAnnual, exists := dividendReport.result[year-1]

		if exists {
			totalGrowth, err := calcAnnualDividendGrowth(&dividendReportAnnual.total, &prevYearDividendReportAnnual.total, rateManager)
			if err != nil {
				return DividendReport{}, err
			}

			dividendReportAnnual.totalGrowthRate = totalGrowth

			for security, securityDividendReportAnnual := range dividendReportAnnual.security {
				prevYearSecurityDividendReportAnnual, exists := prevYearDividendReportAnnual.security[security]

				if exists {
					totalGrowth, err := calcAnnualDividendGrowth(&securityDividendReportAnnual.total, &prevYearSecurityDividendReportAnnual.total, rateManager)
					if err != nil {
						return DividendReport{}, err
					}

					averageUnitPriceGrowth, err := amountRatio(securityDividendReportAnnual.averageUnitPrice, prevYearSecurityDividendReportAnnual.averageUnitPrice, rateManager)
					if err != nil {
						return DividendReport{}, err
					}

					securityDividendReportAnnual.totalGrowthRate = totalGrowth
					securityDividendReportAnnual.averageUnitPriceGrowth = averageUnitPriceGrowth

				} else {
					securityDividendReportAnnual.totalGrowthRate = 1
					securityDividendReportAnnual.averageUnitPriceGrowth = 1
				}

				dividendReportAnnual.security[security] = securityDividendReportAnnual
			}
		} else {
			dividendReportAnnual.totalGrowthRate = 1

			for security, securityDividendReportAnnual := range dividendReportAnnual.security {
				securityDividendReportAnnual.totalGrowthRate = 1
				securityDividendReportAnnual.averageUnitPriceGrowth = 1

				dividendReportAnnual.security[security] = securityDividendReportAnnual
			}
		}

		dividendReport.result[year] = dividendReportAnnual
	}

	return dividendReport, nil
}
