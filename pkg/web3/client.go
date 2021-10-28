package web3

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type apiClient struct {
	server   string
	client   *retryablehttp.Client
	routines int
}

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// APIClientOptions ...
type APIClientOptions struct {

	// NetworkType name of network only support (CardanoTestNet|SolanaDevNet)
	NetworkType string
	// Server url to use Api Endpoint
	server string

	// Max number of routines to use for *All methods
	MaxRoutines int

	RetryWaitMin time.Duration // Minimum time to wait
	RetryWaitMax time.Duration // Maximum time to wait
	RetryMax     int           // Maximum number of retries
}

// NewAPICardanoClient ...
func NewAPICardanoClient(options APIClientOptions) (APICardanoClient, error) {
	if options.NetworkType == "" {
		options.NetworkType = CardanoTestNet
	}
	if options.server == "" {
		options.server = NetworkMap[CardanoTestNet].ApiURL
	}

	retryclient := retryablehttp.NewClient()
	retryclient.Logger = nil

	if options.MaxRoutines == 0 {
		options.MaxRoutines = 10
	}

	client := &apiClient{
		server:   options.server,
		client:   retryclient,
		routines: options.MaxRoutines,
	}

	return client, nil
}

type APICardanoClient interface {
	// InfoByAddress ...
	InfoByAddress(ctx context.Context, address string) (AddrSumary, error)
}

// NewAPISolanaClient ...
func NewAPISolanaClient(options APIClientOptions) (APISolanaClient, error) {
	if options.NetworkType == "" {
		options.NetworkType = SolanaDevNet
	}
	if options.server == "" {
		options.server = NetworkMap[SolanaDevNet].ApiURL
	}

	retryclient := retryablehttp.NewClient()
	retryclient.Logger = nil

	if options.MaxRoutines == 0 {
		options.MaxRoutines = 10
	}

	client := &apiClient{
		server:   options.server,
		client:   retryclient,
		routines: options.MaxRoutines,
	}

	return client, nil
}

type APISolanaClient interface {
	// LastTxByAddress ...
	LastTxByAddress(ctx context.Context, payload PayloadReqJSONRPC) (LastTxByAddr, error)
	// InfoByTx ...
	InfoByTx(ctx context.Context, payload PayloadReqJSONRPC) (TxInfo, error)
}
