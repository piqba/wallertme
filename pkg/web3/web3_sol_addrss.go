package web3

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// LastTxByAddr get last txs by address
type LastTxByAddr struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Result  []ResultTxByAddressSOL `json:"result"`
	ID      int64                  `json:"id"`
}

// ResultTxByAddressSOL result object from getSignaturesForAddress method json-rpc
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

// LastTxByAddress get last TXs by address
func (c *apiClient) LastTxByAddress(ctx context.Context, payload PayloadReqJSONRPC) (lastTx LastTxByAddr, err error) {
	_, span := otel.Tracer(nameSOL).Start(ctx, "LastTxByAddress")
	defer span.End()
	requestUrl, err := url.Parse(c.server)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return lastTx, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload.ToReader())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return lastTx, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return lastTx, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&lastTx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return lastTx, err
	}
	span.SetAttributes(attribute.String("request.api", res.Status))

	return lastTx, nil
}
