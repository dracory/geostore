package geostore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/dracory/neat"
)

// NewStoreOptions define the options for creating a new geostore
type NewStoreOptions struct {
	CountryTableName   string
	StateTableName     string
	TimezoneTableName  string
	DB                 *sql.DB
	AutomigrateEnabled bool
	AutoseedEnabled    bool
}

// NewStore creates a new geostore
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.CountryTableName == "" {
		return nil, errors.New("geostore: CountryTableName is required")
	}

	if opts.StateTableName == "" {
		return nil, errors.New("geostore: StateTableName is required")
	}

	if opts.TimezoneTableName == "" {
		return nil, errors.New("geostore: TimezoneTableName is required")
	}

	if opts.DB == nil {
		return nil, errors.New("geostore: DB is required")
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	store := &storeImplementation{
		countryTableName:   opts.CountryTableName,
		stateTableName:     opts.StateTableName,
		timezoneTableName:  opts.TimezoneTableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		autoseedEnabled:    opts.AutoseedEnabled,
		db:                 neatDB,
		logger:             logger,
	}

	if store.automigrateEnabled {
		err := store.MigrateUp(context.Background())
		if err != nil {
			return nil, err
		}
	}

	if store.autoseedEnabled {
		err := store.Seed(context.Background())
		if err != nil {
			return nil, err
		}
	}

	return store, nil
}
