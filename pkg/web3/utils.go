package web3

import (
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

// handleAPIErrorResponse ...
func handleAPIErrorResponse(res *http.Response) error {
	switch res.StatusCode {
	case 400:
		return &APIError{}
	case 403:
		return &APIError{}
	case 404:
		return &APIError{}
	case 429:
		return &APIError{}
	case 418:
		return &APIError{}
	case 500:
		return &APIError{}
	default:
		return &APIError{}
	}
}

// handleRequest ...
func (c *apiClient) handleRequest(req *http.Request) (res *http.Response, err error) {
	req.Header.Add("Content-Type", "application/json")
	rreq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return
	}
	res, err = c.client.Do(rreq)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		return res, handleAPIErrorResponse(res)
	}

	return res, nil
}
