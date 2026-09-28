package tcmb

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

// BaseURL is the root endpoint for fetching TCMB currency XML files.
var BaseURL = "http://www.tcmb.gov.tr/kurlar"

type CurrencyDay struct {
	ID         string     `json:"id"`
	Date       time.Time  `json:"date"`
	DayNo      string     `json:"dayNo"`
	Currencies []Currency `json:"currencies"`
}

type Currency struct {
	Code            string  `json:"code"`
	CrossOrder      int     `json:"crossOrder"`
	Unit            int     `json:"unit"`
	CurrencyNameTR  string  `json:"currencyNameTr"`
	CurrencyName    string  `json:"currencyName"`
	ForexBuying     float64 `json:"forexBuying"`
	ForexSelling    float64 `json:"forexSelling"`
	BanknoteBuying  float64 `json:"banknoteBuying"`
	BanknoteSelling float64 `json:"banknoteSelling"`
	CrossRateUSD    float64 `json:"crossRateUsd"`
	CrossRateOther  float64 `json:"crossRateOther"`
}

type tarihDate struct {
	XMLName  xml.Name   `xml:"Tarih_Date"`
	Tarih    string     `xml:"Tarih,attr"`
	Date     string     `xml:"Date,attr"`
	BultenNo string     `xml:"Bulten_No,attr"`
	Currency []xmlMoney `xml:"Currency"`
}

type xmlMoney struct {
	Kod             string `xml:"Kod,attr"`
	CrossOrder      string `xml:"CrossOrder,attr"`
	CurrencyKod     string `xml:"CurrencyCode,attr"`
	Unit            string `xml:"Unit"`
	Isim            string `xml:"Isim"`
	CurrencyName    string `xml:"CurrencyName"`
	ForexBuying     string `xml:"ForexBuying"`
	ForexSelling    string `xml:"ForexSelling"`
	BanknoteBuying  string `xml:"BanknoteBuying"`
	BanknoteSelling string `xml:"BanknoteSelling"`
	CrossRateUSD    string `xml:"CrossRateUSD"`
	CrossRateOther  string `xml:"CrossRateOther"`
}

func GetCurrencyDay(date time.Time) (*CurrencyDay, error) {
	originalDate := date

	for i := 0; i < 30; i++ {
		currDay, err := fetchCurrencyDay(date, originalDate)
		if err != nil {
			return nil, err
		}

		if currDay != nil {
			return currDay, nil
		}

		date = date.AddDate(0, 0, -1)
	}

	return nil, fmt.Errorf("no currency data found within 30 days")
}

func fetchCurrencyDay(date time.Time, originalDate time.Time) (*CurrencyDay, error) {
	url := fmt.Sprintf("%s/%s/%s.xml", BaseURL, date.Format("200601"), date.Format("02012006"))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := ParseXMLReader(resp.Body, originalDate)
	if err != nil {
		log.Printf("xml decode error: %v", err)
		return nil, err
	}

	return result, nil
}

// ParseXML parses raw TCMB XML bytes into a CurrencyDay struct.
func ParseXML(data []byte, date time.Time) (*CurrencyDay, error) {
	return ParseXMLReader(bytes.NewReader(data), date)
}

// ParseXMLReader parses TCMB XML from an io.Reader into a CurrencyDay struct.
func ParseXMLReader(r io.Reader, date time.Time) (*CurrencyDay, error) {
	var parsed tarihDate
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&parsed); err != nil {
		return nil, err
	}

	result := &CurrencyDay{
		ID:         date.Format("20060102"),
		Date:       date,
		DayNo:      parsed.BultenNo,
		Currencies: make([]Currency, len(parsed.Currency)),
	}

	for i, curr := range parsed.Currency {
		result.Currencies[i] = Currency{
			Code:            curr.CurrencyKod,
			CrossOrder:      parseInt(curr.CrossOrder),
			Unit:            parseInt(curr.Unit),
			CurrencyNameTR:  curr.Isim,
			CurrencyName:    curr.CurrencyName,
			ForexBuying:     parseFloat(curr.ForexBuying),
			ForexSelling:    parseFloat(curr.ForexSelling),
			BanknoteBuying:  parseFloat(curr.BanknoteBuying),
			BanknoteSelling: parseFloat(curr.BanknoteSelling),
			CrossRateUSD:    parseFloat(curr.CrossRateUSD),
			CrossRateOther:  parseFloat(curr.CrossRateOther),
		}
	}

	return result, nil
}

func parseFloat(value string) float64 {
	if value == "" {
		return 0
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}

	return parsed
}

func parseInt(value string) int {
	if value == "" {
		return 0
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return parsed
}
