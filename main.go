package main

import (
	"fmt"
	"net/http"

	"github.com/mohithchintu/finalyear_project_service/handlers"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.TestHandler)
	mux.HandleFunc("/devices", handlers.CreateDeviceHandler)
	mux.HandleFunc("/sss", handlers.SSSHandler)
	mux.HandleFunc("/authenticate", handlers.ProcessDevicesHandler)

	handler := enableCORS(mux)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", handler)
}
