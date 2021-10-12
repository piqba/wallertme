package web3

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type Balance struct {
	JsonRPC string `json:"jsonrpc,omitempty"`
	ID      int    `json:"id,omitempty"`
	Result  string `json:"result,omitempty"`
}

func (v *Balance) ToJSON() string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func (b *Balance) ParseDataFromHexToInt64() {
	result, err := ConvHexToDec(b.Result)
	b.Result = result
	if err != nil {
		log.Println(err)
	}
}

func (c *apiClient) Balance(ctx context.Context, payload PayloadReq) (blc Balance, err error) {

	requestUrl, err := url.Parse(c.server)
	if err != nil {
		return blc, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload.ToReader())
	if err != nil {
		return blc, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		return blc, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&blc); err != nil {
		return blc, err
	}
	blc.ParseDataFromHexToInt64()
	return blc, nil
}
