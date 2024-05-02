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

func (dh *DividendHistory) constructDividendStatistics(targetCurrency string, rateManager *RateManager) (DividendStatistics, error) {
	dividendStatistics := newDividendStatistics()

	for _, dividend := range *dh {
		err := dividendStatistics.addDividend(&dividend, targetCurrency, rateManager)
		if err != nil {
			return newDividendStatistics(), err
		}
	}

	return dividendStatistics, nil
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

type SecurityAnnualDividend struct {
	total     AnnualDividend
	unitPrice AnnualDividend
}

func (sad *SecurityAnnualDividend) addDividend(dividend *Dividend, targetCurrencyCode string, rateManager *RateManager) error {
	month := dividend.date.Month()

	var err error
	sad.total[month-1], err = addAmount(sad.total[month-1], dividend.total, targetCurrencyCode, rateManager)
	if err != nil {
		return err
	}

	sad.unitPrice[month-1], err = addAmount(sad.unitPrice[month-1], dividend.unitPrice, targetCurrencyCode, rateManager)
	if err != nil {
		return err
	}

	return nil
}

type TotalDividendStatistics struct {
	dividend    AnnualDividend
	totalGrowth float64
}

type SecurityDividendStatistics struct {
	assetType  string
	ticker     string
	name       string
	statistics map[int]struct {
		dividend        SecurityAnnualDividend
		totalGrowth     float64
		unitPriceGrowth float64
	}
}

func (sds *SecurityDividendStatistics) getStatisticsOfYear(year int) (struct {
	dividend        SecurityAnnualDividend
	totalGrowth     float64
	unitPriceGrowth float64
}, bool) {
	statistics, exists := sds.statistics[year]

	if exists {
		return statistics, true
	} else {
		return struct {
			dividend        SecurityAnnualDividend
			totalGrowth     float64
			unitPriceGrowth float64
		}{}, false
	}

}

func (sds *SecurityDividendStatistics) addDividend(dividend *Dividend, targetCurrencyCode string, rateManager *RateManager) error {
	year := dividend.date.Year()

	sad, _ := sds.getStatisticsOfYear(year)

	err := sad.dividend.addDividend(dividend, targetCurrencyCode, rateManager)
	if err != nil {
		return err
	}

	if prevYearStatistics, exists := sds.getStatisticsOfYear(year - 1); exists {
		totalGrowth, err := calcAnnualDividendGrowth(&sad.dividend.total, &prevYearStatistics.dividend.total, rateManager)
		if err != nil {
			return err
		}

		unitPriceGrowth, err := calcAnnualDividendGrowth(&sad.dividend.unitPrice, &prevYearStatistics.dividend.unitPrice, rateManager)
		if err != nil {
			return err
		}

		sad.totalGrowth = totalGrowth
		sad.unitPriceGrowth = unitPriceGrowth
	}

	sds.statistics[year] = sad

	if nextYearStatistics, exists := sds.getStatisticsOfYear(year + 1); exists {
		totalGrowth, err := calcAnnualDividendGrowth(&nextYearStatistics.dividend.total, &sad.dividend.total, rateManager)
		if err != nil {
			return err
		}

		unitPriceGrowth, err := calcAnnualDividendGrowth(&nextYearStatistics.dividend.unitPrice, &sad.dividend.unitPrice, rateManager)
		if err != nil {
			return err
		}

		nextYearStatistics.totalGrowth = totalGrowth
		nextYearStatistics.unitPriceGrowth = unitPriceGrowth

		sds.statistics[year+1] = nextYearStatistics
	}

	return nil
}

type DividendStatistics struct {
	total    map[int]TotalDividendStatistics
	security map[string]SecurityDividendStatistics
}

func newDividendStatistics() DividendStatistics {
	return DividendStatistics{
		total:    make(map[int]TotalDividendStatistics),
		security: make(map[string]SecurityDividendStatistics),
	}
}

func (ds *DividendStatistics) getSecuritydividenStatisticsOfSecurity(assetType, ticker, name string) SecurityDividendStatistics {
	identifier := strings.Join([]string{assetType, ticker, name}, " ")

	sds, exists := ds.security[identifier]
	if exists {
		return sds
	} else {
		return SecurityDividendStatistics{
			assetType: assetType,
			ticker:    ticker,
			name:      name,
			statistics: make(map[int]struct {
				dividend        SecurityAnnualDividend
				totalGrowth     float64
				unitPriceGrowth float64
			}),
		}
	}
}

func (ds *DividendStatistics) addDividend(dividend *Dividend, targetCurrencyCode string, rateManager *RateManager) error {
	err := ds.addDividendToTotal(dividend, targetCurrencyCode, rateManager)
	if err != nil {
		return err
	}

	err = ds.addDividendToSecurity(dividend, targetCurrencyCode, rateManager)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DividendStatistics) addDividendToTotal(dividend *Dividend, targetCurrencyCode string, rateManager *RateManager) error {
	year := dividend.date.Year()

	totalDividendStatistics := ds.total[year]

	month := dividend.date.Month()

	newAmount, err := addAmount(totalDividendStatistics.dividend[month-1], dividend.total, targetCurrencyCode, rateManager)
	if err != nil {
		return nil
	}

	totalDividendStatistics.dividend[month-1] = newAmount
	totalDividendStatistics.totalGrowth = 1

	prevYearStat, exists := ds.total[year-1]

	if exists {
		growth, err := calcAnnualDividendGrowth(&totalDividendStatistics.dividend, &prevYearStat.dividend, rateManager)
		if err != nil {
			return err
		}

		totalDividendStatistics.totalGrowth = growth
	}

	ds.total[year] = totalDividendStatistics

	return nil
}

func (ds *DividendStatistics) addDividendToSecurity(dividend *Dividend, targetCurrencyCode string, rateManager *RateManager) error {
	identifier := dividend.identifier()

	// TODO
	securityDividendStatistics := ds.getSecuritydividenStatisticsOfSecurity(dividend.assetType, dividend.ticker, dividend.name)

	err := securityDividendStatistics.addDividend(dividend, targetCurrencyCode, rateManager)
	if err != nil {
		return err
	}

	ds.security[identifier] = securityDividendStatistics

	return nil
}
