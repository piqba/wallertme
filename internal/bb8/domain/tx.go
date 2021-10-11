package bb8

import (
	"encoding/json"
	"log"
)

// ResultTx ...
type ResultTx struct {
	// Time `time` is a transaction's timestamp
	Time int64 `json:"time,omitempty"`
	// Txfrom `txfrom` sender's Ethereum address
	Txfrom string `json:"txfrom,omitempty"`
	// Txto `txto` recipient's Ethereum address
	Txto string `json:"txto,omitempty"`
	// Gas `gas` indicates `gasUsed`
	Gas int64 `json:"gas,omitempty"`
	// Gasprice `gasprice` indicates `gasPrice`
	Gasprice int64 `json:"gasprice,omitempty"`
	// Block `block` is a transaction's block number
	Block int64 `json:"block,omitempty"`
	// Txhash `txhash` is a transaction's hash
	Txhash string `json:"txhash,omitempty"`
	// Value `value` stores amount of ETH transferred
	Value int64 `json:"value,omitempty"`
	// ContractTo `contract_to` indicates recipient's Ethereum address in case of contract
	ContractTo string `json:"contract_to,omitempty"`
	// ContractValue `contract_value` stores amount of ERC20 transaction in its tokens
	ContractValue string `json:"contract_value,omitempty"`
}

// ToJSON ...
func (rtx *ResultTx) ToJSON() string {
	bytes, err := json.Marshal(rtx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
}
