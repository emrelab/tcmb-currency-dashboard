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

### Today's Currency Rates

```http
GET /api/currencies/today
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
TCBMCurrency/
├── internal/
│   ├── tcmb/       # TCMB data fetching and XML parsing
│   └── http/       # HTTP handlers and REST API
│
├── frontend/       # React + Vite application
├── assets/         # Screenshots and project assets
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
