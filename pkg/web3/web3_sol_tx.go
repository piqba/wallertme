package web3

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type TxInfo struct {
	Jsonrpc string       `json:"jsonrpc"`
	Result  ResultTxInfo `json:"result"`
	ID      int64        `json:"id"`
}

func (r *TxInfo) ToJSON() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

type ResultTxInfo struct {
	BlockTime   int64       `json:"blockTime"`
	Meta        Meta        `json:"meta"`
	Slot        int64       `json:"slot"`
	Transaction Transaction `json:"transaction"`
}

type Meta struct {
	Err               interface{}   `json:"err"`
	Fee               int64         `json:"fee"`
	InnerInstructions []interface{} `json:"innerInstructions"`
	LogMessages       []string      `json:"logMessages"`
	PostBalances      []int64       `json:"postBalances"`
	PostTokenBalances []interface{} `json:"postTokenBalances"`
	PreBalances       []int64       `json:"preBalances"`
	PreTokenBalances  []interface{} `json:"preTokenBalances"`
	Rewards           []interface{} `json:"rewards"`
	Status            Status        `json:"status"`
}

type Status struct {
	Ok interface{} `json:"Ok"`
}

type Transaction struct {
	Message    Message  `json:"message"`
	Signatures []string `json:"signatures"`
}

type Message struct {
	AccountKeys     []string      `json:"accountKeys"`
	Header          Header        `json:"header"`
	Instructions    []Instruction `json:"instructions"`
	RecentBlockhash string        `json:"recentBlockhash"`
}

type Header struct {
	NumReadonlySignedAccounts   int64 `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts int64 `json:"numReadonlyUnsignedAccounts"`
	NumRequiredSignatures       int64 `json:"numRequiredSignatures"`
}

type Instruction struct {
	Accounts       []int64 `json:"accounts"`
	Data           string  `json:"data"`
	ProgramIDIndex int64   `json:"programIdIndex"`
}

func (c *apiClient) InfoByTx(ctx context.Context, payload PayloadReqJSONRPC) (info TxInfo, err error) {

	requestUrl, err := url.Parse(c.server)
	if err != nil {
		return info, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload.ToReader())
	if err != nil {
		return info, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		return info, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&info); err != nil {
		return info, err
	}
	return info, nil
}
