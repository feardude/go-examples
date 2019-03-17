package main

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// Currency represents currency structure
type Currency struct {
	CodeCbr string `xml:"Vcode"`
	CodeEng string `xml:"VcharCode"`
	NameRus string `xml:"Vname"`
	NameEng string `xml:"VEngname"`
}

// ToString returns formatted Currency
func (currency Currency) ToString() string {
	return fmt.Sprintf("%s (%s) – %s:%s",
		currency.NameEng,
		currency.NameRus,
		currency.CodeEng,
		currency.CodeCbr)
}

// FxRate represents currency rate structure
type FxRate struct {
	Date       time.Time `xml:"CursDate"`
	Multiplier int16     `xml:"Vnom"`
	Value      float32   `xml:"Vcurs"`
}

// ToString returns formatted FxRate
func (fxRate FxRate) ToString(code string) string {
	return fmt.Sprintf("%s: %d %s = %.4f RUB",
		fxRate.Date.Format(time.RFC3339),
		fxRate.Multiplier,
		code,
		fxRate.Value)
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

// FxRates parses input XML to []FxRate and returns result
func FxRates(input *[]byte) []FxRate {
	var fxRates struct {
		FxRates []FxRate `xml:"Body>GetCursDynamicXMLResponse>GetCursDynamicXMLResult>ValuteData>ValuteCursDynamic"`
	}
	err := xml.Unmarshal(*input, &fxRates)
	if err == nil {
		return fxRates.FxRates
	}
	panic(err)
}