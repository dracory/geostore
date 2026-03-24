package geostore

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

// const DISCOUNT_TABLE_NAME = "shop_country"

var _ StoreInterface = (*storeImplementation)(nil) // verify it extends the interface

type storeImplementation struct {
	countryTableName   string
	stateTableName     string
	timezoneTableName  string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	autoseedEnabled    bool
	debugEnabled       bool
}

// MigrateUp creates all database tables
func (store *storeImplementation) MigrateUp() error {
	// create country table
	sql, err := store.sqlCountryTableCreate()
	if err != nil {
		log.Println(err)
		return err
	}

	if sql == "" {
		return errors.New("country table create sql is empty")
	}

	_, err = store.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	// create state table
	sql, err = store.sqlStateTableCreate()
	if err != nil {
		log.Println(err)
		return err
	}

	if sql == "" {
		return errors.New("state table create sql is empty")
	}

	_, err = store.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	// create timezone table
	sql, err = store.sqlTimezoneTableCreate()
	if err != nil {
		log.Println(err)
		return err
	}

	if sql == "" {
		return errors.New("timezone table create sql is empty")
	}

	_, err = store.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// MigrateDown drops all database tables
func (store *storeImplementation) MigrateDown() error {
	// Drop tables in reverse order to avoid foreign key constraints
	tables := []string{store.timezoneTableName, store.stateTableName, store.countryTableName}

	for _, table := range tables {
		_, err := store.db.Exec("DROP TABLE IF EXISTS " + table)
		if err != nil {
			log.Printf("Error dropping table %s: %v", table, err)
			return err
		}
	}

	return nil
}

// Seed populates all tables with initial data
func (store *storeImplementation) Seed() error {
	// seed country table
	err := store.seedCountriesIfTableIsEmpty()
	if err != nil {
		log.Println(err)
		return err
	}

	// seed state table
	err = store.seedStatesIfTableEmpty()
	if err != nil {
		log.Println(err)
		return err
	}

	// seed timezone table
	err = store.seedTimezonesIfTableEmpty()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Automigrate is a convenience method that calls MigrateUp
func (store *storeImplementation) Automigrate() error {
	return store.MigrateUp()
}

// Autoseed is a convenience method that calls Seed
func (store *storeImplementation) Autoseed() error {
	return store.Seed()
}

// AutoMigrate maintains backward compatibility - migrates and seeds
func (store *storeImplementation) AutoMigrate() error {
	err := store.MigrateUp()
	if err != nil {
		return err
	}

	return store.Seed()
}

// EnableDebug - enables the debug option
func (store *storeImplementation) EnableDebug(debug bool) {
	store.debugEnabled = debug
}

func (store *storeImplementation) CountryCreate(ctx context.Context, country *Country) error {
	country.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	country.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := country.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.countryTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

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

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.countryTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

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

	sqlStr, params, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []Country{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(database.NewQueryableContext(ctx, store.db), sqlStr, params...)
	if err != nil {
		return []Country{}, err
	}

	list := []Country{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewCountryFromExistingData(modelMap)
		list = append(list, *model)
	})

	return list, nil
}

func (store *storeImplementation) CountrySoftDelete(ctx context.Context, country *Country) error {
	if country == nil {
		return errors.New("country is nil")
	}

	country.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

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

	// country.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := country.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updatable
	delete(dataChanged, "hash")    // Hash is not updatable
	delete(dataChanged, "data")    // Data is not updatable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.countryTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(country.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	country.MarkAsNotDirty()

	return err
}

func (store *storeImplementation) countryQuery(options CountryQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.countryTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.Status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	}

	if options.Iso2 != "" {
		q = q.Where(goqu.C(COLUMN_ISO2_CODE).Eq(options.Iso2))
	}

	if options.Iso3 != "" {
		q = q.Where(goqu.C(COLUMN_ISO3_CODE).Eq(options.Iso3))
	}

	if len(options.IDIn) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
	}

	if options.CountOnly {
		q = q.Select(goqu.COUNT("*"))
	}

	if options.Limit > 0 {
		q = q.Limit(uint(options.Limit))
	}

	if options.Offset > 0 {
		q = q.Offset(uint(options.Offset))
	}

	if options.OrderBy != "" {
		if options.SortOrder == "desc" {
			q = q.Order(goqu.C(options.OrderBy).Desc())
		} else {
			q = q.Order(goqu.C(options.OrderBy).Asc())
		}
	}

	if !options.WithDeleted {
		q = q.Where(goqu.C(COLUMN_DELETED_AT).Eq(sb.NULL_DATETIME))
	}

	return q
}

