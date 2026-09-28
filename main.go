package main

import (
	"fmt"
	"log"
	"net/http"

	api "github.com/emrelab/tcmb-currency-dashboard/internal/http"
)

func withCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", api.HealthHandler)
	mux.HandleFunc("/api/currencies/today", api.TodayCurrenciesHandler)

	fmt.Println("Server running on :8080")

	err := http.ListenAndServe(":8080", withCORS(mux))
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
