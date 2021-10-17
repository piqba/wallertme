package web3

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type Block struct {
	ID      int64     `json:"id"`
	Jsonrpc string    `json:"jsonrpc"`
	Result  ResultBlk `json:"result"`
}

func (v *Block) ToJSON() string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

type ResultBlk struct {
	Number           string               `json:"number"`
	Hash             string               `json:"hash"`
	ParentHash       string               `json:"parentHash"`
	MixHash          string               `json:"mixHash"`
	Nonce            string               `json:"nonce"`
	Sha3Uncles       string               `json:"sha3Uncles"`
	LogsBloom        string               `json:"logsBloom"`
	TransactionsRoot string               `json:"transactionsRoot"`
	StateRoot        string               `json:"stateRoot"`
	ReceiptsRoot     string               `json:"receiptsRoot"`
	Miner            string               `json:"miner"`
	Difficulty       string               `json:"difficulty"`
	TotalDifficulty  string               `json:"totalDifficulty"`
	ExtraData        string               `json:"extraData"`
	Size             string               `json:"size"`
	GasLimit         string               `json:"gasLimit"`
	GasUsed          string               `json:"gasUsed"`
	Timestamp        string               `json:"timestamp"`
	Transactions     []TransactionOnBlock `json:"transactions"`
	Uncles           []interface{}        `json:"uncles"`
}

type TransactionOnBlock struct {
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Input            string `json:"input"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

func (v *ResultBlk) ToJSON() string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func (c *apiClient) BlockByNumberETH(ctx context.Context, payload PayloadReq) (blk Block, err error) {

	requestUrl, err := url.Parse(c.server)
	if err != nil {
		return blk, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload.ToReader())
	if err != nil {
		return blk, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		return blk, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&blk); err != nil {
		return blk, err
	}

	// parse result elements
	// TODO: parsed elements

	return blk, nil
}
