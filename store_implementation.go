package geostore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/dracory/neat"
	contractsorm "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dracory/neat/database/schema/constants"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

var _ StoreInterface = (*storeImplementation)(nil) // verify it extends the interface

type storeImplementation struct {
	countryTableName   string
	stateTableName     string
	timezoneTableName  string
	db                 *neat.Database
	automigrateEnabled bool
	autoseedEnabled    bool
	debugEnabled       bool
	logger             *slog.Logger
}

// MigrateUp creates all database tables
func (store *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	// create country table
	if !store.db.Schema().HasTable(store.countryTableName) {
		err := store.db.Schema().Create(store.countryTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 21)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_STATUS, 20)
			table.String(COLUMN_ISO2_CODE, 2)
			table.String(COLUMN_ISO3_CODE, 3)
			table.String(COLUMN_NAME, 255)
			table.String(COLUMN_CONTINENT, 100)
			table.String(COLUMN_PHONE_PREFIX, 20)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT).Default(constants.MaxSoftDeletedAtDefault)
		})
		if err != nil {
			store.logger.Error("MigrateUp failed for country table", "error", err)
			return err
		}
	}

	// create state table
	if !store.db.Schema().HasTable(store.stateTableName) {
		err := store.db.Schema().Create(store.stateTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 21)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_STATUS, 20)
			table.String(COLUMN_COUNTRY_CODE, 2)
			table.String(COLUMN_STATE_CODE, 5)
			table.String(COLUMN_NAME, 255)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT).Default(constants.MaxSoftDeletedAtDefault)
		})
		if err != nil {
			store.logger.Error("MigrateUp failed for state table", "error", err)
			return err
		}
	}

	// create timezone table
	if !store.db.Schema().HasTable(store.timezoneTableName) {
		err := store.db.Schema().Create(store.timezoneTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 21)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_STATUS, 20)
			table.String(COLUMN_TIMEZONE, 100)
			table.String(COLUMN_ZONE_NAME, 100)
			table.String(COLUMN_GLOBAL_NAME, 100)
			table.String(COLUMN_COUNTRY_CODE, 50)
			table.String(COLUMN_OFFSET, 50)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT).Default(constants.MaxSoftDeletedAtDefault)
		})
		if err != nil {
			store.logger.Error("MigrateUp failed for timezone table", "error", err)
			return err
		}
	}

	return nil
}

// MigrateDown drops all database tables
func (store *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	// Drop tables in reverse order to avoid foreign key constraints
	if store.db.Schema().HasTable(store.timezoneTableName) {
		err := store.db.Schema().Drop(store.timezoneTableName)
		if err != nil {
			store.logger.Error("MigrateDown failed for timezone table", "error", err)
			return err
		}
	}

	if store.db.Schema().HasTable(store.stateTableName) {
		err := store.db.Schema().Drop(store.stateTableName)
		if err != nil {
			store.logger.Error("MigrateDown failed for state table", "error", err)
			return err
		}
	}

	if store.db.Schema().HasTable(store.countryTableName) {
		err := store.db.Schema().Drop(store.countryTableName)
		if err != nil {
			store.logger.Error("MigrateDown failed for country table", "error", err)
			return err
		}
	}

	return nil
}

