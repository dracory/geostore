package geostore

import (
	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// == CLASS =================================================================

type State struct {
	dataobject.DataObject
}

// == CONSTRUCTORS ==========================================================

func NewState() *State {
	state := &State{}
	state.SetID(uid.UuidV4())
	state.SetStatus(COUNTRY_STATUS_ACTIVE)
	state.SetName("")
	state.SetCountryCode("")
	state.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	state.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	state.SetDeletedAt(sb.NULL_DATETIME)
	return state
}

func NewStateFromExistingData(data map[string]string) *State {
	state := &State{}
	state.Hydrate(data)
	return state
}

// == SETTERS AND GETTERS ===================================================

func (o *State) CountryCode() string {
	return o.Get(COLUMN_COUNTRY_CODE)
}

func (o *State) SetCountryCode(countryCodeIso2 string) *State {
	o.Set(COLUMN_COUNTRY_CODE, countryCodeIso2)
	return o
}

func (o *State) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *State) SetCreatedAt(createdAt string) *State {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *State) DeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *State) SetDeletedAt(deletedAt string) *State {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *State) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *State) SetName(name string) *State {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *State) StateCode() string {
	return o.Get(COLUMN_STATE_CODE)
}

func (o *State) SetStateCode(stateCode string) *State {
	o.Set(COLUMN_STATE_CODE, stateCode)
	return o
}

func (o *State) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *State) SetStatus(status string) *State {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *State) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *State) SetUpdatedAt(updatedAt string) *State {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
