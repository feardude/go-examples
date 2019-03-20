package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var cbrCodeToCurrency map[string]Currency
var codeToFxRate map[string][]FxRate

func main() {
	log.Println("Started CBR loader")
	InitDB()
	loadCurrencies()
	loadFxRates()
	ShutdownDB()
	log.Println("Finished CBR loader")
}

func load(body *strings.Reader) []byte {
	client := http.Client{}
	resp, err := client.Post("http://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx",
		"application/soap+xml", body)
	check(err)

	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	return bytes
}

func loadCurrencies() {
	currencies := GetCurrencies()
	if len(currencies) == 0 {
		log.Fatal("No currencies in DB. Shutting down...")
	}

	cbrCodeToCurrency = make(map[string]Currency)
	for _, currency := range currencies {
		cbrCodeToCurrency[currency.CodeCbr] = currency
	}
}

func loadFxRates() {
	codeToFxRate = make(map[string][]FxRate)

	totalCurrencies := len(cbrCodeToCurrency)
	var wg sync.WaitGroup
	wg.Add(totalCurrencies)

	for cbrCode := range cbrCodeToCurrency {
		lastDate := findLastDate(cbrCode)
		go loadFxRate(cbrCode, lastDate, &wg)
	}

	wg.Wait()
}

func findLastDate(cbrCode string) time.Time {
	lastDate := GetLastDate(cbrCode)
	if lastDate.Equal(*new(time.Time)) {
		lastDate = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	return lastDate
}

func loadFxRate(cbrCode string, lastDate time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Loading FX rates for %s from %s\n", cbrCodeToCurrency[cbrCode].CodeEng, lastDate)

	bytes := load(fxRatesRequestBody(cbrCode, lastDate))
	fxRates := FxRates(&bytes)

	log.Printf("Loaded %d FX rates for %s\n", len(fxRates), cbrCodeToCurrency[cbrCode].CodeEng)

	for _, fxRate := range fxRates {
		// Avoid merged codes for single currency
		// Example: request for R01670 (TJS) returns two codes: R01670 and R01670B
		// We store only original R01670
		fxRate.CbrCode = cbrCode
		AddFxRate(fxRate)
	}
}

func fxRatesRequestBody(cbrCode string, lastDate time.Time) *strings.Reader {
	body := `<?xml version="1.0" encoding="utf-8"?>
		<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
		<soap12:Body>
			<GetCursDynamicXML xmlns="http://web.cbr.ru/">
				<FromDate>%s</FromDate>
				<ToDate>%s</ToDate>
				<ValutaCode>%s</ValutaCode>
			</GetCursDynamicXML>
		</soap12:Body>
		</soap12:Envelope>
	`
	now := time.Now()
	from := lastDate.Format(time.RFC3339)
	to := now.Format(time.RFC3339)
	body = fmt.Sprintf(body, from, to, cbrCode)
	body = strings.TrimSpace(body)
	return strings.NewReader(body)
}
