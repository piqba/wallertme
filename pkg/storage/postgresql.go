package storage

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

func WalletFromPostgreSQL(pgx *sqlx.DB) ([]Wallet, error) {
	var wallets []Wallet

	query := `select wallet from wallets where wallet @> '{"is_active": true}'`
	stmt, err := pgx.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var data []byte
		var wallet Wallet
		err = rows.Scan(&data)
		json.Unmarshal(data, &wallet)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	return wallets, nil
}
