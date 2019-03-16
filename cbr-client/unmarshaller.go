package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// Currency represents currency structure
type Currency struct {
	CodeCbr string `xml:"Vcode"`
	CodeEng string `xml:"VcharCode"`
	NameRus string `xml:"Vname"`
	NameEng string `xml:"VEngname"`
}

// ToString returns formatted Currency
func ToString(currency Currency) string {
	return fmt.Sprintf("%s (%s) – %s:%s",
		currency.NameEng,
		currency.NameRus,
		currency.CodeEng,
		currency.CodeCbr)
}

// Currencies parses input XML to []Currency and returns result
func Currencies(input *[]byte) []Currency {
	var currencies struct {
		CbrCurrencies []Currency `xml:"Body>EnumValutesXMLResponse>EnumValutesXMLResult>ValuteData>EnumValutes"`
	}
	err := xml.Unmarshal(*input, &currencies)
	if err == nil {
		return processed(currencies.CbrCurrencies)
	}
	panic(err)
}

func processed(cbrCurrencies []Currency) []Currency {
	processed := make([]Currency, 0)
	for _, c := range cbrCurrencies {
		if c.CodeEng != "" {
			c.CodeCbr = strings.TrimSpace(c.CodeCbr)
			c.CodeEng = strings.TrimSpace(c.CodeEng)
			c.NameEng = strings.TrimSpace(c.NameEng)
			c.NameRus = strings.TrimSpace(c.NameRus)
			processed = append(processed, c)
		}
	}
	return processed
}
