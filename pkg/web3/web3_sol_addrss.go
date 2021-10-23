package web3

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type LastTxByAddr struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Result  []ResultTxByAddressSOL `json:"result"`
	ID      int64                  `json:"id"`
}

type ResultTxByAddressSOL struct {
	BlockTime          int64       `json:"blockTime"`
	ConfirmationStatus string      `json:"confirmationStatus"`
	Err                interface{} `json:"err"`
	Memo               interface{} `json:"memo"`
	Signature          string      `json:"signature"`
	Slot               int64       `json:"slot"`
}

func (r *ResultTxByAddressSOL) ToJSON() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func (c *apiClient) LastTxByAddress(ctx context.Context, payload PayloadReqJSONRPC) (lastTx LastTxByAddr, err error) {

	requestUrl, err := url.Parse(c.server)
	if err != nil {
		return lastTx, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload.ToReader())
	if err != nil {
		return lastTx, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		return lastTx, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&lastTx); err != nil {
		return lastTx, err
	}
	return lastTx, nil
}
