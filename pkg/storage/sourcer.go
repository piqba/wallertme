package storage

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/logger"
)

const (
	// JSONFILE ...
	JSONFILE = "json"
	// DB ...
	DB = "db"
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
	Wallets() (wallets []Wallet, err error)
	WalletsTONotify() (wallets []Wallet, err error)
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

func (s *Source) Wallets() (wallets []Wallet, err error) {

	switch s.SourceType {
	case JSONFILE:
		wallets, err = WalletsFromJsonToStruct(s.Options.PathName, s.Options.FileName)
		if err != nil {
			return nil, errors.Errorf("storage: %v", err)
		}
		return wallets, nil
	case DB:
		wallets, err = WalletFromPostgreSQL(s.Options.Pgx)
		if err != nil {
			return nil, errors.Errorf("storage: %v", err)
		}
		return wallets, nil
	default:
		return nil, errors.Errorf("storage: %v", "Source not found")
	}
}

func (s *Source) WalletsTONotify() (wallets []Wallet, err error) {

	switch s.SourceType {
	case JSONFILE:
		wallets, err = WalletsFromJsonToStruct(s.Options.PathName, s.Options.FileName)
		if err != nil {
			return nil, errors.Errorf("storage: %v", err)
		}
		return wallets, nil
	case DB:
		wallets, err = WalletFromPostgreSQL(s.Options.Pgx)
		if err != nil {
			return nil, errors.Errorf("storage: %v", err)
		}
		return wallets, nil
	default:
		return nil, errors.Errorf("storage: %v", "Source not found")
	}
}
