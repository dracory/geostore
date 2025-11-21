package geostore

import "context"

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool)

	CountryCreate(ctx context.Context, country *Country) error
	CountryDelete(ctx context.Context, country *Country) error
	CountryDeleteByID(ctx context.Context, countryID string) error
	CountryFindByID(ctx context.Context, countryID string) (*Country, error)
	CountryFindByIso2(ctx context.Context, iso2Code string) (*Country, error)
	CountryList(ctx context.Context, options CountryQueryOptions) ([]Country, error)
	CountrySoftDelete(ctx context.Context, discount *Country) error
	CountrySoftDeleteByID(ctx context.Context, discountID string) error
	CountryUpdate(ctx context.Context, country *Country) error
	TimezoneCreate(ctx context.Context, timezone *Timezone) error
	TimezoneList(ctx context.Context, options TimezoneQueryOptions) ([]Timezone, error)
}
