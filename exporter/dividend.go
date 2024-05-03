package main

import (
	"sort"
	"strconv"
	"strings"
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
	account          string
	assetType        string
	ticker           string
	name             string
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
		account:          dividendHistory.GetAccount(),
		assetType:        dividendHistory.GetType(),
		ticker:           dividendHistory.GetTicker(),
		name:             dividendHistory.GetName(),
		count:            dividendHistory.GetCount(),
		unitPrice:        unitPrice,
		totalBeforeTaxed: totalBeforeTaxed,
		totalTaxes:       totalTaxes,
		total:            total,
	}, nil
}

func (d *Dividend) identifier() string {
	return strings.Join([]string{d.assetType, d.ticker, d.name}, " ")
}

type DividendReportAnnual struct {
	total           AnnualDividend
	totalGrowthRate float64
	security        map[Security]struct {
		total                  AnnualDividend
		totalGrowthRate        float64
		unitPrice              AnnualDividend
		averageUnitPrice       currency.Amount
		averageUnitPriceGrowth float64
	}
}

func newDividendReportAnnual() DividendReportAnnual {
	return DividendReportAnnual{
		security: make(map[Security]struct {
			total                  AnnualDividend
			totalGrowthRate        float64
			unitPrice              AnnualDividend
			averageUnitPrice       currency.Amount
			averageUnitPriceGrowth float64
		}),
	}
}

type DividendReport struct {
	actual map[int]DividendReportAnnual
}

func newDividendReport() DividendReport {
	return DividendReport{actual: make(map[int]DividendReportAnnual)}
}

func (dh *DividendHistory) constructDividendReport(targetCurrencyCode string, rateManager *RateManager) (DividendReport, error) {
	dividendReport := newDividendReport()

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

		// TODO
		security := newSecurity(dividend.assetType, dividend.ticker, dividend.name)
		securityDividendReportAnnual := dividendReportAnnual.security[security]

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

		dividendReportAnnual.security[security] = securityDividendReportAnnual
		dividendReport.actual[year] = dividendReportAnnual
	}

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
