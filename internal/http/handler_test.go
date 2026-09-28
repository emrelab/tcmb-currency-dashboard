package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emrelab/tcmb-currency-dashboard/internal/tcmb"
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
</Tarih_Date>`

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	HealthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expectedBody := "OK"
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestTodayCurrenciesHandler_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sampleTCMBXML))
	}))
	defer server.Close()

	originalBaseURL := tcmb.BaseURL
	tcmb.BaseURL = server.URL
	defer func() { tcmb.BaseURL = originalBaseURL }()

	req := httptest.NewRequest(http.MethodGet, "/api/currencies/today", nil)
	rec := httptest.NewRecorder()

	TodayCurrenciesHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}

	var data tcmb.CurrencyDay
	if err := json.NewDecoder(rec.Body).Decode(&data); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if data.DayNo != "2026/182" {
		t.Errorf("expected DayNo '2026/182', got %q", data.DayNo)
	}
	if len(data.Currencies) != 1 {
		t.Fatalf("expected 1 currency, got %d", len(data.Currencies))
	}
	if data.Currencies[0].Code != "USD" {
		t.Errorf("expected currency code 'USD', got %q", data.Currencies[0].Code)
	}
}

func TestTodayCurrenciesHandler_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	originalBaseURL := tcmb.BaseURL
	tcmb.BaseURL = server.URL
	defer func() { tcmb.BaseURL = originalBaseURL }()

	req := httptest.NewRequest(http.MethodGet, "/api/currencies/today", nil)
	rec := httptest.NewRecorder()

	TodayCurrenciesHandler(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
