package storage

import (
	"context"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	// JSONFILE ...
	JSONFILE = "json"
	// DB ...
	DB = "db"
	// namePgxClient is the Tracer namePgxClient used to identify this instrumentation library.
	namePgxClient = "storage.pgx.client"
	// nameSourcer is the Tracer nameSourcer used to identify this instrumentation library.
	nameSourcer = "storage.sourcer"
	// nameSourcerPostgres is the Tracer nameSourcerPostgres used to identify this instrumentation library.
	nameSourcerPostgres = "storage.sourcer.postgres"
	// nameSourcerJson is the Tracer nameSourcerJson used to identify this instrumentation library.
	nameSourcerJson = "storage.sourcer.jsonfile"
)

type OptionsSource struct {
	FileName string
	PathName string
	Pgx      *sqlx.DB
}
type Source struct {
	SourceType string
	Options    OptionsSource
}

func NewSource(sourceType string, options OptionsSource) Source {
	return Source{
		SourceType: sourceType,
		Options:    options,
	}
}

// Sourcer ...
type Sourcer interface {
	Wallets(ctx context.Context) (wallets []Wallet, err error)
	WalletsTONotify(ctx context.Context) (wallets []Wallet, err error)
}

// Wallet ...
type Wallet struct {
	Address         string            `json:"address,omitempty"`
	Symbol          string            `json:"symbol,omitempty"`
	IsActive        bool              `json:"is_active,omitempty"`
	NotifierService []NotifierService `json:"notifier_service,omitempty"`
	NetworkType     string            `json:"network_type,omitempty"`
}

func (w *Wallet) ToJSON() string {
	bytes, err := json.Marshal(w)
	if err != nil {
		logger.LogError(errors.Errorf("storage: %v", err).Error())

	}
	return string(bytes)
}

// NotifierService ...
type NotifierService struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}

func (s *Source) Wallets(ctx context.Context) (wallets []Wallet, err error) {
	_, span := otel.Tracer(nameSourcer).Start(ctx, "Wallets")
	defer span.End()
	switch s.SourceType {
	case JSONFILE:
		wallets, err = WalletsFromJsonToStruct(ctx, s.Options.PathName, s.Options.FileName)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, errors.Errorf("storage: %v", err)
		}
		span.SetAttributes(attribute.String("sourcer.get.wallets", "Success"))

		return wallets, nil
	case DB:
		wallets, err = WalletFromPostgreSQL(ctx, s.Options.Pgx)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, errors.Errorf("storage: %v", err)
		}
		span.SetAttributes(attribute.String("sourcer.get.wallets", "Success"))

		return wallets, nil
	default:
		err = errors.Errorf("storage: %v", "Source not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

}

func (s *Source) WalletsTONotify(ctx context.Context) (wallets []Wallet, err error) {
	_, span := otel.Tracer(nameSourcer).Start(ctx, "Wallets")
	defer span.End()
	switch s.SourceType {
	case JSONFILE:
		wallets, err = WalletsFromJsonToStruct(ctx, s.Options.PathName, s.Options.FileName)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, errors.Errorf("storage: %v", err)
		}
		return wallets, nil
	case DB:
		wallets, err = WalletFromPostgreSQL(ctx, s.Options.Pgx)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, errors.Errorf("storage: %v", err)
		}
		return wallets, nil
	default:
		err = errors.Errorf("storage: %v", "Source not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
}
