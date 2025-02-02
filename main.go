package handler

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Text string `json:"text"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/get" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("GET request received"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	case http.MethodPost:
		if r.URL.Path == "/post" {
			var msg Message
			err := json.NewDecoder(r.Body).Decode(&msg)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid request body"))
				return
			}
			response := map[string]string{"message": "POST request received with message: " + msg.Text}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Invalid request method"))
	}
}
