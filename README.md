# Geo Store <a href="https://gitpod.io/#https://github.com/dracory/geostore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

![tests](https://github.com/dracory/geostore/workflows/tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/geostore)](https://goreportcard.com/report/github.com/dracory/geostore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/geostore)](https://pkg.go.dev/github.com/dracory/geostore)

## Usage:

```golang
store, err := NewStore(NewStoreOptions{
  DB:                 db,
  CountryTableName:   "geo_country",
  StateTableName:     "geo_state",
  TimezoneTableName:  "geo_timezone",
  AutomigrateEnabled: true,
})

if err != nil {
  t.Fatal("unexpected error:", err)
}

if store == nil {
  t.Fatal("unexpected nil store")
}

country, errFind := store.CountryFindByIso2("BG")

if errFind != nil {
  t.Fatal("unexpected error:", errFind)
}

if country == nil {
  t.Fatal("Country MUST NOT be nil")
}

log.Print(country.Name()) // Prints the name of the country
```
