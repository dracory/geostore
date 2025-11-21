package geostore

type TimezoneQueryOptions struct {
	ID          string
	Status      string
	StatusIn    []string
	CountryCode string
	Timezone    string
	Offset      int
	Limit       int
	SortOrder   string
	OrderBy     string
	CountOnly   bool
	WithDeleted bool
}
