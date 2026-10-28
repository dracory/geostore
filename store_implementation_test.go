package geostore

import (
	"context"
	"database/sql"
	"strings"
	"testing"

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

// setupTestDB initializes a test database using the sb package for table creation
// seed parameter controls whether to populate tables with initial data
func setupTestDB(t *testing.T, store StoreInterface, seed bool) {
	t.Helper()

	err := store.MigrateUp(context.Background())

	if err != nil {
		t.Fatal("MigrateUp failed: " + err.Error())
	}

	if seed {
		err = store.Seed(context.Background())
		if err != nil {
			t.Fatalf("Failed to seed tables: %v", err)
		}
	}
}

func TestStoreCountryCreate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "test_country_create",
		StateTableName:     "test_country_create_state",
		TimezoneTableName:  "test_country_create_timezone",
		AutomigrateEnabled: false,
		AutoseedEnabled:    false,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	setupTestDB(t, store, true)

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

	setupTestDB(t, store, true)

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

	if !strings.Contains(countryFound.GetSoftDeletedAt(), MAX_DATETIME) {
		t.Fatal("Country MUST NOT be soft deleted", countryFound.GetSoftDeletedAt())
	}
}

func TestStoreCountrySoftDelete(t *testing.T) {
	db := initDB()

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

	setupTestDB(t, store, true)

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

	if country.GetSoftDeletedAt() != MAX_DATETIME {
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

	if strings.Contains(countryFindWithDeleted[0].GetSoftDeletedAt(), MAX_DATETIME) {
		t.Fatal("Exam MUST be soft deleted", country.GetSoftDeletedAt())
	}

}

func TestStoreStateAutomigrate(t *testing.T) {
	db := initDB()

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

	setupTestDB(t, store, true)

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

	setupTestDB(t, store, true)

	state := NewState().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("Unknown").
		SetStateCode("UN")

	err = store.StateCreate(context.Background(), state)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreTimezoneCreate(t *testing.T) {
	db := initDB()

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

	setupTestDB(t, store, true)

	ctx := context.Background()

	timezone := NewTimezone().
		SetStatus("active").
		SetTimezone("America/New_York").
		SetZoneName("Eastern").
		SetGlobalName("Eastern Time").
		SetCountryCode("US").
		SetOffset("-05:00")

	err = store.TimezoneCreate(ctx, timezone)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreTimezoneList(t *testing.T) {
	db := initDB()

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

	setupTestDB(t, store, true)

	ctx := context.Background()

	// Create test timezone
	timezone := NewTimezone().
		SetStatus("active").
		SetTimezone("America/New_York").
		SetZoneName("Eastern").
		SetGlobalName("Eastern Time").
		SetCountryCode("US").
		SetOffset("-05:00")

	err = store.TimezoneCreate(ctx, timezone)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Test listing all timezones
	timezones, err := store.TimezoneList(ctx, TimezoneQueryOptions{})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(timezones) == 0 {
		t.Fatal("there must be at least one timezone")
	}

	// Test filtering by country code
	timezonesUS, err := store.TimezoneList(ctx, TimezoneQueryOptions{
		CountryCode: "US",
		Limit:       10,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(timezonesUS) == 0 {
		t.Fatal("there must be at least one US timezone")
	}

	// Test filtering by timezone name
	timezonesNY, err := store.TimezoneList(ctx, TimezoneQueryOptions{
		Timezone: "America/New_York",
		Limit:    1,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(timezonesNY) == 0 {
		t.Fatal("there must be at least one America/New_York timezone")
	}

	if timezonesNY[0].Timezone() != "America/New_York" {
		t.Fatal("timezone name must match", timezonesNY[0].Timezone())
	}
}

func TestStoreCountryUpdate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "geo_country_update",
		StateTableName:     "geo_state_update",
		TimezoneTableName:  "geo_timezone_update",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	setupTestDB(t, store, true)

	ctx := context.Background()

	// Create country
	country := NewCountry().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("Testland").
		SetIsoCode2("TL").
		SetIsoCode3("TLD")

	err = store.CountryCreate(ctx, country)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Update country
	country.SetName("Updated Testland")
	country.SetPhonePrefix("+123")

	err = store.CountryUpdate(ctx, country)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Verify update
	updatedCountry, err := store.CountryFindByID(ctx, country.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if updatedCountry.Name() != "Updated Testland" {
		t.Fatal("country name should be updated", updatedCountry.Name())
	}

	if updatedCountry.PhonePrefix() != "+123" {
		t.Fatal("phone prefix should be updated", updatedCountry.PhonePrefix())
	}
}

func TestStoreCountryFindByIso2(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "geo_country_iso2",
		StateTableName:     "geo_state_iso2",
		TimezoneTableName:  "geo_timezone_iso2",
		AutomigrateEnabled: false,
		AutoseedEnabled:    false,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	setupTestDB(t, store, false)

	ctx := context.Background()

	// Create country
	country := NewCountry().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("IsoTest").
		SetIsoCode2("IT").
		SetIsoCode3("IST")

	err = store.CountryCreate(ctx, country)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Find by ISO2
	foundCountry, err := store.CountryFindByIso2(ctx, "IT")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if foundCountry == nil {
		t.Fatal("country should be found")
	}

	if foundCountry.IsoCode2() != "IT" {
		t.Fatal("ISO2 code should match", foundCountry.IsoCode2())
	}

	if foundCountry.Name() != "IsoTest" {
		t.Fatal("name should match IsoTest", foundCountry.Name())
	}
}

func TestStoreCountryDelete(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "geo_country_delete",
		StateTableName:     "geo_state_delete",
		TimezoneTableName:  "geo_timezone_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	setupTestDB(t, store, true)

	ctx := context.Background()

	// Create country
	country := NewCountry().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("DeleteTest").
		SetIsoCode2("DT").
		SetIsoCode3("DLT")

	err = store.CountryCreate(ctx, country)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	countryID := country.ID()

	// Delete country
	err = store.CountryDelete(ctx, country)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Verify deletion
	foundCountry, err := store.CountryFindByID(ctx, countryID)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if foundCountry != nil {
		t.Fatal("country should be deleted")
	}
}

func TestStoreCountryDeleteByID(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "geo_country_delete_id",
		StateTableName:     "geo_state_delete_id",
		TimezoneTableName:  "geo_timezone_delete_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	setupTestDB(t, store, true)

	ctx := context.Background()

	// Create country
	country := NewCountry().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("DeleteIDTest").
		SetIsoCode2("DI").
		SetIsoCode3("DID")

	err = store.CountryCreate(ctx, country)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	countryID := country.ID()

	// Delete country by ID
	err = store.CountryDeleteByID(ctx, countryID)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Verify deletion
	foundCountry, err := store.CountryFindByID(ctx, countryID)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if foundCountry != nil {
		t.Fatal("country should be deleted")
	}
}

func TestStoreCountryQueryOptions(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "geo_country_query",
		StateTableName:     "geo_state_query",
		TimezoneTableName:  "geo_timezone_query",
		AutomigrateEnabled: false,
		AutoseedEnabled:    false,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	setupTestDB(t, store, false)

	ctx := context.Background()

	// Create test countries
	country1 := NewCountry().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("ActiveCountry").
		SetIsoCode2("AC").
		SetIsoCode3("ACT").
		SetContinent("Asia")

	country2 := NewCountry().
		SetStatus("inactive").
		SetName("InactiveCountry").
		SetIsoCode2("IC").
		SetIsoCode3("ICT").
		SetContinent("Europe")

	err = store.CountryCreate(ctx, country1)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CountryCreate(ctx, country2)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Test filter by status
	activeCountries, err := store.CountryList(ctx, CountryQueryOptions{
		Status: COUNTRY_STATUS_ACTIVE,
		Limit:  10,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(activeCountries) != 1 {
		t.Fatal("should find exactly 1 active country", len(activeCountries))
	}

	if activeCountries[0].Status() != COUNTRY_STATUS_ACTIVE {
		t.Fatal("country should be active", activeCountries[0].Status())
	}

	// Test filter by ISO2 code
	asiaCountries, err := store.CountryList(ctx, CountryQueryOptions{
		Iso2:  "AC",
		Limit: 10,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(asiaCountries) != 1 {
		t.Fatal("should find exactly 1 country with ISO2 AC", len(asiaCountries))
	}

	if asiaCountries[0].IsoCode2() != "AC" {
		t.Fatal("country should have ISO2 code AC", asiaCountries[0].IsoCode2())
	}

	// Test limit
	allCountries, err := store.CountryList(ctx, CountryQueryOptions{
		Limit: 1,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(allCountries) != 1 {
		t.Fatal("should limit to 1 country", len(allCountries))
	}
}

func TestStoreStateQueryOptions(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "geo_country_state_query",
		StateTableName:     "geo_state_state_query",
		TimezoneTableName:  "geo_timezone_state_query",
		AutomigrateEnabled: false,
		AutoseedEnabled:    false,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	setupTestDB(t, store, false)

	ctx := context.Background()

	// Create test states
	state1 := NewState().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("ActiveState").
		SetStateCode("AS").
		SetCountryCode("US")

	state2 := NewState().
		SetStatus("inactive").
		SetName("InactiveState").
		SetStateCode("IS").
		SetCountryCode("CA")

	err = store.StateCreate(context.Background(), state1)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.StateCreate(context.Background(), state2)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Test filter by country code
	usStates, err := store.StateList(ctx, StateQueryOptions{
		CountryCode: "US",
		Limit:       10,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(usStates) != 1 {
		t.Fatal("should find exactly 1 US state", len(usStates))
	}

	if usStates[0].CountryCode() != "US" {
		t.Fatal("state should be in US", usStates[0].CountryCode())
	}

	// Test filter by status
	activeStates, err := store.StateList(ctx, StateQueryOptions{
		Status: COUNTRY_STATUS_ACTIVE,
		Limit:  10,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(activeStates) != 1 {
		t.Fatal("should find exactly 1 active state", len(activeStates))
	}

	if activeStates[0].Status() != COUNTRY_STATUS_ACTIVE {
		t.Fatal("state should be active", activeStates[0].Status())
	}
}

func TestStoreStatesCreate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "geo_country_bulk",
		StateTableName:     "geo_state_bulk",
		TimezoneTableName:  "geo_timezone_bulk",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	setupTestDB(t, store, true)

	// Create multiple states
	states := []*State{
		NewState().SetStatus(COUNTRY_STATUS_ACTIVE).SetName("State1").SetStateCode("S1").SetCountryCode("US"),
		NewState().SetStatus(COUNTRY_STATUS_ACTIVE).SetName("State2").SetStateCode("S2").SetCountryCode("US"),
		NewState().SetStatus(COUNTRY_STATUS_ACTIVE).SetName("State3").SetStateCode("S3").SetCountryCode("CA"),
	}

	err = store.StatesCreate(context.Background(), states)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Verify all states were created
	ctx := context.Background()
	allStates, err := store.StateList(ctx, StateQueryOptions{
		Limit: 10,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(allStates) < 3 {
		t.Fatal("should find at least 3 states", len(allStates))
	}
}

func TestStoreMigrateUp(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "test_migrate_up_country",
		StateTableName:     "test_migrate_up_state",
		TimezoneTableName:  "test_migrate_up_timezone",
		AutomigrateEnabled: false,
		AutoseedEnabled:    false,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	// Test MigrateUp creates tables
	err = store.MigrateUp(context.Background())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Verify tables exist by checking we can query them
	ctx := context.Background()

	// Check country table exists
	countries, err := store.CountryList(ctx, CountryQueryOptions{Limit: 1})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if countries == nil {
		t.Fatal("countries should not be nil")
	}

	// Check state table exists
	states, err := store.StateList(ctx, StateQueryOptions{Limit: 1})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if states == nil {
		t.Fatal("states should not be nil")
	}

	// Check timezone table exists
	timezones, err := store.TimezoneList(ctx, TimezoneQueryOptions{Limit: 1})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if timezones == nil {
		t.Fatal("timezones should not be nil")
	}
}

func TestStoreMigrateDown(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "test_migrate_down_country",
		StateTableName:     "test_migrate_down_state",
		TimezoneTableName:  "test_migrate_down_timezone",
		AutomigrateEnabled: false,
		AutoseedEnabled:    false,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	// First migrate up to create tables
	err = store.MigrateUp(context.Background())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Verify tables exist
	ctx := context.Background()
	_, err = store.CountryList(ctx, CountryQueryOptions{Limit: 1})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Test MigrateDown drops tables
	err = store.MigrateDown(context.Background())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Verify tables are gone - should get errors when querying
	_, err = store.CountryList(ctx, CountryQueryOptions{Limit: 1})
	if err == nil {
		t.Fatal("expected error when querying dropped table")
	}
}

func TestStoreSeed(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "test_seed_country",
		StateTableName:     "test_seed_state",
		TimezoneTableName:  "test_seed_timezone",
		AutomigrateEnabled: false,
		AutoseedEnabled:    false,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	// First migrate up to create tables
	err = store.MigrateUp(context.Background())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Test Seed populates data
	err = store.Seed(context.Background())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Verify seeded data exists
	ctx := context.Background()

	// Check countries are seeded
	countries, err := store.CountryList(ctx, CountryQueryOptions{})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(countries) == 0 {
		t.Fatal("countries should be seeded")
	}

	// Check states are seeded
	states, err := store.StateList(ctx, StateQueryOptions{})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(states) == 0 {
		t.Fatal("states should be seeded")
	}

	// Check timezones are seeded
	timezones, err := store.TimezoneList(ctx, TimezoneQueryOptions{})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(timezones) == 0 {
		t.Fatal("timezones should be seeded")
	}
}

func TestStoreAutomigrate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "test_automigrate_country",
		StateTableName:     "test_automigrate_state",
		TimezoneTableName:  "test_automigrate_timezone",
		AutomigrateEnabled: false,
		AutoseedEnabled:    false,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	// Test MigrateUp only migrates (doesn't seed)
	err = store.MigrateUp(context.Background())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Verify tables exist but are empty
	ctx := context.Background()

	countries, err := store.CountryList(ctx, CountryQueryOptions{})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(countries) != 0 {
		t.Fatal("countries should be empty after Automigrate only", len(countries))
	}
}

func TestStoreAutoseed(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "test_autoseed_country",
		StateTableName:     "test_autoseed_state",
		TimezoneTableName:  "test_autoseed_timezone",
		AutomigrateEnabled: false,
		AutoseedEnabled:    false,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	// First migrate up to create tables
	err = store.MigrateUp(context.Background())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Test Seed only seeds (doesn't migrate)
	err = store.Seed(context.Background())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Verify data is seeded
	ctx := context.Background()

	countries, err := store.CountryList(ctx, CountryQueryOptions{})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(countries) == 0 {
		t.Fatal("countries should be seeded after Autoseed")
	}
}

func TestStoreAutoMigrateBackwardCompatibility(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "test_auto_migrate_compat_country",
		StateTableName:     "test_auto_migrate_compat_state",
		TimezoneTableName:  "test_auto_migrate_compat_timezone",
		AutomigrateEnabled: false,
		AutoseedEnabled:    false,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	// Test MigrateUp + Seed maintains backward compatibility (migrates + seeds)
	err = store.MigrateUp(context.Background())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.Seed(context.Background())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Verify both tables exist and data is seeded
	ctx := context.Background()

	countries, err := store.CountryList(ctx, CountryQueryOptions{})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(countries) == 0 {
		t.Fatal("countries should be seeded after AutoMigrate")
	}

	states, err := store.StateList(ctx, StateQueryOptions{})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(states) == 0 {
		t.Fatal("states should be seeded after AutoMigrate")
	}
}

func TestCountryList_CacheServesAndInvalidates(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "test_country_cache",
		StateTableName:     "test_country_cache_state",
		TimezoneTableName:  "test_country_cache_timezone",
		AutomigrateEnabled: false,
		AutoseedEnabled:    false,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	setupTestDB(t, store, true)
	ctx := context.Background()

	// warm the cache
	country, err := store.CountryFindByIso2(ctx, "US")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if country == nil {
		t.Fatal("expected seeded US country")
	}

	// update must invalidate the snapshot
	country.SetName("United States of America")
	if err := store.CountryUpdate(ctx, country); err != nil {
		t.Fatal("unexpected error:", err)
	}
	updated, err := store.CountryFindByIso2(ctx, "US")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if updated == nil || updated.Name() != "United States of America" {
		t.Fatal("expected updated name after CountryUpdate")
	}

	// soft delete must invalidate and hide the row from default queries
	if err := store.CountrySoftDelete(ctx, updated); err != nil {
		t.Fatal("unexpected error:", err)
	}
	deleted, err := store.CountryFindByIso2(ctx, "US")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if deleted != nil {
		t.Fatal("expected soft-deleted country to be hidden")
	}
	withDeleted, err := store.CountryList(ctx, CountryQueryOptions{Iso2: "US", WithDeleted: true})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(withDeleted) != 1 {
		t.Fatal("expected soft-deleted country with WithDeleted option")
	}

	// create must invalidate the snapshot
	newCountry := NewCountry()
	newCountry.SetIsoCode2("XX")
	newCountry.SetIsoCode3("XXX")
	newCountry.SetName("Testland")
	newCountry.SetStatus(COUNTRY_STATUS_ACTIVE)
	if err := store.CountryCreate(ctx, newCountry); err != nil {
		t.Fatal("unexpected error:", err)
	}
	created, err := store.CountryFindByIso2(ctx, "XX")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if created == nil || created.Name() != "Testland" {
		t.Fatal("expected newly created country to be visible")
	}
}

func TestMigrateUp_IDColumnLength(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		CountryTableName:   "test_id_length_country",
		StateTableName:     "test_id_length_state",
		TimezoneTableName:  "test_id_length_timezone",
		AutomigrateEnabled: false,
		AutoseedEnabled:    false,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.MigrateUp(ctx)
	if err != nil {
		t.Fatal("MigrateUp failed:", err)
	}

	// 36-char UUIDv7 string (e.g., 018f3a5b-7c12-7000-8000-000000000000)
	country := NewCountry().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("36CharIDCountry")
	if len(country.ID()) != 36 {
		t.Fatalf("expected country ID length 36, got %d (%s)", len(country.ID()), country.ID())
	}
	err = store.CountryCreate(ctx, country)
	if err != nil {
		t.Fatal("CountryCreate failed with 36-char ID:", err)
	}

	state := NewState().
		SetStatus(COUNTRY_STATUS_ACTIVE).
		SetName("36CharIDState").
		SetStateCode("36")
	if len(state.ID()) != 36 {
		t.Fatalf("expected state ID length 36, got %d (%s)", len(state.ID()), state.ID())
	}
	err = store.StateCreate(ctx, state)
	if err != nil {
		t.Fatal("StateCreate failed with 36-char ID:", err)
	}

	timezone := NewTimezone().
		SetStatus("active").
		SetTimezone("Test/36CharID").
		SetZoneName("Test").
		SetGlobalName("Test Time").
		SetCountryCode("TS").
		SetOffset("+00:00")
	if len(timezone.ID()) != 36 {
		t.Fatalf("expected timezone ID length 36, got %d (%s)", len(timezone.ID()), timezone.ID())
	}
	err = store.TimezoneCreate(ctx, timezone)
	if err != nil {
		t.Fatal("TimezoneCreate failed with 36-char ID:", err)
	}

	// Verify retrieval
	foundCountry, err := store.CountryFindByID(ctx, country.ID())
	if err != nil || foundCountry == nil {
		t.Fatal("failed to find country by 36-char ID:", err)
	}

	foundStates, err := store.StateList(ctx, StateQueryOptions{ID: state.ID()})
	if err != nil || len(foundStates) != 1 {
		t.Fatal("failed to find state by 36-char ID:", err)
	}

	foundTimezones, err := store.TimezoneList(ctx, TimezoneQueryOptions{ID: timezone.ID()})
	if err != nil || len(foundTimezones) != 1 {
		t.Fatal("failed to find timezone by 36-char ID:", err)
	}
}
