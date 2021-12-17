package r2d2

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/web3"
)

const (
	// TxSender OUT
	TxSender = "sender"
	// TxReceiver IN
	TxReceiver = "receiver"
)

// ResultLastTxADA ...
type ResultLastTxADA struct {
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
func (rtx *ResultLastTxADA) ToJSON() string {
	bytes, err := json.Marshal(rtx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
}

// parseField ...
func (tx *ResultLastTxADA) parseField() (float64, float64, string) {
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

	// timestampUnix, err := strconv.ParseInt(tx.CtbTimeIssued, 10, 64)
	// if err != nil {
	// 	logger.LogError(err.Error())
	// }
	// //Unix Timestamp to time.Time
	// timeT := time.Unix(timestampUnix, 0)
	return newBalance, newAmmount, tx.CtbTimeIssued
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

// TemplateTelegram ...
func (tx *ResultLastTxADA) TemplateTelegram() string {
	newBalance, newAmmount, timeT := tx.parseField()
	var msg string
	if tx.TypeTx == TxSender {

		msg = "💱Symbol:%s\nTxID: %s\n📡 Address: %s\n 💰 Balance: %v  ₳\n💵 Ammount: %v ₳\n⬅️ TypeTx: %s\n💳 From: %s\n💳 TO: %s\n⏰ Time: %s"
	} else {
		msg = "💱Symbol: %s\nTxID: %s\n📡 Address: %s\n💰 Balance: %v ₳\n💵 Ammount: %v ₳\n➡️ TypeTx: %s\n💳 From: %s\n💳 TO: %s\n⏰ Time: %s"
	}
	return fmt.Sprintf(
		msg,
		"ADA",
		fmt.Sprintf(web3.NetworkMap[web3.CardanoTestNet].ExplorerURL, tx.CtbID),
		tx.TruncateAddress(tx.Addr),
		newBalance,
		newAmmount,
		tx.TypeTx,
		tx.TruncateAddress(tx.FromAddr),
		tx.TruncateAddress(tx.ToAddr),
		timeT,
	)
}

// TemplateDiscord ...
func (tx *ResultLastTxADA) TemplateDiscord() string {
	newBalance, newAmmount, timeT := tx.parseField()
	var msg string
	if tx.TypeTx == TxSender {

		msg = "💱Symbol: **`%s`**\n🆔 [Show TxID](%s)\n📡 Address: **%s**\n 💰 Balance: `%v  ₳`\n💵 Ammount: `%v  ₳`\n⬅️ TypeTx: `%s`\n💳 From: **%s**\n💳 TO: **%s**\n⏰ Time: `%s`"
	} else {
		msg = "💱 Symbol: **`%s`**\n🆔 [Show TxID](%s)\n📡 Address: **%s**\n💰 Balance: `%v  ₳`\n💵 Ammount: `%v  ₳`\n➡️ TypeTx: `%s`\n💳 From: **%s**\n💳 TO: **%s**\n⏰ Time: `%s`"
	}
	return fmt.Sprintf(
		msg,
		"ADA",
		fmt.Sprintf(web3.NetworkMap[web3.CardanoTestNet].ExplorerURL, tx.CtbID),
		tx.TruncateAddress(tx.Addr),
		newBalance,
		newAmmount,
		tx.TypeTx,
		tx.TruncateAddress(tx.FromAddr),
		tx.TruncateAddress(tx.ToAddr),
		timeT,
	)
}

// TemplateSMTP ...
func (tx *ResultLastTxADA) TemplateSMTP() string {
	newBalance, newAmmount, timeT := tx.parseField()
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
		<a href="%s">🆔 show TxID</a>
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
		<a href="%s">🆔 show TxID</a>
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
		"ADA",
		fmt.Sprintf(web3.NetworkMap[web3.CardanoTestNet].ExplorerURL, tx.CtbID),
		tx.TruncateAddress(tx.Addr),
		newBalance,
		newAmmount,
		tx.TypeTx,
		tx.TruncateAddress(tx.FromAddr),
		tx.TruncateAddress(tx.ToAddr),
		timeT,
	)

}
