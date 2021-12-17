package bb8

import (
	"encoding/json"
	"log"
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
