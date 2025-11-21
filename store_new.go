package geostore

import (
	"database/sql"
	"errors"

	"github.com/dracory/sb"
)

type NewStoreOptions struct {
	DB                 *sql.DB
	DbDriverName       string
	CountryTableName   string
	StateTableName     string
	TimezoneTableName  string
	AutomigrateEnabled bool
}

func NewStore(options NewStoreOptions) (*Store, error) {
	if options.CountryTableName == "" {
		return nil, errors.New("geo store: CountryTableName is required")
	}

	if options.CountryTableName == "" {
		return nil, errors.New("geo store: StateTableName is required")
	}

	if options.TimezoneTableName == "" {
		return nil, errors.New("geo store: TimezoneTableName is required")
	}

	if options.DB == nil {
		return nil, errors.New("shop store: DB is required")
	}

	if options.DbDriverName == "" {
		options.DbDriverName = sb.DatabaseDriverName(options.DB)
	}

	store := &Store{
		db:                 options.DB,
		dbDriverName:       options.DbDriverName,
		countryTableName:   options.CountryTableName,
		stateTableName:     options.StateTableName,
		timezoneTableName:  options.TimezoneTableName,
		automigrateEnabled: options.AutomigrateEnabled,
	}

	if store.automigrateEnabled {
		err := store.AutoMigrate()

		if err != nil {
			return nil, err
		}
	}

	return store, nil
}