// Seed populates all tables with initial data
func (store *storeImplementation) Seed(ctx context.Context, tx ...*sql.Tx) error {
	// seed country table
	err := store.seedCountriesIfTableIsEmpty(ctx)
	if err != nil {
		store.logger.Error("Seed failed for countries", "error", err)
		return err
	}

	// seed state table
	err = store.seedStatesIfTableEmpty(ctx)
	if err != nil {
		store.logger.Error("Seed failed for states", "error", err)
		return err
	}

	// seed timezone table
	err = store.seedTimezonesIfTableEmpty(ctx)
	if err != nil {
		store.logger.Error("Seed failed for timezones", "error", err)
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (store *storeImplementation) EnableDebug(debug bool) {
	store.debugEnabled = debug
	if debug {
		store.db.EnableDebug()
		store.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		store.db.DisableDebug()
		store.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
}

// GetCountryTableName returns the country table name
func (store *storeImplementation) GetCountryTableName() string {
	return store.countryTableName
}

// SetCountryTableName sets the country table name
func (store *storeImplementation) SetCountryTableName(countryTableName string) {
	store.countryTableName = countryTableName
}

// GetStateTableName returns the state table name
func (store *storeImplementation) GetStateTableName() string {
	return store.stateTableName
}

// SetStateTableName sets the state table name
func (store *storeImplementation) SetStateTableName(stateTableName string) {
	store.stateTableName = stateTableName
}

// GetTimezoneTableName returns the timezone table name
func (store *storeImplementation) GetTimezoneTableName() string {
	return store.timezoneTableName
}

// SetTimezoneTableName sets the timezone table name
func (store *storeImplementation) SetTimezoneTableName(timezoneTableName string) {
	store.timezoneTableName = timezoneTableName
}

// == COUNTRY CRUD ==========================================================

func (store *storeImplementation) CountryCreate(ctx context.Context, country *Country) error {
	country.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	country.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := store.countryToMap(country)

	err := store.db.Query().Table(store.countryTableName).Create(data)
	if err != nil {
		return err
	}

	country.MarkAsNotDirty()
	return nil
}

func (store *storeImplementation) CountryDelete(ctx context.Context, country *Country) error {
	if country == nil {
		return errors.New("country is nil")
	}

	return store.CountryDeleteByID(ctx, country.ID())
}

func (store *storeImplementation) CountryDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("country id is empty")
	}

	_, err := store.db.Query().
		Table(store.countryTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()

	return err
}

func (store *storeImplementation) CountryFindByID(ctx context.Context, id string) (*Country, error) {
	if id == "" {
		return nil, errors.New("country id is empty")
	}

	list, err := store.CountryList(ctx, CountryQueryOptions{
		ID:    id,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

func (store *storeImplementation) CountryFindByIso2(ctx context.Context, iso2Code string) (*Country, error) {
	if iso2Code == "" {
		return nil, errors.New("country iso2 code is empty")
	}

	list, err := store.CountryList(ctx, CountryQueryOptions{
		Status: COUNTRY_STATUS_ACTIVE,
		Iso2:   iso2Code,
		Limit:  1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

func (store *storeImplementation) CountryNameFindByIso2(ctx context.Context, iso2Code string) (string, error) {
	country, err := store.CountryFindByIso2(ctx, iso2Code)

	if err != nil {
		return "", err
	}

	if country == nil {
		return "", nil
	}

	return country.Name(), nil
}

func (store *storeImplementation) CountryList(ctx context.Context, options CountryQueryOptions) ([]Country, error) {
	q := store.countryQuery(options)

	var results []map[string]any
	err := q.Get(&results)
	if err != nil {
		return []Country{}, err
	}

	list := []Country{}
	lo.ForEach(results, func(result map[string]any, index int) {
		model := store.mapToCountry(result)
		list = append(list, *model)
	})

	return list, nil
}

func (store *storeImplementation) CountrySoftDelete(ctx context.Context, country *Country) error {
	if country == nil {
		return errors.New("country is nil")
	}

	country.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.CountryUpdate(ctx, country)
}

func (store *storeImplementation) CountrySoftDeleteByID(ctx context.Context, id string) error {
	country, err := store.CountryFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.CountrySoftDelete(ctx, country)
}

func (store *storeImplementation) CountryUpdate(ctx context.Context, country *Country) error {
	if country == nil {
		return errors.New("country is nil")
	}

	data := store.countryToMap(country)
	delete(data, COLUMN_ID) // ID is not updatable

	// Check if any meaningful field has changed
	if country.originalData != nil {
		hasChanges := false
		for k, v := range data {
			if k == COLUMN_ID || k == COLUMN_CREATED_AT || k == COLUMN_UPDATED_AT {
				continue
			}
			if country.originalData[k] != v {
				hasChanges = true
				break
			}
		}
		if !hasChanges {
			return nil
		}
	}

	country.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	data = store.countryToMap(country)
	delete(data, COLUMN_ID) // ID is not updatable

	_, err := store.db.Query().
		Table(store.countryTableName).
		Where(COLUMN_ID+" = ?", country.ID()).
		Update(data)

	if err != nil {
		return err
	}

	country.MarkAsNotDirty()
	return nil
}

func (store *storeImplementation) countryQuery(options CountryQueryOptions) contractsorm.Query {
	q := store.db.Query().Table(store.countryTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if options.Status != "" {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status)
	}

	if options.Iso2 != "" {
		q = q.Where(COLUMN_ISO2_CODE+" = ?", options.Iso2)
	}

	if options.Iso3 != "" {
		q = q.Where(COLUMN_ISO3_CODE+" = ?", options.Iso3)
	}

	if len(options.IDIn) > 0 {
		q = q.WhereIn(COLUMN_ID, lo.ToAnySlice(options.IDIn))
	}

	if len(options.StatusIn) > 0 {
		q = q.WhereIn(COLUMN_STATUS, lo.ToAnySlice(options.StatusIn))
	}

	if options.Limit > 0 {
		q = q.Limit(options.Limit)
	}

	if options.Offset > 0 {
		q = q.Offset(options.Offset)
	}

	if options.OrderBy != "" {
		direction := options.SortOrder
		if direction == "" {
			direction = "asc"
		}
		q = q.OrderBy(options.OrderBy, direction)
	}

	if !options.WithDeleted {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", MAX_DATETIME)
	}

	return q
}

// == STATE CRUD ============================================================

func (store *storeImplementation) StateCreate(ctx context.Context, state *State) error {
	state.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	state.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := store.stateToMap(state)

	err := store.db.Query().Table(store.stateTableName).Create(data)
	if err != nil {
		return err
	}

	state.MarkAsNotDirty()
	return nil
}

func (store *storeImplementation) StatesCreate(ctx context.Context, states []*State) error {
	now := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	for _, state := range states {
		state.SetCreatedAt(now)
		state.SetUpdatedAt(now)
	}

	data := make([]map[string]any, len(states))
	for i, state := range states {
		data[i] = store.stateToMap(state)
	}

	const batchSize = 100
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		err := store.db.Query().Table(store.stateTableName).Create(data[i:end])
		if err != nil {
			return err
		}
	}

	for _, state := range states {
		state.MarkAsNotDirty()
	}

	return nil
}

func (store *storeImplementation) StateList(ctx context.Context, options StateQueryOptions) ([]State, error) {
	q := store.stateQuery(options)

	var results []map[string]any
	err := q.Get(&results)
	if err != nil {
		return []State{}, err
	}

	list := []State{}
	lo.ForEach(results, func(result map[string]any, index int) {
		model := store.mapToState(result)
		list = append(list, *model)
	})

	return list, nil
}

func (store *storeImplementation) stateQuery(options StateQueryOptions) contractsorm.Query {
	q := store.db.Query().Table(store.stateTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if options.Status != "" {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status)
	}

	if options.CountryCode != "" {
		q = q.Where(COLUMN_COUNTRY_CODE+" = ?", options.CountryCode)
	}

	if options.StateCode != "" {
		q = q.Where(COLUMN_STATE_CODE+" = ?", options.StateCode)
	}

	if len(options.StatusIn) > 0 {
		q = q.WhereIn(COLUMN_STATUS, lo.ToAnySlice(options.StatusIn))
	}

	if options.Limit > 0 {
		q = q.Limit(options.Limit)
	}

	if options.Offset > 0 {
		q = q.Offset(options.Offset)
	}

	if options.OrderBy != "" {
		direction := options.SortOrder
		if direction == "" {
			direction = "asc"
		}
		q = q.OrderBy(options.OrderBy, direction)
	}

	if !options.WithDeleted {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", MAX_DATETIME)
	}

	return q
}

// == TIMEZONE CRUD =========================================================

func (store *storeImplementation) TimezoneCreate(ctx context.Context, timezone *Timezone) error {
	timezone.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	timezone.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := store.timezoneToMap(timezone)

	err := store.db.Query().Table(store.timezoneTableName).Create(data)
	if err != nil {
		return err
	}

	timezone.MarkAsNotDirty()
	return nil
}

func (store *storeImplementation) TimezoneList(ctx context.Context, options TimezoneQueryOptions) ([]Timezone, error) {
	q := store.timezoneQuery(options)

	var results []map[string]any
	err := q.Get(&results)
	if err != nil {
		return []Timezone{}, err
	}

	list := []Timezone{}
	lo.ForEach(results, func(result map[string]any, index int) {
		model := store.mapToTimezone(result)
		list = append(list, *model)
	})

	return list, nil
}

func (store *storeImplementation) timezoneQuery(options TimezoneQueryOptions) contractsorm.Query {
	q := store.db.Query().Table(store.timezoneTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if options.Status != "" {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status)
	}

	if options.CountryCode != "" {
		q = q.Where(COLUMN_COUNTRY_CODE+" = ?", options.CountryCode)
	}

	if options.Timezone != "" {
		q = q.Where(COLUMN_TIMEZONE+" = ?", options.Timezone)
	}

	if len(options.StatusIn) > 0 {
		q = q.WhereIn(COLUMN_STATUS, lo.ToAnySlice(options.StatusIn))
	}

	if options.Limit > 0 {
		q = q.Limit(options.Limit)
	}

	if options.Offset > 0 {
		q = q.Offset(options.Offset)
	}

	if options.OrderBy != "" {
		direction := options.SortOrder
		if direction == "" {
			direction = "asc"
		}
		q = q.OrderBy(options.OrderBy, direction)
	}

	if !options.WithDeleted {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", MAX_DATETIME)
	}

	return q
}

// == HELPERS ===============================================================

func (store *storeImplementation) countryToMap(country *Country) map[string]any {
	return map[string]any{
		COLUMN_ID:              country.ID(),
		COLUMN_STATUS:          country.Status(),
		COLUMN_ISO2_CODE:       country.IsoCode2(),
		COLUMN_ISO3_CODE:       country.IsoCode3(),
		COLUMN_NAME:            country.Name(),
		COLUMN_CONTINENT:       country.Continent(),
		COLUMN_PHONE_PREFIX:    country.PhonePrefix(),
		COLUMN_CREATED_AT:      country.CreatedAt(),
		COLUMN_UPDATED_AT:      country.UpdatedAt(),
		COLUMN_SOFT_DELETED_AT: country.GetSoftDeletedAt(),
	}
}

func (store *storeImplementation) stateToMap(state *State) map[string]any {
	return map[string]any{
		COLUMN_ID:              state.ID(),
		COLUMN_STATUS:          state.Status(),
		COLUMN_COUNTRY_CODE:    state.CountryCode(),
		COLUMN_STATE_CODE:      state.StateCode(),
		COLUMN_NAME:            state.Name(),
		COLUMN_CREATED_AT:      state.CreatedAt(),
		COLUMN_UPDATED_AT:      state.UpdatedAt(),
		COLUMN_SOFT_DELETED_AT: state.GetSoftDeletedAt(),
	}
}

func (store *storeImplementation) timezoneToMap(timezone *Timezone) map[string]any {
	return map[string]any{
		COLUMN_ID:              timezone.ID(),
		COLUMN_STATUS:          timezone.Status(),
		COLUMN_TIMEZONE:        timezone.Timezone(),
		COLUMN_ZONE_NAME:       timezone.ZoneName(),
		COLUMN_GLOBAL_NAME:     timezone.GlobalName(),
		COLUMN_COUNTRY_CODE:    timezone.CountryCode(),
		COLUMN_OFFSET:          timezone.Offset(),
		COLUMN_CREATED_AT:      timezone.CreatedAt(),
		COLUMN_UPDATED_AT:      timezone.UpdatedAt(),
		COLUMN_SOFT_DELETED_AT: timezone.GetSoftDeletedAt(),
	}
}

func (store *storeImplementation) mapToCountry(data map[string]any) *Country {
	stringData := make(map[string]string)
	for k, v := range data {
		if v != nil {
			stringData[k] = toString(v)
		} else {
			stringData[k] = ""
		}
	}
	return NewCountryFromExistingData(stringData)
}

func (store *storeImplementation) mapToState(data map[string]any) *State {
	stringData := make(map[string]string)
	for k, v := range data {
		if v != nil {
			stringData[k] = toString(v)
		} else {
			stringData[k] = ""
		}
	}
	return NewStateFromExistingData(stringData)
}

func (store *storeImplementation) mapToTimezone(data map[string]any) *Timezone {
	stringData := make(map[string]string)
	for k, v := range data {
		if v != nil {
			stringData[k] = toString(v)
		} else {
			stringData[k] = ""
		}
	}
	return NewTimezoneFromExistingData(stringData)
}

func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case time.Time:
		if val.IsZero() {
			return ""
		}
		return carbon.CreateFromStdTime(val).ToDateTimeString()
	case []byte:
		return string(val)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%f", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
