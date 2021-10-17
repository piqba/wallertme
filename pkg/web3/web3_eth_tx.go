package web3

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type Transaction struct {
	Jsonrpc string   `json:"jsonrpc"`
	ID      int64    `json:"id"`
	Result  ResultTx `json:"result"`
}

type ResultTx struct {
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

func (r *ResultTx) ToJSON() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
func (r *ResultTx) ParseDataFromHexToInt64() {
	r.BlockNumber, _ = ConvHexToDec(r.BlockNumber)
	r.Gas, _ = ConvHexToDec(r.Gas)
	r.GasPrice, _ = ConvHexToDec(r.GasPrice)
	r.Nonce, _ = ConvHexToDec(r.Nonce)
	r.TransactionIndex, _ = ConvHexToDec(r.TransactionIndex)
	r.Value, _ = ConvHexToDec(r.Value)
	r.Type, _ = ConvHexToDec(r.Type)
}

type TransactionReceipt struct {
	Jsonrpc string    `json:"jsonrpc"`
	ID      int64     `json:"id"`
	Result  ResultTxR `json:"result"`
}

type ResultTxR struct {
	BlockHash         string      `json:"blockHash"`
	BlockNumber       string      `json:"blockNumber"`
	ContractAddress   interface{} `json:"contractAddress"`
	CumulativeGasUsed string      `json:"cumulativeGasUsed"`
	EffectiveGasPrice string      `json:"effectiveGasPrice"`
	From              string      `json:"from"`
	GasUsed           string      `json:"gasUsed"`
	Logs              []Log       `json:"logs"`
	LogsBloom         string      `json:"logsBloom"`
	Status            string      `json:"status"`
	To                string      `json:"to"`
	TransactionHash   string      `json:"transactionHash"`
	TransactionIndex  string      `json:"transactionIndex"`
	Type              string      `json:"type"`
}

type Log struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

func (r *ResultTxR) ToJSON() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func (c *apiClient) TransactionByHash(ctx context.Context, payload PayloadReq) (tx Transaction, err error) {

	requestUrl, err := url.Parse(c.server)
	if err != nil {
		return Transaction{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload.ToReader())
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

func (c *apiClient) TransactionReceipt(ctx context.Context, payload PayloadReq) (tx TransactionReceipt, err error) {

	requestUrl, err := url.Parse(c.server)
	if err != nil {
		return TransactionReceipt{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload.ToReader())
	if err != nil {
		return TransactionReceipt{}, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		return TransactionReceipt{}, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&tx); err != nil {
		return TransactionReceipt{}, err
	}

	// parse result elements
	// TODO: make parse

	return tx, nil
}
