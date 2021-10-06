package web3

import "encoding/json"

func UnmarshalTransaction(data []byte) (Transaction, error) {
	var r Transaction
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Transaction) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

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
