package r2d2

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/piqba/wallertme/pkg/logger"
)

type ResultLastTxByAddr struct {
	Addr          string `json:"addr,omitempty"`
	CtbID         string `json:"ctbId,omitempty"`
	CtbTimeIssued string `json:"ctbTimeIssued,omitempty"`
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

// TruncateAddress ...
func (rtx *ResultLastTxByAddr) TruncateAddress(address string) string {
	prefix := address[0:16]
	sufix := address[len(address)-16:]
	cleanAddress := prefix + "..." + sufix
	return cleanAddress
}

func (tx *ResultLastTxByAddr) Hummanify() string {
	balance, err := strconv.ParseInt(tx.Balance, 10, 64)
	if err != nil {
		logger.LogError(err.Error())
	}
	newBalance := float64(balance) / 1_000_000
	ammount, err := strconv.ParseInt(tx.Ammount, 10, 64)
	if err != nil {
		logger.LogError(err.Error())
	}
	newAmmount := float64(ammount) / 1_000_000

	timestampUnix, err := strconv.ParseInt(tx.CtbTimeIssued, 10, 64)
	if err != nil {
		logger.LogError(err.Error())
	}
	//Unix Timestamp to time.Time
	timeT := time.Unix(timestampUnix, 0)
	msg := `
	📡 Address: %s

	🆔 TxID: %s

	💰 Balance: %v

	💵 Ammount: %v

	↔️ TypeTx: %s

	💳 From: %s

	💳 TO: %s

	⏰ Time: %s

`
	return fmt.Sprintf(
		msg,
		tx.TruncateAddress(tx.Addr),
		tx.CtbID,
		newBalance,
		newAmmount,
		tx.TypeTx,
		tx.TruncateAddress(tx.FromAddr),
		tx.TruncateAddress(tx.ToAddr),
		timeT.String(),
	)
}
