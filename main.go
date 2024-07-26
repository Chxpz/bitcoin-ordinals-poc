package main

import "github.com/chxpz/bitcoin-ordinals-poc/getincription"

func main() {
	txHash := "2333b105bced19fcb5fef2f73539fdebb3416374c9ea0ed23277bb7d1e3ff4a8"
	getincription.GetInscriptionDetails(txHash)
}
