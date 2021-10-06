package web3

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Transaction struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Result  Result `json:"result"`
}

type Result struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

func (r *Result) ToJSON() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
func (r *Result) ParseDataFromHexToInt64() {
	r.BlockNumber, _ = convHexToDec(r.BlockNumber)
	r.Gas, _ = convHexToDec(r.Gas)
	r.GasPrice, _ = convHexToDec(r.GasPrice)
	r.Nonce, _ = convHexToDec(r.Nonce)
	r.TransactionIndex, _ = convHexToDec(r.TransactionIndex)
	r.Value, _ = convHexToDec(r.Value)
	r.Type, _ = convHexToDec(r.Type)
}

func (c *apiClient) TransactionByHash(ctx context.Context, hash string) (tx Transaction, err error) {
	payload := strings.NewReader(
		fmt.Sprintf(
			`{
				"jsonrpc":"2.0",
				"method":"eth_getTransactionByHash",
				"params":[
					"%s"
				],
				"id":1
			}`,
			hash,
		),
	)
	requestUrl, err := url.Parse(c.server)
	if err != nil {
		return Transaction{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload)
	if err != nil {
		return Transaction{}, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		return Transaction{}, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&tx); err != nil {
		return Transaction{}, err
	}

	// parse result elements
	tx.Result.ParseDataFromHexToInt64()

	return tx, nil
}
