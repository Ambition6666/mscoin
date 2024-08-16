package domain

import "strings"

const CNY = "CNY"
const JPY = "JPY"
const HKD = "HKD"

type ExchangeRateDomain struct {
}

func (d *ExchangeRateDomain) GetUsdRate(unit string) float64 {
	upper := strings.ToUpper(unit)
	if upper == CNY {
		return 7.00
	} else if upper == JPY {
		return 110.02
	} else if upper == HKD {
		return 7.8497
	}
	return 0
}

func NewExchangeRateDomain() *ExchangeRateDomain {
	return &ExchangeRateDomain{}
}
