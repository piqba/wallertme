package web3

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
// localhost rpc
// GanacheDevNet = "http://127.0.0.1:8545"
// // Cardano ecosystem
// CardanoTestNet = "https://explorer-api.testnet.dandelion.link"
// CardanoMainNet = "https://explorer-api.mainnet.dandelion.link"

// // Solana Ecosystem
// SolanaDevNet = "https://api.devnet.solana.com"
)

var (
	CardanoTestNet = BlochainNet{
		NameNet:     "CardanoTestNet",
		ApiURL:      "https://explorer-api.testnet.dandelion.link",
		ExplorerURL: "https://explorer.cardano-testnet.iohkdev.io/en/transaction?id=%s",
	}

	SolanaDevNet = BlochainNet{
		NameNet:     "SolanaDevNet",
		ApiURL:      "https://api.devnet.solana.com",
		ExplorerURL: "https://explorer.solana.com/tx/%s?cluster=devnet",
	}
)

type BlochainNet struct {
	NameNet     string `json:"name_net,omitempty"`
	ApiURL      string `json:"api_url,omitempty"`
	ExplorerURL string `json:"explorer_url,omitempty"`
}

// APIError is used to describe errors from the API.
type APIError struct {
	Response interface{}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error, %+v", e.Response)
}

type PayloadReqJSONRPC struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int64         `json:"id"`
}

func (p *PayloadReqJSONRPC) ToReader() *strings.Reader {
	byte, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	return strings.NewReader(string(byte))
}
