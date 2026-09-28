package tcmb

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const sampleTCMBXML = `<?xml version="1.0" encoding="UTF-8"?>
<Tarih_Date Tarih="28.09.2026" Date="09/28/2026" Bulten_No="2026/182">
	<Currency CrossOrder="0" Kod="USD" CurrencyCode="USD">
		<Unit>1</Unit>
		<Isim>ABD DOLARI</Isim>
		<CurrencyName>US DOLLAR</CurrencyName>
		<ForexBuying>48.9008</ForexBuying>
		<ForexSelling>48.9889</ForexSelling>
		<BanknoteBuying>48.8665</BanknoteBuying>
		<BanknoteSelling>49.0623</BanknoteSelling>
		<CrossRateUSD/>
		<CrossRateOther/>
	</Currency>
	<Currency CrossOrder="1" Kod="AUD" CurrencyCode="AUD">
		<Unit>1</Unit>
		<Isim>AVUSTRALYA DOLARI</Isim>
		<CurrencyName>AUSTRALIAN DOLLAR</CurrencyName>
		<ForexBuying>34.2317</ForexBuying>
		<ForexSelling>34.4550</ForexSelling>
		<BanknoteBuying>34.0743</BanknoteBuying>
		<BanknoteSelling>34.6617</BanknoteSelling>
		<CrossRateUSD>1.4252</CrossRateUSD>
		<CrossRateOther/>
	</Currency>
	<Currency CrossOrder="9" Kod="EUR" CurrencyCode="EUR">
		<Unit>1</Unit>
		<Isim>EURO</Isim>
		<CurrencyName>EURO</CurrencyName>
		<ForexBuying>52.1234</ForexBuying>
		<ForexSelling>52.2345</ForexSelling>
		<BanknoteBuying>52.0123</BanknoteBuying>
		<BanknoteSelling>52.3456</BanknoteSelling>
		<CrossRateUSD/>
		<CrossRateOther>1.0658</CrossRateOther>
	</Currency>
	<Currency CrossOrder="15" Kod="JPY" CurrencyCode="JPY">
		<Unit>100</Unit>
		<Isim>JAPON YENİ</Isim>
		<CurrencyName>JAPANESE YEN</CurrencyName>
		<ForexBuying>32.5012</ForexBuying>
		<ForexSelling>32.7180</ForexSelling>
		<BanknoteBuying>32.3812</BanknoteBuying>
		<BanknoteSelling>32.8450</BanknoteSelling>
		<CrossRateUSD>150.45</CrossRateUSD>
		<CrossRateOther/>
	</Currency>
</Tarih_Date>`

func TestParseXML_Valid(t *testing.T) {
	testDate := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	currDay, err := ParseXML([]byte(sampleTCMBXML), testDate)
	if err != nil {
		t.Fatalf("unexpected error parsing valid XML: %v", err)
	}

	if currDay.ID != "20260928" {
		t.Errorf("expected ID '20260928', got %q", currDay.ID)
	}
	if !currDay.Date.Equal(testDate) {
		t.Errorf("expected Date %v, got %v", testDate, currDay.Date)
	}
	if currDay.DayNo != "2026/182" {
		t.Errorf("expected DayNo '2026/182', got %q", currDay.DayNo)
	}
	if len(currDay.Currencies) != 4 {
		t.Fatalf("expected 4 currencies, got %d", len(currDay.Currencies))
	}

	usd := currDay.Currencies[0]
	if usd.Code != "USD" {
		t.Errorf("expected USD code, got %q", usd.Code)
	}
	if usd.CrossOrder != 0 {
		t.Errorf("expected crossOrder 0, got %d", usd.CrossOrder)
	}
	if usd.Unit != 1 {
		t.Errorf("expected unit 1, got %d", usd.Unit)
	}
	if usd.CurrencyNameTR != "ABD DOLARI" {
		t.Errorf("expected 'ABD DOLARI', got %q", usd.CurrencyNameTR)
	}
	if usd.CurrencyName != "US DOLLAR" {
		t.Errorf("expected 'US DOLLAR', got %q", usd.CurrencyName)
	}
	if usd.ForexBuying != 48.9008 {
		t.Errorf("expected ForexBuying 48.9008, got %f", usd.ForexBuying)
	}
	if usd.ForexSelling != 48.9889 {
		t.Errorf("expected ForexSelling 48.9889, got %f", usd.ForexSelling)
	}
	if usd.BanknoteBuying != 48.8665 {
		t.Errorf("expected BanknoteBuying 48.8665, got %f", usd.BanknoteBuying)
	}
	if usd.BanknoteSelling != 49.0623 {
		t.Errorf("expected BanknoteSelling 49.0623, got %f", usd.BanknoteSelling)
	}
	if usd.CrossRateUSD != 0 {
		t.Errorf("expected empty CrossRateUSD to be 0, got %f", usd.CrossRateUSD)
	}
	if usd.CrossRateOther != 0 {
		t.Errorf("expected empty CrossRateOther to be 0, got %f", usd.CrossRateOther)
	}

	jpy := currDay.Currencies[3]
	if jpy.Code != "JPY" {
		t.Errorf("expected JPY code, got %q", jpy.Code)
	}
	if jpy.Unit != 100 {
		t.Errorf("expected unit 100, got %d", jpy.Unit)
	}
	if jpy.CrossRateUSD != 150.45 {
		t.Errorf("expected CrossRateUSD 150.45, got %f", jpy.CrossRateUSD)
	}

	eur := currDay.Currencies[2]
	if eur.CrossRateOther != 1.0658 {
		t.Errorf("expected CrossRateOther 1.0658, got %f", eur.CrossRateOther)
	}
}

