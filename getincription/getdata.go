package getincription

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/wire"
	log "github.com/sirupsen/logrus"
)

func GetData(txHash string) {

	rawTxHex, err := _getRawTransactionHex(txHash)
	if err != nil {
		log.Fatalf("Failed to fetch raw transaction hex: %v", err)
	}

	rawTx, err := _decodeRawTransaction(rawTxHex)
	if err != nil {
		log.Fatalf("Failed to decode raw transaction: %v", err)
	}

	transactionInscriptions := parser.ParseInscriptionsFromTransaction(rawTx)
	if len(transactionInscriptions) == 0 {
		log.Info("No inscriptions found.")
	} else {
		for _, ins := range transactionInscriptions {
			contentType := string(ins.Inscription.ContentType)
			log.Infof("Inscription found at index: %d, offset: %d, type: %s, length: %d",
				ins.TxInIndex, ins.TxInOffset, contentType, ins.Inscription.ContentLength)

			// Display the content of the inscription
			if contentType == "text/plain;charset=utf-8" {
				// Attempt to access the content field
				data := ins.Inscription.ContentBody
				if data != nil {
					fmt.Printf("Inscription Content: %s\n", string(data))
				} else {
					log.Warn("Inscription data is nil.")
				}
			}
		}
	}
}

func _getRawTransactionHex(txid string) (string, error) {
	url := fmt.Sprintf("https://mempool.space/testnet/api/tx/%s/hex", txid)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	rawHex, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(rawHex), nil
}

func _decodeRawTransaction(rawTxHex string) (*wire.MsgTx, error) {
	rawTx := wire.NewMsgTx(wire.TxVersion)
	txBytes, err := hex.DecodeString(rawTxHex)
	if err != nil {
		return nil, fmt.Errorf("error decoding raw transaction hex: %v", err)
	}

	err = rawTx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		return nil, fmt.Errorf("error deserializing transaction: %v", err)
	}

	return rawTx, nil
}
