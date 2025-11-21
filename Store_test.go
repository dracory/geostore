package geostore

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/dracory/sb"
	_ "modernc.org/sqlite"
)

func initDB() *sql.DB {
	dsn := ":memory:" + "?cache=shared&mode=memory"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}

	// Enable foreign key constraints for SQLite
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		panic(err)
	}

	return db
}

// setupTestDB initializes a test database with the required tables
func setupTestDB(t *testing.T, db *sql.DB, tablePrefix string) {
	t.Helper()

	// Create country table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + tablePrefix + `_country (
			id VARCHAR(40) PRIMARY KEY,
			status VARCHAR(20),
			iso2_code VARCHAR(2),
			iso3_code VARCHAR(3),
			name VARCHAR(255),
			continent VARCHAR(100),
			phone_prefix VARCHAR(20),
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create country table: %v", err)
	}

	// Create state table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + tablePrefix + `_state (
			id VARCHAR(40) PRIMARY KEY,
			status VARCHAR(20),
			country_code VARCHAR(2),
			state_code VARCHAR(5),
			name VARCHAR(255),
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create state table: %v", err)
	}

	// Create timezone table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + tablePrefix + `_timezone (
			id VARCHAR(40) PRIMARY KEY,
			country_code VARCHAR(2),
			zone_name VARCHAR(100),
			global_name VARCHAR(100),
			time_zone VARCHAR(100),
			gtm_offset VARCHAR(20),
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create timezone table: %v", err)
	}
}

func TestStoreCountryCreate(t *testing.T) {
	db := initDB()
	setupTestDB(t, db, "test_country_create")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "test_country_create_country",
		StateTableName:     "test_country_create_state",
		TimezoneTableName:  "test_country_create_timezone",
		AutomigrateEnabled: false, // We're handling table creation manually
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	// enable SQL debug logging for this test to help diagnose sqlite syntax errors
	store.EnableDebug(true)

	ctx := context.Background()

	country := NewCountry().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("Unknown")

	err = store.CountryCreate(ctx, country)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreCountryFindByID(t *testing.T) {
	db := initDB()
	setupTestDB(t, db, "test_country_find_by_id")

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

	ctx := context.Background()

	country := NewCountry().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("Unknown").
		SetIsoCode2("UN").
		SetIsoCode3("UNK")

	err = store.CountryCreate(ctx, country)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	countryFound, errFind := store.CountryFindByID(ctx, country.ID())

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
	db := initDB()
	setupTestDB(t, db, "test_country_soft_delete")

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

	ctx := context.Background()

	country := NewCountry().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("Unknown").
		SetIsoCode2("UN").
		SetIsoCode3("UNK")

	err = store.CountryCreate(ctx, country)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CountrySoftDeleteByID(ctx, country.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if country.DeletedAt() != sb.NULL_DATETIME {
		t.Fatal("Country MUST NOT be soft deleted")
	}

	countryFound, errFind := store.CountryFindByID(ctx, country.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if countryFound != nil {
		t.Fatal("Country MUST be nil")
	}

	countryFindWithDeleted, err := store.CountryList(ctx, CountryQueryOptions{
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
	db := initDB()
	setupTestDB(t, db, "test_state_create")

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

	ctx := context.Background()

	states, err := store.StateList(ctx, StateQueryOptions{})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(states) < 10 {
		t.Fatal("there must be states in the database")
	}
}

func TestStoreStateCreate(t *testing.T) {
	db := initDB()
	setupTestDB(t, db, "test_state_create2")

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
