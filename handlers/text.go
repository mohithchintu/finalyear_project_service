package handlers

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Text string `json:"Message"`
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	response := Message{Text: "Hello, Welcone to final year services"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
