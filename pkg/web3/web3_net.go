package web3

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type Version struct {
	JsonRPC string `json:"jsonrpc,omitempty"`
	ID      int    `json:"id,omitempty"`
	Result  string `json:"result,omitempty"`
}

func (v *Version) ToJSON() string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

type Listening struct {
	JsonRPC string `json:"jsonrpc,omitempty"`
	ID      int    `json:"id,omitempty"`
	Result  string `json:"result,omitempty"`
}

func (l *Listening) ToJSON() string {
	bytes, err := json.Marshal(l)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func (c *apiClient) Version(ctx context.Context) (version Version, err error) {
	payload := strings.NewReader(`{
		"jsonrpc":"2.0",
		"method":"net_version",
		"params":[],
		"id":67
	}`)
	requestUrl, err := url.Parse(c.server)
	if err != nil {
		return Version{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload)
	if err != nil {
		return Version{}, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		return Version{}, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&version); err != nil {
		return Version{}, err
	}
	return version, nil
}
func (c *apiClient) Listening(ctx context.Context) (listening Listening, err error) {
	payload := strings.NewReader(`{
		"jsonrpc":"2.0",
		"method":"net_listening",
		"params":[],
		"id":67
	}`)
	requestUrl, err := url.Parse(c.server)
	if err != nil {
		return Listening{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload)
	if err != nil {
		return Listening{}, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		return Listening{}, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&listening); err != nil {
		return Listening{}, err
	}
	return listening, nil
}
