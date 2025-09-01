package geostore

import (
	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// == CLASS =================================================================

type Timezone struct {
	dataobject.DataObject
}

// == CONSTRUCTORS ==========================================================

func NewTimezone() *Timezone {
	tz := &Timezone{}
	tz.SetID(uid.HumanUid())
	tz.SetStatus(TIMEZONE_STATUS_ACTIVE)
	tz.SetTimezone("")
	tz.SetZoneName("")
	tz.SetGlobalName("")
	tz.SetCountryCode("")
	tz.SetOffset("")
	tz.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	tz.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	tz.SetDeletedAt(sb.NULL_DATETIME)
	return tz
}

func NewTimezoneFromExistingData(data map[string]string) *Timezone {
	tz := &Timezone{}
	tz.Hydrate(data)
	return tz
}

// == SETTERS AND GETTERS ===================================================

func (o *Timezone) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *Timezone) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *Timezone) SetCreatedAt(createdAt string) *Timezone {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *Timezone) CountryCode() string {
	return o.Get(COLUMN_COUNTRY_CODE)
}

func (o *Timezone) SetCountryCode(countryCode string) *Timezone {
	o.Set(COLUMN_COUNTRY_CODE, countryCode)
	return o
}

func (o *Timezone) DeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *Timezone) SetDeletedAt(deletedAt string) *Timezone {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *Timezone) GlobalName() string {
	return o.Get(COLUMN_GLOBAL_NAME)
}

func (o *Timezone) SetGlobalName(globalName string) *Timezone {
	o.Set(COLUMN_GLOBAL_NAME, globalName)
	return o
}

func (o *Timezone) Offset() string {
	return o.Get(COLUMN_OFFSET)
}

func (o *Timezone) SetOffset(offset string) *Timezone {
	o.Set(COLUMN_OFFSET, offset)
	return o
}

func (o *Timezone) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *Timezone) SetStatus(status string) *Timezone {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *Timezone) Timezone() string {
	return o.Get(COLUMN_TIMEZONE)
}

func (o *Timezone) SetTimezone(timezone string) *Timezone {
	o.Set(COLUMN_TIMEZONE, timezone)
	return o
}

func (o *Timezone) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *Timezone) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt(), carbon.UTC)
}

func (o *Timezone) SetUpdatedAt(updatedAt string) *Timezone {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *Timezone) ZoneName() string {
	return o.Get(COLUMN_ZONE_NAME)
}

func (o *Timezone) SetZoneName(zoneName string) *Timezone {
	o.Set(COLUMN_ZONE_NAME, zoneName)
	return o
}
