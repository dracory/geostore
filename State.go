package geostore

import (
	"time"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// == CLASS =================================================================

type State struct {
	orm.ShortID

	StatusField      string `db:"status"`
	CountryCodeField string `db:"country_code"`
	StateCodeField   string `db:"state_code"`
	NameField        string `db:"name"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate

	originalData map[string]string
}

// == CONSTRUCTORS ==========================================================

func NewState() *State {
	state := &State{}
	state.SetID(neatuid.GenerateShortID())
	state.SetStatus(STATE_STATUS_ACTIVE)
	state.SetName("")
	state.SetCountryCode("")
	state.SetStateCode("")
	state.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	state.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	state.SetSoftDeletedAt(MAX_DATETIME)
	state.MarkAsNotDirty()
	return state
}

func NewStateFromExistingData(data map[string]string) *State {
	state := &State{}
	state.SetID(data[COLUMN_ID])
	state.SetStatus(data[COLUMN_STATUS])
	state.SetCountryCode(data[COLUMN_COUNTRY_CODE])
	state.SetStateCode(data[COLUMN_STATE_CODE])
	state.SetName(data[COLUMN_NAME])
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		state.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_UPDATED_AT]; ok {
		state.SetUpdatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		state.SetSoftDeletedAt(v)
	}
	state.MarkAsNotDirty()
	return state
}

// == SETTERS AND GETTERS ===================================================

func (o *State) CountryCode() string {
	return o.CountryCodeField
}

func (o *State) SetCountryCode(countryCodeIso2 string) *State {
	o.CountryCodeField = countryCodeIso2
	return o
}

func (o *State) CreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

func (o *State) CreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

func (o *State) SetCreatedAt(createdAt string) *State {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

func (o *State) GetSoftDeletedAt() string {
	if o.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletedAt).ToDateTimeString()
}

func (o *State) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletedAt)
}

func (o *State) SetSoftDeletedAt(softDeletedAt string) *State {
	if softDeletedAt == "" {
		o.SoftDeletedAt = time.Time{}
		return o
	}
	o.SoftDeletedAt = carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
	return o
}

func (o *State) Name() string {
	return o.NameField
}

func (o *State) SetName(name string) *State {
	o.NameField = name
	return o
}

func (o *State) StateCode() string {
	return o.StateCodeField
}

func (o *State) SetStateCode(stateCode string) *State {
	o.StateCodeField = stateCode
	return o
}

func (o *State) Status() string {
	return o.StatusField
}

func (o *State) SetStatus(status string) *State {
	o.StatusField = status
	return o
}

func (o *State) UpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

func (o *State) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

func (o *State) SetUpdatedAt(updatedAt string) *State {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

func (o *State) ID() string {
	return o.ShortID.ID
}

func (o *State) SetID(id string) *State {
	o.ShortID.ID = id
	return o
}

// == DATA METHODS ==========================================================

func (o *State) Data() map[string]string {
	return map[string]string{
		COLUMN_ID:              o.ID(),
		COLUMN_STATUS:          o.Status(),
		COLUMN_COUNTRY_CODE:    o.CountryCode(),
		COLUMN_STATE_CODE:      o.StateCode(),
		COLUMN_NAME:            o.Name(),
		COLUMN_CREATED_AT:      o.CreatedAt(),
		COLUMN_UPDATED_AT:      o.UpdatedAt(),
		COLUMN_SOFT_DELETED_AT: o.GetSoftDeletedAt(),
	}
}

func (o *State) MarkAsNotDirty(columns ...string) {
	o.originalData = o.Data()
}
