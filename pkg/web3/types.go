package web3

import "fmt"

const (
	PolygonMainNet = "https://polygon-rpc.com/"
	PolygonTestNet = "https://matic-mumbai.chainstacklabs.com/"
)

// APIError is used to describe errors from the API.
type APIError struct {
	Response interface{}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error, %+v", e.Response)
}
