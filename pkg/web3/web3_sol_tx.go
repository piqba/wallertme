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

// TxInfo return information related to a TX
type TxInfo struct {
	Jsonrpc string       `json:"jsonrpc"`
	Result  ResultTxInfo `json:"result"`
	ID      int64        `json:"id"`
}

func (r *TxInfo) ToJSON() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

// ResultTxInfo return result information by an specific TX
type ResultTxInfo struct {
	BlockTime   int64       `json:"blockTime"`
	Meta        Meta        `json:"meta"`
	Slot        int64       `json:"slot"`
	Transaction Transaction `json:"transaction"`
}

// Meta object with tx information
type Meta struct {
	Err               interface{}   `json:"err"`
	Fee               int64         `json:"fee"`
	InnerInstructions []interface{} `json:"innerInstructions"`
	LogMessages       []string      `json:"logMessages"`
	PostBalances      []int64       `json:"postBalances"`
	PostTokenBalances []interface{} `json:"postTokenBalances"`
	PreBalances       []int64       `json:"preBalances"`
	PreTokenBalances  []interface{} `json:"preTokenBalances"`
	Rewards           []interface{} `json:"rewards"`
	Status            Status        `json:"status"`
}

type Status struct {
	Ok interface{} `json:"Ok"`
}

// Transaction ...
type Transaction struct {
	Message    Message  `json:"message"`
	Signatures []string `json:"signatures"`
}

// Message ...
type Message struct {
	AccountKeys     []string      `json:"accountKeys"`
	Header          Header        `json:"header"`
	Instructions    []Instruction `json:"instructions"`
	RecentBlockhash string        `json:"recentBlockhash"`
}

// Header ...
type Header struct {
	NumReadonlySignedAccounts   int64 `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts int64 `json:"numReadonlyUnsignedAccounts"`
	NumRequiredSignatures       int64 `json:"numRequiredSignatures"`
}

// Instruction ...
type Instruction struct {
	Accounts       []int64 `json:"accounts"`
	Data           string  `json:"data"`
	ProgramIDIndex int64   `json:"programIdIndex"`
}

// InfoByTx get info by TX
func (c *apiClient) InfoByTx(ctx context.Context, payload PayloadReqJSONRPC) (info TxInfo, err error) {
	_, span := otel.Tracer(nameSOL).Start(ctx, "InfoByTx")
	defer span.End()
	requestUrl, err := url.Parse(c.server)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return info, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload.ToReader())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return info, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return info, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&info); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return info, err
	}
	span.SetAttributes(attribute.String("request.api", res.Status))

	return info, nil
}
