package sendinscription

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

const apiURL = "https://testnet-api.ordinalsbot.com/order"

type InscriptionOrder struct {
	Files          []File `json:"files"`
	ReceiveAddress string `json:"receiveAddress"`
	Fee            int    `json:"fee"`
	LowPostage     bool   `json:"lowPostage"`
}

type File struct {
	Name    string `json:"name"`
	Size    int    `json:"size"`
	DataURL string `json:"dataURL"`
}

type InscriptionResponse struct {
	Status string          `json:"status"`
	Charge Charge          `json:"charge"`
	Error  json.RawMessage `json:"error,omitempty"`
}

type Charge struct {
	ID               string           `json:"id"`
	Address          string           `json:"address"`
	Amount           int64            `json:"amount"`
	LightningInvoice LightningInvoice `json:"lightning_invoice"`
	CreatedAt        int64            `json:"created_at"`
}

type LightningInvoice struct {
	ExpiresAt int64  `json:"expires_at"`
	PayReq    string `json:"payreq"`
}

func Run() {

	data := "Hello, Ordinals!"
	dataURL := "data:plain/text;base64," + base64.StdEncoding.EncodeToString([]byte(data))
	files := []File{
		{
			Name:    "my-text-inscription-file.txt",
			Size:    len(data),
			DataURL: dataURL,
		},
	}

	address := "tb1qe5mfmh5p8z355zwmnq6j59r8qjseecmcf8dkvd" // Substitua por um endereço de teste válido
	fee := 2                                                // Taxa de transação em satoshis por byte
	lowPostage := true                                      // Usar postagem baixa para economizar custos

	payload := InscriptionOrder{
		Files:          files,
		ReceiveAddress: address,
		Fee:            fee,
		LowPostage:     lowPostage,
	}

	resp, err := sendInscription(payload)
	if err != nil {
		fmt.Println("Error sending inscription:", err)
		return
	}

	fmt.Printf("Inscription created, charge details: %+v\n", resp.Charge)
}

func sendInscription(payload InscriptionOrder) (*InscriptionResponse, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error encoding payload: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var response InscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if response.Status != "ok" {
		var errorMsg string
		if len(response.Error) > 0 {
			errorMsg = string(response.Error)
		} else {
			errorMsg = "unknown error"
		}
		return nil, fmt.Errorf("error in API response: %s", errorMsg)
	}

	return &response, nil
}
