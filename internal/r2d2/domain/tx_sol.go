package r2d2

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/piqba/wallertme/pkg/logger"
)

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

func (tx *ResultLastTxSOL) TemplateTelegram() string {
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

	timestampUnix, err := strconv.ParseInt(tx.Timestamp, 10, 64)
	if err != nil {
		logger.LogError(err.Error())
	}
	//Unix Timestamp to time.Time
	timeT := time.Unix(timestampUnix, 0)
	var msg string
	if tx.TypeTx == TxSender {

		msg = "💱Symbol:%s\nTxID: https://explorer.solana.com/tx/%s?cluster=devnet\n📡 Address: %s\n 💰 Balance: %v  ₳\n💵 Ammount: %v ₳\n⬅️ TypeTx: %s\n💳 From: %s\n💳 TO: %s\n⏰ Time: %s"
	} else {
		msg = "💱Symbol: %s\nTxID: https://explorer.solana.com/tx/%s?cluster=devnet\n📡 Address: %s\n💰 Balance: %v ₳\n💵 Ammount: %v ₳\n➡️ TypeTx: %s\n💳 From: %s\n💳 TO: %s\n⏰ Time: %s"
	}
	return fmt.Sprintf(
		msg,
		"SOL",
		tx.TxID,
		tx.TruncateAddress(tx.Addr),
		newBalance,
		newAmmount,
		tx.TypeTx,
		tx.TruncateAddress(tx.FromAddr),
		tx.TruncateAddress(tx.ToAddr),
		timeT.String(),
	)
}

func (tx *ResultLastTxSOL) TemplateDiscord() string {
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

	timestampUnix, err := strconv.ParseInt(tx.Timestamp, 10, 64)
	if err != nil {
		logger.LogError(err.Error())
	}
	//Unix Timestamp to time.Time
	timeT := time.Unix(timestampUnix, 0)
	var msg string
	if tx.TypeTx == TxSender {

		msg = "💱Symbol: **`%s`**\n🆔 [Show TxID](https://explorer.solana.com/tx/%s?cluster=devnet)\n📡 Address: **%s**\n 💰 Balance: `%v  ₳`\n💵 Ammount: `%v  ₳`\n⬅️ TypeTx: `%s`\n💳 From: **%s**\n💳 TO: **%s**\n⏰ Time: `%s`"
	} else {
		msg = "💱 Symbol: **`%s`**\n🆔 [Show TxID](https://explorer.solana.com/tx/%s?cluster=devnet)\n📡 Address: **%s**\n💰 Balance: `%v  ₳`\n💵 Ammount: `%v  ₳`\n➡️ TypeTx: `%s`\n💳 From: **%s**\n💳 TO: **%s**\n⏰ Time: `%s`"
	}
	return fmt.Sprintf(
		msg,
		"SOL",
		tx.TxID,
		tx.TruncateAddress(tx.Addr),
		newBalance,
		newAmmount,
		tx.TypeTx,
		tx.TruncateAddress(tx.FromAddr),
		tx.TruncateAddress(tx.ToAddr),
		timeT.String(),
	)
}

func (tx *ResultLastTxSOL) TemplateSMTP() string {
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

	timestampUnix, err := strconv.ParseInt(tx.Timestamp, 10, 64)
	if err != nil {
		logger.LogError(err.Error())
	}
	//Unix Timestamp to time.Time
	timeT := time.Unix(timestampUnix, 0)
	var msg string
	if tx.TypeTx == TxSender {

		msg = `
		<!DOCTYPE HTML PULBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
		<html>
		<head>
		<meta http-equiv="content-type" content="text/html"; charset=ISO-8859-1">
		</head>
		<body>
		<p>💱Symbol: %s</p>
		<br>
		<a href="https://explorer.solana.com/tx/%s?cluster=devnet">🆔 show TxID</a>
		<br>
		<p>📡 Address: %s</p>
		<br>
		<strong>💰 Balance: %v  ₳</strong>
		<br>
		<strong>💵 Ammount: %v ₳</strong>
		<br>
		<p>⬅️ TypeTx: %s</p>
		<br>
		<code>💳 From: %s</code>
		<br>
		<code>💳 TO: %s</code>
		<br>
		<p>⏰ Time: %s</p>
		<br>
		</body>
		</html>
		`
	} else {
		msg = `
		<!DOCTYPE HTML PULBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
		<html>
		<head>
		<meta http-equiv="content-type" content="text/html"; charset=ISO-8859-1">
		</head>
		<body>
		<p>💱Symbol: %s</p>
		<br>
		<a href="https://explorer.solana.com/tx/%s?cluster=devnet">🆔 show TxID</a>
		<br>
		<p>📡 Address: %s</p>
		<br>
		<strong>💰 Balance: %v  ₳</strong>
		<br>
		<strong>💵 Ammount: %v ₳</strong>
		<br>
		<p>➡️ TypeTx: %s</p>
		<br>
		<code>💳 From: %s</code>
		<br>
		<code>💳 TO: %s</code>
		<br>
		<p>⏰ Time: %s</p>
		<br>
		</body>
		</html>
		`
	}
	return fmt.Sprintf(
		msg,
		"SOL",
		tx.TxID,
		tx.TruncateAddress(tx.Addr),
		newBalance,
		newAmmount,
		tx.TypeTx,
		tx.TruncateAddress(tx.FromAddr),
		tx.TruncateAddress(tx.ToAddr),
		timeT.String(),
	)

}
