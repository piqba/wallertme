package web3

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

var (
	// CardanoTestNet ...
	CardanoTestNet = "CardanoTestNet"
	// CardanoMainNet ...
	CardanoMainNet = "CardanoMainNet"
	// SolanaDevNet ...
	SolanaDevNet = "SolanaDevNet"

	// NetworkMap is a map that contain network information from this blockchain supports
	NetworkMap = map[string]BlockchainNet{
		CardanoTestNet: {
			NameNet:     "CardanoTestNet",
			ApiURL:      "https://explorer-api.testnet.dandelion.link",
			ExplorerURL: "https://explorer.cardano-testnet.iohkdev.io/en/transaction?id=%s",
		},
		SolanaDevNet: {
			NameNet:     "SolanaDevNet",
			ApiURL:      "https://api.devnet.solana.com",
			ExplorerURL: "https://explorer.solana.com/tx/%s?cluster=devnet",
		},
	}
)

// BlockchainNet struct that define behaivor of our logic
type BlockchainNet struct {
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

// PayloadReqJSONRPC it`s an object that define how the client can be make a simple request with JSON - RPC standard
type PayloadReqJSONRPC struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int64         `json:"id"`
}

// ToReader convert string result to reader interfaces
func (p *PayloadReqJSONRPC) ToReader() *strings.Reader {
	byte, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	return strings.NewReader(string(byte))
}
