package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/max-samoylov/go-examples/cbr-loader/app/store"
)

var cbrCodeToCurrency map[string]Currency
var codeToFxRate map[string][]FxRate

func main() {
	store.Init()
	loadCurrencies()
	loadFxRates()
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

func loadFxRates() {
	codeToFxRate = make(map[string][]FxRate)
	for cbrCode := range cbrCodeToCurrency {
		loadFxRate(cbrCode)
	}
}

func loadFxRate(cbrCode string) {
	bytes := load(fxRatesRequestBody(cbrCode))
	fxRates := FxRates(&bytes)

	for _, fxRate := range fxRates {
		code := cbrCodeToCurrency[cbrCode].CodeEng
		codeToFxRate[code] = append(codeToFxRate[code], fxRate)
		fmt.Println(fxRate.ToString(code))
	}
}

func fxRatesRequestBody(cbrCode string) *strings.Reader {
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
	body = fmt.Sprintf(body, from, to, cbrCode)
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
