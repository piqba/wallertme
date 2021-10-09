package web3

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	LATEST   = "latest"
	PENDING  = "pending"
	EARLIEST = "earliest"
)

type FilterRequest struct {
	FromBlock string   `json:"fromBlock"`
	ToBlock   string   `json:"toBlock"`
	Address   []string `json:"address"`
	Topics    []string `json:"topics"`
}
type FilterResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Result  string `json:"result"`
}

func New(addrs, topics []string) FilterRequest {
	return FilterRequest{
		FromBlock: LATEST,
		ToBlock:   LATEST,
		Address:   addrs,
		Topics:    topics,
	}
}

func (c *apiClient) Filter(ctx context.Context, payload PayloadReq) (filterRes FilterResponse, err error) {

	requestUrl, err := url.Parse(c.server)
	if err != nil {
		return FilterResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload.ToReader())
	if err != nil {
		return FilterResponse{}, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		return FilterResponse{}, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&filterRes); err != nil {
		return FilterResponse{}, err
	}

	return filterRes, nil
}
