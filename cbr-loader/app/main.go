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
	Init()
	loadCurrencies()
	loadFxRates()
	Shutdown()
}

func load(body *strings.Reader) []byte {
	client := http.Client{}
	resp, err := client.Post("http://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx",
		"application/soap+xml", body)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	return bytes
}

func loadCurrencies() {
	bytes := load(currenciesRefDataRequestBody())
	currencies := Currencies(&bytes)

	cbrCodeToCurrency = make(map[string]Currency)
	for _, currency := range currencies {
		cbrCodeToCurrency[currency.CodeCbr] = currency
		AddCurrency(currency)
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
	log.Printf("Loading %s from %s\n", cbrCodeToCurrency[cbrCode].CodeEng, lastDate)

	defer wg.Done()

	bytes := load(fxRatesRequestBody(cbrCode, lastDate))
	fxRates := FxRates(&bytes)

	for _, fxRate := range fxRates {
		// code := cbrCodeToCurrency[fxRate.CbrCode].CodeEng
		// codeToFxRate[code] = append(codeToFxRate[code], fxRate)
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
	if now.After(lastDate) {
		from := lastDate.Format(time.RFC3339)
		to := now.Format(time.RFC3339)
		body = fmt.Sprintf(body, from, to, cbrCode)
		body = strings.TrimSpace(body)
		return strings.NewReader(body)
	}
	return nil
}

func currenciesRefDataRequestBody() *strings.Reader {
	body := `
		<?xml version="1.0" encoding="utf-8"?>
		<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
		<soap12:Body>
			<EnumValutesXML xmlns="http://web.cbr.ru/">
				<Seld>false</Seld>
			</EnumValutesXML>
		</soap12:Body>
		</soap12:Envelope>
	`
	body = strings.TrimSpace(body)
	return strings.NewReader(body)
}
