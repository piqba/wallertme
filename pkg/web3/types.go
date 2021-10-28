package web3

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

var (
	CardanoTestNet = "CardanoTestNet"
	CardanoMainNet = "CardanoMainNet"

	SolanaDevNet = "SolanaDevNet"

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
