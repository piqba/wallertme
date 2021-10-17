package web3

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type VersionETH struct {
	JsonRPC string `json:"jsonrpc,omitempty"`
	ID      int    `json:"id,omitempty"`
	Result  string `json:"result,omitempty"`
}

func (v *VersionETH) ToJSON() string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

type ListeningETH struct {
	JsonRPC string `json:"jsonrpc,omitempty"`
	ID      int    `json:"id,omitempty"`
	Result  string `json:"result,omitempty"`
}

func (l *ListeningETH) ToJSON() string {
	bytes, err := json.Marshal(l)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func (c *apiClient) VersionETH(ctx context.Context) (version VersionETH, err error) {
	payload := strings.NewReader(`{
		"jsonrpc":"2.0",
		"method":"net_version",
		"params":[],
		"id":67
	}`)
	requestUrl, err := url.Parse(c.server)
	if err != nil {
		return VersionETH{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload)
	if err != nil {
		return VersionETH{}, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		return VersionETH{}, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&version); err != nil {
		return VersionETH{}, err
	}
	return version, nil
}
func (c *apiClient) ListeningETH(ctx context.Context) (listening ListeningETH, err error) {
	payload := strings.NewReader(`{
		"jsonrpc":"2.0",
		"method":"net_listening",
		"params":[],
		"id":67
	}`)
	requestUrl, err := url.Parse(c.server)
	if err != nil {
		return ListeningETH{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload)
	if err != nil {
		return ListeningETH{}, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		return ListeningETH{}, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&listening); err != nil {
		return ListeningETH{}, err
	}
	return listening, nil
}
