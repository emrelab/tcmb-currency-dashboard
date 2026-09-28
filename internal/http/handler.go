package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/emrelab/tcmb-currency-dashboard/internal/tcmb"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func TodayCurrenciesHandler(w http.ResponseWriter, r *http.Request) {
	data, err := tcmb.GetCurrencyDay(time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
