package bb8

import (
	"encoding/json"
	"log"
)

// ResultLastTxSOL return the last TX by address SOL symbol
type ResultLastTxSOL struct {
	Addr      string `json:"addr,omitempty"`
	TxID      string `json:"tx_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	FromAddr  string `json:"from_addr,omitempty"`
	ToAddr    string `json:"to_addr,omitempty"`
	Balance   string `json:"balance,omitempty"`
	Ammount   string `json:"ammount,omitempty"`
	TypeTx    string `json:"type_tx,omitempty"`
}

// ToJSON ...
func (rtx *ResultLastTxSOL) ToJSON() string {
	bytes, err := json.Marshal(rtx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
}

// ToMAP ...
func (rtx *ResultLastTxSOL) ToMAP() (toHashMap map[string]interface{}, err error) {

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
func (rtx *ResultLastTxSOL) TruncateAddress(address string) string {
	prefix := address[0:8]
	sufix := address[len(address)-8:]
	cleanAddress := prefix + "..." + sufix
	return cleanAddress
}
