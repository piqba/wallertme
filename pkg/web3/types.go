package web3

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	PolygonMainNet    = "https://polygon-rpc.com/"
	PolygonMainNetWSS = "wss://ws-matic-mainnet.chainstacklabs.com"
	PolygonTestNet    = "https://matic-mumbai.chainstacklabs.com/"
	GanacheDevNet     = "http://127.0.0.1:8545"
)
const (
	TRANSFER   = "Transfer"
	DEPOSIT    = "Deposit"
	WITHDRAWAL = "Withdrawal"
	APROVAL    = "Approval"
)

var (
	Topics = []Topic{
		{
			Name:  TRANSFER,
			Value: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		}, {
			Name:  DEPOSIT,
			Value: "0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c",
		}, {
			Name:  WITHDRAWAL,
			Value: "0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65",
		}, {
			Name:  APROVAL,
			Value: "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
		},
	}
)

// APIError is used to describe errors from the API.
type APIError struct {
	Response interface{}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error, %+v", e.Response)
}

// Topic
type Topic struct {
	Name  string
	Value string
}

type PayloadReq struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int64         `json:"id"`
}

func (p *PayloadReq) ToReader() *strings.Reader {
	byte, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	return strings.NewReader(string(byte))
}
