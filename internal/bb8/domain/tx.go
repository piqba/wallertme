package bb8

import (
	"encoding/json"
	"log"

	"github.com/piqba/wallertme/pkg/web3"
)

type Txer interface {
	InfoByAddress(address string) (ResultInfoByAddr, error)
}

// ResultInfoByAddr ...
type ResultInfoByAddr struct {
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

type ResultLastTxByAddr struct {
	Addr          string `json:"addr,omitempty"`
	CtbID         string `json:"ctbId,omitempty"`
	CtbTimeIssued int64  `json:"ctbTimeIssued,omitempty"`
	FromAddr      string `json:"from_addr,omitempty"`
	ToAddr        string `json:"to_addr,omitempty"`
	Balance       string `json:"balance,omitempty"`
	Ammount       string `json:"ammount,omitempty"`
	TypeTx        string `json:"type_tx,omitempty"`
}

// ToJSON ...
func (rtx *ResultLastTxByAddr) ToJSON() string {
	bytes, err := json.Marshal(rtx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
}

func (rtx *ResultLastTxByAddr) ToMAP() (toHashMap map[string]interface{}, err error) {

	fromStruct, err := json.Marshal(rtx)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(fromStruct, &toHashMap); err != nil {
		return toHashMap, err
	}

	return toHashMap, nil
}

// ToJSON ...
func (rtx *ResultInfoByAddr) ToJSON() string {
	bytes, err := json.Marshal(rtx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
}
func (rtx *ResultInfoByAddr) ToMAP() (toHashMap map[string]interface{}, err error) {

	fromStruct, err := json.Marshal(rtx)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(fromStruct, &toHashMap); err != nil {
		return toHashMap, err
	}

	return toHashMap, nil
}
