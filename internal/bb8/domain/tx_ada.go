package bb8

import (
	"encoding/json"
	"log"

	"github.com/piqba/wallertme/pkg/web3"
)

// ResultLastTxADA return the last TX by ADA symbol (Cardano blockchain)
type ResultLastTxADA struct {
	Addr          string `json:"addr,omitempty"`
	CtbID         string `json:"ctbId,omitempty"`
	CtbTimeIssued string `json:"ctbTimeIssued,omitempty"`
	FromAddr      string `json:"from_addr,omitempty"`
	ToAddr        string `json:"to_addr,omitempty"`
	Balance       string `json:"balance,omitempty"`
	Ammount       string `json:"ammount,omitempty"`
	TypeTx        string `json:"type_tx,omitempty"`
	Symbol        string `json:"symbol,omitempty"`
}

// ToJSON ...
func (rtx *ResultLastTxADA) ToJSON() string {
	bytes, err := json.Marshal(rtx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
}

// ToMAP ...
func (rtx *ResultLastTxADA) ToMAP() (toHashMap map[string]interface{}, err error) {

	fromStruct, err := json.Marshal(rtx)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(fromStruct, &toHashMap); err != nil {
		return toHashMap, err
	}

	return toHashMap, nil
}

// TruncateAddress ...
func (rtx *ResultLastTxADA) TruncateAddress(address string) string {
	prefix := address[0:16]
	sufix := address[len(address)-16:]
	cleanAddress := prefix + "..." + sufix
	return cleanAddress
}

// ResultInfoForADA information for ADA address
type ResultInfoForADA struct {
	Address   string          `json:"address,omitempty"`
	Type      string          `json:"type,omitempty"`
	BlockNO   int64           `json:"block_no,omitempty"`
	BlockHash string          `json:"block_hash,omitempty"`
	TxTotal   int64           `json:"tx_total,omitempty"`
	Balance   string          `json:"balance,omitempty"`
	TotalIn   string          `json:"total_in,omitempty"`
	TotalOut  string          `json:"total_out,omitempty"`
	TotalFee  string          `json:"total_fee,omitempty"`
	TxList    []web3.CATxList `json:"tx_list,omitempty"`
}

// ToJSON ...
func (rtx *ResultInfoForADA) ToJSON() string {
	bytes, err := json.Marshal(rtx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
}

// ToMAP ...
func (rtx *ResultInfoForADA) ToMAP() (toHashMap map[string]interface{}, err error) {

	fromStruct, err := json.Marshal(rtx)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(fromStruct, &toHashMap); err != nil {
		return toHashMap, err
	}

	return toHashMap, nil
}
