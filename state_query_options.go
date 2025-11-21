package geostore

type StateQueryOptions struct {
	ID          string
	Status      string
	StatusIn    []string
	CountryCode string
	StateCode   string
	Offset      int
	Limit       int
	SortOrder   string
	OrderBy     string
	CountOnly   bool
	WithDeleted bool
}
