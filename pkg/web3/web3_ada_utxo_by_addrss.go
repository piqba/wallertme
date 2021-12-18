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

// This data was fetch from the new endpoint api using graphql

type TxByAddrADAV2 struct {
	Data Data `json:"data"`
}

func (r *TxByAddrADAV2) ToJSON() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

type Data struct {
	Utxos []Utxo `json:"utxos"`
}

type Utxo struct {
	Value       string         `json:"value"`
	TxHash      string         `json:"txHash"`
	Transaction TransactionADA `json:"transaction"`
}

type TransactionADA struct {
	Block       Block  `json:"block"`
	Fee         int64  `json:"fee"`
	TotalOutput string `json:"totalOutput"`
	IncludedAt  string `json:"includedAt"`
	Inputs      []Put  `json:"inputs"`
	Outputs     []Put  `json:"outputs"`
}

type Block struct {
	Number int64  `json:"number"`
	Hash   string `json:"hash"`
}

type Put struct {
	Address string `json:"address"`
	Value   string `json:"value"`
}

// LastTxByAddressADA get last TXs by address
func (c *apiClient) LastTxByAddressADA(ctx context.Context, payload PayloadReqJSONGQL) (lastTx TxByAddrADAV2, err error) {
	_, span := otel.Tracer(nameSOL).Start(ctx, "LastTxByAddressADA")
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
