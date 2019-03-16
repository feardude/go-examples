package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type cbrCurrency struct {
	CbrCode string `xml:"Vcode"`
	EngCode string `xml:"VcharCode"`
	Name    string `xml:"Vname"`
}

// Unmarshal parses input XML to currency
func unmarshal(input *[]byte) {
	var currencies struct {
		CbrCurrencies []cbrCurrency `xml:"Body>EnumValutesXMLResponse>EnumValutesXMLResult>ValuteData>EnumValutes"`
	}
	err := xml.Unmarshal(*input, &currencies)
	if err == nil {
		currencies.CbrCurrencies = filter(currencies.CbrCurrencies)
		for _, c := range currencies.CbrCurrencies {
			fmt.Println(toString(c))
		}
		return
	}
	panic(err)
}

func toString(cbrCurrency cbrCurrency) string {
	return fmt.Sprintf("%s – %s:%s",
		strings.TrimSpace(cbrCurrency.Name),
		strings.TrimSpace(cbrCurrency.EngCode),
		strings.TrimSpace(cbrCurrency.CbrCode))
}

func filter(cbrCurrencies []cbrCurrency) []cbrCurrency {
	filtered := make([]cbrCurrency, 0)
	for _, currency := range cbrCurrencies {
		if currency.EngCode != "" {
			filtered = append(filtered, currency)
		}
	}
	return filtered
}
