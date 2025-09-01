package geostore

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/dracory/sb"
	_ "github.com/mattn/go-sqlite3"
)

func initDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreCountryCreate(t *testing.T) {
	db := initDB("test_country_create.db")

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

	country := NewCountry().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("Unknown")

	err = store.CountryCreate(country)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreCountryFindByID(t *testing.T) {
	db := initDB("test_country_find_by_id.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "geo_country_find_by_id",
		StateTableName:     "geo_state_find_by_id",
		TimezoneTableName:  "geo_timezone_find_by_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	country := NewCountry().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("Unknown").
		SetIsoCode2("UN").
		SetIsoCode3("UNK")

	err = store.CountryCreate(country)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	countryFound, errFind := store.CountryFindByID(country.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if countryFound == nil {
		t.Fatal("Country MUST NOT be nil")
	}

	if countryFound.Name() != "Unknown" {
		t.Fatal("Country title MUST BE 'Unknown', found: ", countryFound.Name())
	}

	if countryFound.Status() != COUNTRY_STATUS_ACTIVE {
		t.Fatal("Country status MUST BE 'active', found: ", countryFound.Status())
	}

	if countryFound.IsoCode2() != "UN" {
		t.Fatal("Country iso_code2 MUST BE 'UN', found: ", countryFound.IsoCode2())
	}

	if countryFound.IsoCode3() != "UNK" {
		t.Fatal("Country iso_code3 MUST BE 'UNK', found: ", countryFound.IsoCode3())
	}

	if !strings.Contains(countryFound.DeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("Country MUST NOT be soft deleted", countryFound.DeletedAt())
	}
}

func TestStoreCountrySoftDelete(t *testing.T) {
	db := initDB("test_country_soft_delete.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "geo_country_find_by_id",
		StateTableName:     "geo_state_find_by_id",
		TimezoneTableName:  "geo_timezone_find_by_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	country := NewCountry().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("Unknown").
		SetIsoCode2("UN").
		SetIsoCode3("UNK")

	err = store.CountryCreate(country)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CountrySoftDeleteByID(country.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if country.DeletedAt() != sb.NULL_DATETIME {
		t.Fatal("Country MUST NOT be soft deleted")
	}

	countryFound, errFind := store.CountryFindByID(country.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if countryFound != nil {
		t.Fatal("Country MUST be nil")
	}

	countryFindWithDeleted, err := store.CountryList(CountryQueryOptions{
		ID:          country.ID(),
		Limit:       1,
		WithDeleted: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(countryFindWithDeleted) == 0 {
		t.Fatal("Exam MUST be soft deleted")
	}

	if strings.Contains(countryFindWithDeleted[0].DeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("Exam MUST be soft deleted", country.DeletedAt())
	}

}

func TestStoreStateAutomigrate(t *testing.T) {
	db := initDB("test_state_create.db")

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

	states, err := store.StateList(StateQueryOptions{})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(states) < 10 {
		t.Fatal("there must be states in the database")
	}
}

func TestStoreStateCreate(t *testing.T) {
	db := initDB("test_state_create.db")

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

	state := NewState().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("Unknown").
		SetStateCode("UN")

	err = store.StateCreate(state)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}
