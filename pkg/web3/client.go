package web3

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// apiClient core struct for our api client logic
type apiClient struct {
	server   string
	client   *retryablehttp.Client
	routines int
}

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// APIClientOptions for management all clients
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

// NewAPICardanoClient is a function in charge of to proccess all logic from cardano API
//
//
// More info visit here -> https://explorer-api.testnet.dandelion.link
func NewAPICardanoClient(options APIClientOptions) (APICardanoClient, error) {
	_, span := otel.Tracer(nameClientAPI).Start(context.Background(), "NewAPICardanoClient")
	defer span.End()
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
	span.SetAttributes(attribute.String("create.client", client.server))

	return client, nil
}

// APICardanoClient define all methods that can be used on our program
type APICardanoClient interface {
	// InfoByAddress ...
	InfoByAddress(ctx context.Context, address string) (AddrSumary, error)
	LastTxByAddressADA(ctx context.Context, payload PayloadReqJSONGQL) (TxByAddrADAV2, error)
}

// NewAPISolanaClient is a function in charge of to proccess all logic from solana API
//
//
//  More info visit here -> https://api.devnet.solana.com
func NewAPISolanaClient(options APIClientOptions) (APISolanaClient, error) {
	_, span := otel.Tracer(nameClientAPI).Start(context.Background(), "NewAPISolanaClient")
	defer span.End()
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
	span.SetAttributes(attribute.String("create.client", client.server))
	return client, nil
}

// APISolanaClient define all methods that can be used on our program

type APISolanaClient interface {
	// LastTxByAddress ...
	LastTxByAddress(ctx context.Context, payload PayloadReqJSONRPC) (LastTxByAddr, error)
	// InfoByTx ...
	InfoByTx(ctx context.Context, payload PayloadReqJSONRPC) (TxInfo, error)
}
