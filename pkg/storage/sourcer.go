package storage

import (
	"github.com/piqba/wallertme/pkg/errors"
)

const (
	// JSONFILE ...
	JSONFILE = "json"
	// POSTGRESQL ...
	POSTGRESQL = "postgresql"
)

type OptionsSource struct {
	FileName string
	PathName string
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

func (s *Source) Wallets() (wallets []Wallet, err error) {

	switch s.SourceType {
	case JSONFILE:
		wallets, err = WalletsFromJsonToStruct(s.Options.PathName, s.Options.FileName)
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
	default:
		return nil, errors.Errorf("storage: %v", "Source not found")
	}
}
