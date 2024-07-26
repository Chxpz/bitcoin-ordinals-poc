package main

import (
	"log"
	"net/http"

	"github.com/chxpz/bitcoin-ordinals-poc/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// txHash := "2333b105bced19fcb5fef2f73539fdebb3416374c9ea0ed23277bb7d1e3ff4a8"
	// ret, err := getincription.GetInscriptionDetails(txHash)

	// if err != nil {
	// 	log.Fatalf("Error to get data", err)
	// }

	// log.Println(ret)

	r := mux.NewRouter()

	r.HandleFunc("/getinscription", handlers.GetInscription).Methods("GET")
	// r.HandleFunc("/createinscription", handlers.CreateInscription).Methods("POST")

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
