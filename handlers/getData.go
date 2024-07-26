package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chxpz/bitcoin-ordinals-poc/getincription"
)

func GetInscription(w http.ResponseWriter, r *http.Request) {
	txHash := r.URL.Query().Get("txHash")

	if txHash == "" {
		http.Error(w, "txHash is required", http.StatusBadRequest)
		return
	}

	log.Default().Printf("Getting inscription details for txHash: %s", txHash)

	ret, err := getincription.GetInscriptionDetails(txHash)
	if err != nil {
		http.Error(w, "Error retrieving data", http.StatusInternalServerError)
		log.Printf("Error getting inscription details: %v", err)
		return
	}

	// Convert the ReturnData struct to JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
