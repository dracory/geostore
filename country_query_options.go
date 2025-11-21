package geostore

type CountryQueryOptions struct {
	ID          string
	IDIn        []string
	Status      string
	StatusIn    []string
	Iso2        string
	Iso3        string
	Offset      int
	Limit       int
	SortOrder   string
	OrderBy     string
	CountOnly   bool
	WithDeleted bool
}