func (store *storeImplementation) StateCreate(state *State) error {
	state.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	state.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := state.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.stateTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	state.MarkAsNotDirty()

	return nil
}

func (store *storeImplementation) StatesCreate(states []*State) error {
	for index, state := range states {
		state.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
		state.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
		states[index] = state
	}

	const batchSize = 50

	for i := 0; i < len(states); i += batchSize {
		end := min(i+batchSize, len(states))

		batch := states[i:end]
		rows := []map[string]string{}

		for _, state := range batch {
			data := state.Data()
			rows = append(rows, data)
		}

		sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
			Insert(store.stateTableName).
			Prepared(true).
			Rows(rows).
			ToSQL()

		if errSql != nil {
			return errSql
		}

		if store.debugEnabled {
			log.Println(sqlStr)
		}

		_, err := store.db.Exec(sqlStr, params...)

		if err != nil {
			return err
		}

		for _, state := range batch {
			state.MarkAsNotDirty()
		}
	}

	return nil
}

func (store *storeImplementation) StateList(ctx context.Context, options StateQueryOptions) ([]State, error) {
	q := store.stateQuery(options)

	sqlStr, params, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []State{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(database.NewQueryableContext(ctx, store.db), sqlStr, params...)
	if err != nil {
		return []State{}, err
	}

	list := []State{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewStateFromExistingData(modelMap)
		list = append(list, *model)
	})

	return list, nil
}

func (store *storeImplementation) stateQuery(options StateQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.stateTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.Status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	}

	if options.CountryCode != "" {
		q = q.Where(goqu.C(COLUMN_COUNTRY_CODE).Eq(options.CountryCode))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
	}

	if options.CountOnly {
		q = q.Select(goqu.COUNT("*"))
	}

	if options.Limit > 0 {
		q = q.Limit(uint(options.Limit))
	}

	if options.Offset > 0 {
		q = q.Offset(uint(options.Offset))
	}

	if options.OrderBy != "" {
		if options.SortOrder == "desc" {
			q = q.Order(goqu.C(options.OrderBy).Desc())
		} else {
			q = q.Order(goqu.C(options.OrderBy).Asc())
		}
	}

	if !options.WithDeleted {
		q = q.Where(goqu.C(COLUMN_DELETED_AT).Eq(sb.NULL_DATETIME))
	}

	return q
}

func (store *storeImplementation) TimezoneCreate(ctx context.Context, timezone *Timezone) error {
	timezone.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	timezone.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := timezone.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.timezoneTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	if err != nil {
		return err
	}

	timezone.MarkAsNotDirty()

	return nil
}

func (store *storeImplementation) TimezoneList(ctx context.Context, options TimezoneQueryOptions) ([]Timezone, error) {
	q := store.timezoneQuery(options)

	sqlStr, params, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []Timezone{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(database.NewQueryableContext(ctx, store.db), sqlStr, params...)
	if err != nil {
		return []Timezone{}, err
	}

	list := []Timezone{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewTimezoneFromExistingData(modelMap)
		list = append(list, *model)
	})

	return list, nil
}

func (store *storeImplementation) timezoneQuery(options TimezoneQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.timezoneTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.Status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	}

	if options.CountryCode != "" {
		q = q.Where(goqu.C(COLUMN_COUNTRY_CODE).Eq(options.CountryCode))
	}

	if options.Timezone != "" {
		q = q.Where(goqu.C(COLUMN_TIMEZONE).Eq(options.Timezone))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
	}

	if options.CountOnly {
		q = q.Select(goqu.COUNT("*"))
	}

	if options.Limit > 0 {
		q = q.Limit(uint(options.Limit))
	}

	if options.Offset > 0 {
		q = q.Offset(uint(options.Offset))
	}

	if options.OrderBy != "" {
		if options.SortOrder == "desc" {
			q = q.Order(goqu.C(options.OrderBy).Desc())
		} else {
			q = q.Order(goqu.C(options.OrderBy).Asc())
		}
	}

	if !options.WithDeleted {
		q = q.Where(goqu.C(COLUMN_DELETED_AT).Eq(sb.NULL_DATETIME))
	}

	return q
}
