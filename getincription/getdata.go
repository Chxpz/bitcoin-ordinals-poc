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

type ReturnData struct {
	ContentType     string `json:"contentType"`
	TxInIndex       string `json:"txInIndex"`
	TxInOffset      string `json:"txInOffset"`
	ContentLength   string `json:"contentLength"`
	InscriptionData string `json:"inscriptionData"`
}

func GetInscriptionDetails(txHash string) (ReturnData, error) {
	var returnData ReturnData

	rawTxHex, err := _getRawTransactionHex(txHash)
	if err != nil {
		log.Fatalf("Failed to fetch raw transaction hex: %v", err)
		return returnData, err
	}

	rawTx, err := _decodeRawTransaction(rawTxHex)
	if err != nil {
		log.Fatalf("Failed to decode raw transaction: %v", err)
		return returnData, err
	}

	transactionInscriptions := parser.ParseInscriptionsFromTransaction(rawTx)
	if len(transactionInscriptions) == 0 {
		log.Info("No inscriptions found.")
		return returnData, nil
	}

	for _, ins := range transactionInscriptions {
		contentType := string(ins.Inscription.ContentType)

		if contentType == "text/plain;charset=utf-8" {
			data := ins.Inscription.ContentBody
			if data != nil {
				returnData = ReturnData{
					ContentType:     contentType,
					TxInIndex:       fmt.Sprintf("%d", ins.TxInIndex),
					TxInOffset:      fmt.Sprintf("%d", ins.TxInOffset),
					ContentLength:   fmt.Sprintf("%d", ins.Inscription.ContentLength),
					InscriptionData: string(data),
				}
				return returnData, nil
			} else {
				log.Warn("Inscription data is nil.")
			}
		}
	}

	return returnData, nil
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
