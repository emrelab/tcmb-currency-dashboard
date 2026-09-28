# TCMB Currency Dashboard

A full-stack currency dashboard built with **Go** and **React**, using official exchange-rate data published by the **Central Bank of the Republic of Türkiye (TCMB)**.

The application fetches currency data from TCMB's XML endpoint, processes it with a Go backend, and exposes the data through a REST API consumed by a React frontend.

![TCMB Currency Dashboard](./assets/ui.jpeg)

## Features

* Live TCMB exchange-rate data
* REST API built with Go
* XML parsing with Go's `encoding/xml`
* React + Vite frontend
* Currency search and filtering
* USD, EUR and GBP highlighted rates
* Clean dashboard interface
* Health-check endpoint

## Tech Stack

### Backend

* Go
* `net/http`
* `encoding/xml`
* REST API

### Frontend

* React
* Vite
* Axios
* CSS

## Architecture

```text
TCMB XML API
     │
     ▼
┌─────────────────┐
│   Go Backend    │
│                 │
│ XML Parsing     │
│ Data Processing │
│ REST API        │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ React Frontend  │
│                 │
│ Dashboard       │
│ Search          │
│ Filtering       │
└─────────────────┘
```

## API

### Health Check

```http
GET /api/health
```

**Response (`200 OK`):**

```text
OK
```

### Today's Currency Rates

```http
GET /api/currencies/today
```

**Response (`200 OK` - `application/json`):**

```json
{
  "id": "20260929",
  "date": "2026-09-29T00:00:00Z",
  "dayNo": "2026/182",
  "currencies": [
    {
      "code": "USD",
      "crossOrder": 0,
      "unit": 1,
      "currencyNameTr": "ABD DOLARI",
      "currencyName": "US DOLLAR",
      "forexBuying": 48.9008,
      "forexSelling": 48.9889,
      "banknoteBuying": 48.8665,
      "banknoteSelling": 49.0623,
      "crossRateUsd": 0,
      "crossRateOther": 0
    },
    {
      "code": "EUR",
      "crossOrder": 9,
      "unit": 1,
      "currencyNameTr": "EURO",
      "currencyName": "EURO",
      "forexBuying": 55.6307,
      "forexSelling": 55.731,
      "banknoteBuying": 55.5918,
      "banknoteSelling": 55.8146,
      "crossRateUsd": 0,
      "crossRateOther": 1.1376
    }
  ]
}
```

## Getting Started

### Backend

```bash
go run .
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Data Source

Currency data is provided by the official **Türkiye Cumhuriyet Merkez Bankası (TCMB)**.

https://www.tcmb.gov.tr/kurlar/

## Project Structure

```text
tcmb-currency-dashboard/
├── internal/
│   ├── tcmb/       # TCMB data fetching and XML parsing
│   └── http/       # HTTP handlers and REST API
│
├── frontend/       # React + Vite application
├── assets/         # Screenshots and project assets
├── LICENSE         # MIT License
└── README.md
```

## Purpose

This project was built to practice working with real-world financial data and to explore:

* Go HTTP servers
* XML parsing
* REST API design
* React frontend development
* Backend/frontend integration
* External API consumption

## Status

This project is actively open for improvements and experimentation.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
