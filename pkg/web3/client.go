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

	// Server url to use
	Server string

	// Max number of routines to use for *All methods
	MaxRoutines int

	RetryWaitMin time.Duration // Minimum time to wait
	RetryWaitMax time.Duration // Maximum time to wait
	RetryMax     int           // Maximum number of retries
}

// NewAPIEthClient ...
func NewAPIEthClient(options APIClientOptions) (APIClient, error) {
	if options.Server == "" {
		options.Server = GanacheDevNet
	}

	retryclient := retryablehttp.NewClient()
	retryclient.Logger = nil

	if options.MaxRoutines == 0 {
		options.MaxRoutines = 10
	}

	client := &apiClient{
		server:   options.Server,
		client:   retryclient,
		routines: options.MaxRoutines,
	}

	return client, nil
}

type APIClient interface {
	// VersionETH ...
	VersionETH(ctx context.Context) (Version, error)
	// ListeningETH ...
	ListeningETH(ctx context.Context) (Listening, error)
	// TransactionByHashETH ...
	TransactionByHashETH(ctx context.Context, payload PayloadReq) (Transaction, error)
	// TransactionReceiptETH ...
	TransactionReceiptETH(ctx context.Context, payload PayloadReq) (TransactionReceipt, error)
	// BalanceETH ...
	BalanceETH(ctx context.Context, payload PayloadReq) (Balance, error)
	// BlockByNumberETH ...
	BlockByNumberETH(ctx context.Context, payload PayloadReq) (Block, error)
}
