package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	client := http.Client{}
	resp, err := client.Post("http://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx", "application/soap+xml", getRequestCurrencyCodes())
	// resp, err := client.Post("http://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx", "application/soap+xml", getRequestCurrencyRates())

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bytes))
}

func getRequestCurrencyRates() *strings.Reader {
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
	body = fmt.Sprintf(body, "2019-03-15T00:00:00.000", "2019-03-15T00:00:00.000", "R01235")
	body = strings.TrimSpace(body)
	return strings.NewReader(body)
}

func getRequestCurrencyCodes() *strings.Reader {
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