func TestParseXML_Malformed(t *testing.T) {
	testDate := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	malformedXML := []byte("<Tarih_Date><Currency>unterminated")

	currDay, err := ParseXML(malformedXML, testDate)
	if err == nil {
		t.Errorf("expected error parsing malformed XML, got nil")
	}
	if currDay != nil {
		t.Errorf("expected nil CurrencyDay on malformed XML, got %+v", currDay)
	}
}

func TestParseXML_Incomplete(t *testing.T) {
	testDate := time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC)
	incompleteXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<Tarih_Date Tarih="28.09.2026" Date="09/28/2026" Bulten_No="2026/182">
	<Currency CurrencyCode="XDR">
		<Unit>1</Unit>
		<Isim>ÖZEL ÇEKME HAKKI (SDR)</Isim>
		<CurrencyName>SPECIAL DRAWING RIGHT (SDR)</CurrencyName>
		<ForexBuying>64.1234</ForexBuying>
		<ForexSelling/>
	</Currency>
</Tarih_Date>`)

	currDay, err := ParseXML(incompleteXML, testDate)
	if err != nil {
		t.Fatalf("unexpected error parsing incomplete XML: %v", err)
	}

	if len(currDay.Currencies) != 1 {
		t.Fatalf("expected 1 currency, got %d", len(currDay.Currencies))
	}

	xdr := currDay.Currencies[0]
	if xdr.Code != "XDR" {
		t.Errorf("expected code XDR, got %q", xdr.Code)
	}
	if xdr.ForexBuying != 64.1234 {
		t.Errorf("expected ForexBuying 64.1234, got %f", xdr.ForexBuying)
	}
	if xdr.ForexSelling != 0 {
		t.Errorf("expected empty ForexSelling to be 0, got %f", xdr.ForexSelling)
	}
	if xdr.BanknoteBuying != 0 {
		t.Errorf("expected missing BanknoteBuying to be 0, got %f", xdr.BanknoteBuying)
	}
}

func TestParseFloatAndInt(t *testing.T) {
	testsFloat := []struct {
		input    string
		expected float64
	}{
		{"", 0},
		{"invalid", 0},
		{"0", 0},
		{"48.9008", 48.9008},
		{"-3.14", -3.14},
	}

	for _, tt := range testsFloat {
		if got := parseFloat(tt.input); got != tt.expected {
			t.Errorf("parseFloat(%q) = %f; want %f", tt.input, got, tt.expected)
		}
	}

	testsInt := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"abc", 0},
		{"0", 0},
		{"1", 1},
		{"100", 100},
		{"-5", -5},
	}

	for _, tt := range testsInt {
		if got := parseInt(tt.input); got != tt.expected {
			t.Errorf("parseInt(%q) = %d; want %d", tt.input, got, tt.expected)
		}
	}
}

func TestGetCurrencyDay_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sampleTCMBXML))
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	testDate := time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC)
	currDay, err := GetCurrencyDay(testDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if currDay == nil {
		t.Fatal("expected non-nil CurrencyDay")
	}
	if currDay.DayNo != "2026/182" {
		t.Errorf("expected DayNo '2026/182', got %q", currDay.DayNo)
	}
	if len(currDay.Currencies) != 4 {
		t.Errorf("expected 4 currencies, got %d", len(currDay.Currencies))
	}
}

func TestGetCurrencyDay_WeekendFallback(t *testing.T) {
	// Request on Sunday 2026-09-27 -> 404
	// Saturday 2026-09-26 -> 404
	// Friday 2026-09-25 -> 200 OK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "25092026.xml") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(sampleTCMBXML))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	sunday := time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC)
	currDay, err := GetCurrencyDay(sunday)
	if err != nil {
		t.Fatalf("unexpected error with weekend fallback: %v", err)
	}

	if currDay == nil {
		t.Fatal("expected non-nil CurrencyDay")
	}
	// ID and Date should reflect the original requested Sunday
	if currDay.ID != "20260927" {
		t.Errorf("expected ID '20260927', got %q", currDay.ID)
	}
	if !currDay.Date.Equal(sunday) {
		t.Errorf("expected Date %v, got %v", sunday, currDay.Date)
	}
	if len(currDay.Currencies) != 4 {
		t.Errorf("expected 4 currencies, got %d", len(currDay.Currencies))
	}
}

func TestGetCurrencyDay_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	testDate := time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC)
	currDay, err := GetCurrencyDay(testDate)
	if err == nil {
		t.Fatal("expected error on server 500, got nil")
	}
	if currDay != nil {
		t.Errorf("expected nil CurrencyDay on server error, got %+v", currDay)
	}
}

func TestGetCurrencyDay_ExhaustLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	testDate := time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC)
	currDay, err := GetCurrencyDay(testDate)
	if err == nil {
		t.Fatal("expected error when all 30 days return 404, got nil")
	}
	if currDay != nil {
		t.Errorf("expected nil CurrencyDay when exhausted, got %+v", currDay)
	}
}
