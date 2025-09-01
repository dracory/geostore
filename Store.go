package geostore

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

// const DISCOUNT_TABLE_NAME = "shop_country"

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

type Store struct {
	countryTableName   string
	stateTableName     string
	timezoneTableName  string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debugEnabled       bool
}

// AutoMigrate auto migrate
func (store *Store) AutoMigrate() error {
	// create country table
	sql := store.sqlCountryTableCreate()

	if sql == "" {
		return errors.New("country table create sql is empty")
	}

	_, err := store.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	// create state table
	sql = store.sqlStateTableCreate()

	if sql == "" {
		return errors.New("state table create sql is empty")
	}

	_, err = store.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	// create timezone table
	sql = store.sqlTimezoneTableCreate()

	if sql == "" {
		return errors.New("timezone table create sql is empty")
	}

	_, err = store.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	// seed country table
	err = store.seedCountriesIfTableIsEmpty()

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

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

func (store *Store) CountryCreate(country *Country) error {
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

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	country.MarkAsNotDirty()

	return nil
}

func (store *Store) CountryDelete(country *Country) error {
	if country == nil {
		return errors.New("country is nil")
	}

	return store.CountryDeleteByID(country.ID())
}

func (store *Store) CountryDeleteByID(id string) error {
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

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

func (store *Store) CountryFindByID(id string) (*Country, error) {
	if id == "" {
		return nil, errors.New("country id is empty")
	}

	list, err := store.CountryList(CountryQueryOptions{
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

func (store *Store) CountryFindByIso2(iso2Code string) (*Country, error) {
	if iso2Code == "" {
		return nil, errors.New("country iso2 code is empty")
	}

	list, err := store.CountryList(CountryQueryOptions{
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

func (store *Store) CountryNameFindByIso2(iso2Code string) (string, error) {
	country, err := store.CountryFindByIso2(iso2Code)

	if err != nil {
		return "", err
	}

	if country == nil {
		return "", nil
	}

	return country.Name(), nil
}

func (store *Store) CountryList(options CountryQueryOptions) ([]Country, error) {
	q := store.countryQuery(options)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []Country{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr)
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

func (store *Store) CountrySoftDelete(country *Country) error {
	if country == nil {
		return errors.New("country is nil")
	}

	country.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.CountryUpdate(country)
}

func (store *Store) CountrySoftDeleteByID(id string) error {
	country, err := store.CountryFindByID(id)

	if err != nil {
		return err
	}

	return store.CountrySoftDelete(country)
}

func (store *Store) CountryUpdate(country *Country) error {
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

	_, err := store.db.Exec(sqlStr, params...)

	country.MarkAsNotDirty()

	return err
}

func (store *Store) countryQuery(options CountryQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.countryTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.Status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
	}

	if options.Iso2 != "" {
		q = q.Where(goqu.C(COLUMN_ISO2_CODE).Eq(options.Iso2))
	}

	if options.Iso3 != "" {
		q = q.Where(goqu.C(COLUMN_ISO3_CODE).Eq(options.Iso3))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	sortOrder := "desc"
	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.OrderBy != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy).Desc())
		}
	}

	if !options.WithDeleted {
		q = q.Where(goqu.C(COLUMN_DELETED_AT).Eq(sb.NULL_DATETIME))
	}

	return q
}

func (store *Store) StateCreate(state *State) error {
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

func (store *Store) StatesCreate(states []*State) error {
	for index, state := range states {
		state.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
		state.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
		states[index] = state
	}

	rows := []map[string]string{}

	for _, state := range states {
		data := state.Data()
		rows = append(rows, data)
	}

	limit := 500
	batches := [][]map[string]string{}

	for i := 0; i < len(rows); i += limit {
		batch := rows[i:min(i+limit, len(rows))]
		batches = append(batches, batch)
	}

	for _, batch := range batches {
		sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
			Insert(store.stateTableName).
			Prepared(true).
			Rows(batch).
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
	}

	for _, state := range states {
		state.MarkAsNotDirty()
	}

	return nil
}

func (store *Store) StateList(options StateQueryOptions) ([]State, error) {
	q := store.stateQuery(options)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []State{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr)
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

func (store *Store) stateQuery(options StateQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.stateTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.Status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
	}

	if options.CountryCode != "" {
		q = q.Where(goqu.C(COLUMN_COUNTRY_CODE).Eq(options.CountryCode))
	}

	if options.StateCode != "" {
		q = q.Where(goqu.C(COLUMN_STATE_CODE).Eq(options.StateCode))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	sortOrder := "desc"
	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.OrderBy != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy).Desc())
		}
	}

	if !options.WithDeleted {
		q = q.Where(goqu.C(COLUMN_DELETED_AT).Eq(sb.NULL_DATETIME))
	}

	return q
}

func (store *Store) TimezoneCreate(timezone *Timezone) error {
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

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	timezone.MarkAsNotDirty()

	return nil
}

func (store *Store) TimezoneList(options TimezoneQueryOptions) ([]Timezone, error) {
	q := store.timezoneQuery(options)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []Timezone{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr)
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

func (store *Store) timezoneQuery(options TimezoneQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.timezoneTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.Status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
	}

	if options.CountryCode != "" {
		q = q.Where(goqu.C(COLUMN_COUNTRY_CODE).Eq(options.CountryCode))
	}

	if options.Timezone != "" {
		q = q.Where(goqu.C(COLUMN_TIMEZONE).Eq(options.Timezone))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	sortOrder := "desc"
	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.OrderBy != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy).Desc())
		}
	}

	if !options.WithDeleted {
		q = q.Where(goqu.C(COLUMN_DELETED_AT).Eq(sb.NULL_DATETIME))
	}

	return q
}
