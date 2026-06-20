package geostore

import (
	"context"
	"database/sql"
)

type StoreInterface interface {
	// GetCountryTableName returns the country table name
	GetCountryTableName() string
	// SetCountryTableName sets the country table name
	SetCountryTableName(countryTableName string)

	// GetStateTableName returns the state table name
	GetStateTableName() string
	// SetStateTableName sets the state table name
	SetStateTableName(stateTableName string)

	// GetTimezoneTableName returns the timezone table name
	GetTimezoneTableName() string
	// SetTimezoneTableName sets the timezone table name
	SetTimezoneTableName(timezoneTableName string)

	// MigrateUp creates all database tables
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error

	// MigrateDown drops all database tables
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error

	// Seed populates all tables with initial data
	Seed(ctx context.Context, tx ...*sql.Tx) error

	EnableDebug(debug bool)

	CountryCreate(ctx context.Context, country *Country) error
	CountryDelete(ctx context.Context, country *Country) error
	CountryDeleteByID(ctx context.Context, countryID string) error
	CountryFindByID(ctx context.Context, countryID string) (*Country, error)
	CountryFindByIso2(ctx context.Context, iso2Code string) (*Country, error)
	CountryList(ctx context.Context, options CountryQueryOptions) ([]Country, error)
	CountrySoftDelete(ctx context.Context, country *Country) error
	CountrySoftDeleteByID(ctx context.Context, countryID string) error
	CountryUpdate(ctx context.Context, country *Country) error

	StateCreate(ctx context.Context, state *State) error
	StatesCreate(ctx context.Context, states []*State) error
	StateList(ctx context.Context, options StateQueryOptions) ([]State, error)

	TimezoneCreate(ctx context.Context, timezone *Timezone) error
	TimezoneList(ctx context.Context, options TimezoneQueryOptions) ([]Timezone, error)
}
