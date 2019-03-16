package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var cbrCodeToCurrency map[string]Currency
var codeToFxRate map[string]FxRate

func main() {
	loadCurrencies()
	loadFxRate()
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
	}
}

func loadFxRate() {
	bytes := load(fxRatesRequestBody())
	fxRates := FxRates(&bytes)

	codeToFxRate = make(map[string]FxRate)
	for _, fxRate := range fxRates {
		codeToFxRate["R01235"] = fxRate
		fmt.Println(fxRate.ToString())
	}
}

func fxRatesRequestBody() *strings.Reader {
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
	minus3days, _ := time.ParseDuration("-72h")
	from := time.Now().Add(minus3days).Format(time.RFC3339)
	to := time.Now().Format(time.RFC3339)
	body = fmt.Sprintf(body, from, to, "R01235")
	body = strings.TrimSpace(body)
	return strings.NewReader(body)
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
