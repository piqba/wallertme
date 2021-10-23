package web3

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	// localhost rpc
	GanacheDevNet = "http://127.0.0.1:8545"
	// Cardano ecosystem
	CardanoTestNet = "https://explorer-api.testnet.dandelion.link"
	CardanoMainNet = "https://explorer-api.mainnet.dandelion.link"

	// Solana Ecosystem
	SolanaDevNet = "https://api.devnet.solana.com"
)

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
