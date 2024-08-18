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

type DividendReportAnnual struct {
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

func newDividendReportAnnual() DividendReportAnnual {
	return DividendReportAnnual{
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

func estimateDividend(securityEstimatedUnitPrice map[string]struct {
	year             int
	averageUnitPrice currency.Amount
}, securityDividendMonth map[string][12]bool, assetCount map[Security]float64, targetCurrencyCode string, rateManager *RateManager) (DividendEstimation, error) {

	totalDividend := AnnualDividend{}
	securityDividend := make(map[string]struct {
		total     AnnualDividend
		unitPrice AnnualDividend
	}, 0)

	for security, count := range assetCount {
		estimatedUnitPrice, exists := securityEstimatedUnitPrice[security.identifier()]

		if exists {
			securityEstimation := securityDividend[security.identifier()]

			estimatedDividendPerMonth, err := estimatedUnitPrice.averageUnitPrice.Mul(strconv.FormatFloat(count, 'f', -1, 64))
			if err != nil {
				return DividendEstimation{}, err
			}

			dividendMonth, _ := securityDividendMonth[security.identifier()]

			securityEstimation.total = constructConstantAnnualDividend(dividendMonth, estimatedDividendPerMonth, targetCurrencyCode)
			securityEstimation.unitPrice = constructConstantAnnualDividend(dividendMonth, estimatedUnitPrice.averageUnitPrice, targetCurrencyCode)

			totalDividend, err = addAnnualDividend(totalDividend, securityEstimation.total, targetCurrencyCode, rateManager)
			if err != nil {
				return DividendEstimation{}, err
			}

			securityDividend[security.identifier()] = securityEstimation
		}
	}

	return DividendEstimation{totalDividend, securityDividend}, nil
}

func estimateSecurityDividendMonth(annualDividend map[int]AnnualDividend) [12]bool {
	dividendMonth := [12]bool{}

	monthAndYears := monthsOfPastYear(true)
	for _, monthAndYear := range monthAndYears {
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
	actual   map[int]DividendReportAnnual
	estimate DividendEstimation
}

func newDividendReport() DividendReport {
	return DividendReport{actual: make(map[int]DividendReportAnnual)}
}

func (dh *DividendHistory) constructDividendReport(assetCount map[Security]float64, targetCurrencyCode string, rateManager *RateManager) (DividendReport, error) {
	dividendReport := newDividendReport()

	latestAverageUnitPrice := make(map[string]struct {
		year             int
		averageUnitPrice currency.Amount
	}, 0)

	dividendMonth := make(map[string][12]bool, 0)

	for _, dividend := range *dh {
		year := dividend.date.Year()
		monthIndex := dividend.date.Month() - 1

		dividendReportAnnual, exists := dividendReport.actual[year]
		if !exists {
			dividendReportAnnual = newDividendReportAnnual()
		}

		annualTotal, err := addAmount(dividend.total, dividendReportAnnual.total[monthIndex], targetCurrencyCode, rateManager)
		if err != nil {
			return DividendReport{}, err
		}

		dividendReportAnnual.total[monthIndex] = annualTotal

		securityDividendReportAnnual := dividendReportAnnual.security[dividend.security.identifier()]

		securityAnnualTotal, err := addAmount(dividend.total, securityDividendReportAnnual.total[monthIndex], targetCurrencyCode, rateManager)
		if err != nil {
			return DividendReport{}, err
		}
		securityDividendReportAnnual.total[monthIndex] = securityAnnualTotal

		securityAnnualUnitPrice, err := addAmount(dividend.unitPrice, securityDividendReportAnnual.unitPrice[monthIndex], targetCurrencyCode, rateManager)
		if err != nil {
			return DividendReport{}, err
		}
		securityDividendReportAnnual.unitPrice[monthIndex] = securityAnnualUnitPrice

		securityAnnualAverageUnitPrice, err := securityDividendReportAnnual.unitPrice.calcAverageUnitPrice()
		if err != nil {
			return DividendReport{}, err
		}
		securityDividendReportAnnual.averageUnitPrice = securityAnnualAverageUnitPrice

		dividendReportAnnual.security[dividend.security.identifier()] = securityDividendReportAnnual
		dividendReport.actual[year] = dividendReportAnnual

		latest, exists := latestAverageUnitPrice[dividend.security.identifier()]
		if exists {
			if dividend.date.Year() >= latest.year {
				latest.averageUnitPrice = securityAnnualAverageUnitPrice
			}
		} else {
			latest.averageUnitPrice = securityAnnualAverageUnitPrice
		}
		latestAverageUnitPrice[dividend.security.identifier()] = latest
	}

	for security, _ := range latestAverageUnitPrice {
		securityAnnualDividend := make(map[int]AnnualDividend, 0)

		for year, annualDividend := range dividendReport.actual {
			if _, exists := annualDividend.security[security]; exists {
				securityAnnualDividend[year] = annualDividend.security[security].total
			}
		}

		dividendMonth[security] = estimateSecurityDividendMonth(securityAnnualDividend)
	}

	dividendEstimation, err := estimateDividend(latestAverageUnitPrice, dividendMonth, assetCount, targetCurrencyCode, rateManager)
	if err != nil {
		return DividendReport{}, err
	}
	dividendReport.estimate = dividendEstimation

	// Calculate growth
	for year, dividendReportAnnual := range dividendReport.actual {
		prevYearDividendReportAnnual, exists := dividendReport.actual[year-1]

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

		dividendReport.actual[year] = dividendReportAnnual
	}

	return dividendReport, nil
}

type AnnualDividend [12]currency.Amount

func constructConstantAnnualDividend(dividendMonth [12]bool, amount currency.Amount, targetCurrencyCode string) AnnualDividend {
	annualDividend := AnnualDividend{}

	for monthIndex, isDividendMonth := range dividendMonth {
		if isDividendMonth {
			annualDividend[monthIndex] = amount
		} else {
			currency.NewAmount("0", targetCurrencyCode)
		}
	}

	return annualDividend
}

func calcAnnualDividendGrowth(thisYear *AnnualDividend, prevYear *AnnualDividend, rateManager *RateManager) (float64, error) {
	thisYearTotal, err := thisYear.calcTotal()
	if err != nil {
		return 1, err
	}

	prevYearTotal, err := prevYear.calcTotal()
	if err != nil {
		return 1, err
	}

	return amountRatio(thisYearTotal, prevYearTotal, rateManager)
}

func addAnnualDividend(a AnnualDividend, b AnnualDividend, targetCurrencyCode string, rateManager *RateManager) (AnnualDividend, error) {
	added := AnnualDividend{}

	for i := 0; i < 12; i++ {
		addedAmount, err := addAmount(a[i], b[i], targetCurrencyCode, rateManager)
		if err != nil {
			return AnnualDividend{}, err
		}
		added[i] = addedAmount
	}

	return added, nil
}

func (ad *AnnualDividend) calcTotal() (currency.Amount, error) {
	total := currency.Amount{}

	for _, dividend := range ad {
		newTotal, err := total.Add(dividend)
		if err != nil {
			return currency.Amount{}, err
		}

		total = newTotal
	}

	return total, nil
}

func (ad *AnnualDividend) calcAverageUnitPrice() (currency.Amount, error) {
	count := 0

	for _, amount := range ad {
		if !amount.IsZero() {
			count++
		}
	}

	if count == 0 {
		return ad.calcTotal()
	}

	total, err := ad.calcTotal()
	if err != nil {
		return currency.Amount{}, err
	}

	return total.Div(strconv.FormatInt(int64(count), 10))
}
