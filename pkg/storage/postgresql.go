package storage

import (
	"context"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func WalletFromPostgreSQL(ctx context.Context, pgx *sqlx.DB) ([]Wallet, error) {
	_, span := otel.Tracer(nameSourcerPostgres).Start(ctx, "WalletFromPostgreSQL")
	defer span.End()
	var wallets []Wallet

	query := `select wallet from wallets where wallet @> '{"is_active": true}'`
	stmt, err := pgx.Prepare(query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	for rows.Next() {
		var data []byte
		var wallet Wallet
		err = rows.Scan(&data)
		json.Unmarshal(data, &wallet)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	span.SetAttributes(attribute.Int("postgresql.query.result", len(wallets)))

	return wallets, nil
}
