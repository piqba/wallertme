package web3

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// resourceAddress endpoint resource for get summary by address
	resourceAddress = "api/addresses/summary"
)

// AddrSumary it`s a result for the cardano API
type AddrSumary struct {
	Result Right `json:"Right"`
}

// Right main object
type Right struct {
	CAAddress     string     `json:"caAddress"`
	CAType        string     `json:"caType"`
	CAChainTip    CAChainTip `json:"caChainTip"`
	CATxNum       int64      `json:"caTxNum"`
	CABalance     CABalance  `json:"caBalance"`
	CATotalInput  CABalance  `json:"caTotalInput"`
	CATotalOutput CABalance  `json:"caTotalOutput"`
	CATotalFee    CABalance  `json:"caTotalFee"`
	CATxList      []CATxList `json:"caTxList"`
}

// ToJSON convert to JSON this struct
func (r *Right) ToJSON() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

// CABalance get balance
type CABalance struct {
	GetCoin string `json:"getCoin"`
}

// CAChainTip get general info from Chain
type CAChainTip struct {
	CTBlockNo   int64  `json:"ctBlockNo"`
	CTSlotNo    int64  `json:"ctSlotNo"`
	CTBlockHash string `json:"ctBlockHash"`
}

// CATxList get info from TXs
type CATxList struct {
	CtbID         string    `json:"ctbId"`
	CtbTimeIssued int64     `json:"ctbTimeIssued"`
	CtbInputs     []CtbPut  `json:"ctbInputs"`
	CtbOutputs    []CtbPut  `json:"ctbOutputs"`
	CtbInputSum   CABalance `json:"ctbInputSum"`
	CtbOutputSum  CABalance `json:"ctbOutputSum"`
	CtbFees       CABalance `json:"ctbFees"`
}

// CtbPut general info from Put object
type CtbPut struct {
	CtaAddress string    `json:"ctaAddress"`
	CtaAmount  CABalance `json:"ctaAmount"`
	CtaTxHash  string    `json:"ctaTxHash"`
	CtaTxIndex int64     `json:"ctaTxIndex"`
}

// InfoByAddress get info by address
func (c *apiClient) InfoByAddress(ctx context.Context, address string) (AddrSumary, error) {
	requestUrl, err := url.Parse(fmt.Sprintf("%s/%s/%s", c.server, resourceAddress, address))
	if err != nil {
		return AddrSumary{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return AddrSumary{}, err
	}
	res, err := c.handleRequest(req)
	if err != nil {
		return AddrSumary{}, err
	}
	defer res.Body.Close()
	sumary := AddrSumary{}
	if err = json.NewDecoder(res.Body).Decode(&sumary); err != nil {
		return AddrSumary{}, err
	}
	return sumary, nil
}
