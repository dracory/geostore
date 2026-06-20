package geostore

import (
	"time"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// == CLASS =================================================================

type Timezone struct {
	orm.ShortID

	StatusField      string `db:"status"`
	TimezoneField    string `db:"timezone"`
	ZoneNameField    string `db:"zone_name"`
	GlobalNameField  string `db:"global_name"`
	CountryCodeField string `db:"country_code"`
	OffsetField      string `db:"offset"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate

	originalData map[string]string
}

// == CONSTRUCTORS ==========================================================

func NewTimezone() *Timezone {
	tz := &Timezone{}
	tz.SetID(neatuid.GenerateShortID())
	tz.SetStatus(TIMEZONE_STATUS_ACTIVE)
	tz.SetTimezone("")
	tz.SetZoneName("")
	tz.SetGlobalName("")
	tz.SetCountryCode("")
	tz.SetOffset("")
	tz.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	tz.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	tz.SetSoftDeletedAt(MAX_DATETIME)
	tz.MarkAsNotDirty()
	return tz
}

func NewTimezoneFromExistingData(data map[string]string) *Timezone {
	tz := &Timezone{}
	tz.SetID(data[COLUMN_ID])
	tz.SetStatus(data[COLUMN_STATUS])
	tz.SetTimezone(data[COLUMN_TIMEZONE])
	tz.SetZoneName(data[COLUMN_ZONE_NAME])
	tz.SetGlobalName(data[COLUMN_GLOBAL_NAME])
	tz.SetCountryCode(data[COLUMN_COUNTRY_CODE])
	tz.SetOffset(data[COLUMN_OFFSET])
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		tz.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_UPDATED_AT]; ok {
		tz.SetUpdatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		tz.SetSoftDeletedAt(v)
	}
	tz.MarkAsNotDirty()
	return tz
}

// == SETTERS AND GETTERS ===================================================

func (o *Timezone) CreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

func (o *Timezone) CreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

func (o *Timezone) SetCreatedAt(createdAt string) *Timezone {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

func (o *Timezone) CountryCode() string {
	return o.CountryCodeField
}

func (o *Timezone) SetCountryCode(countryCode string) *Timezone {
	o.CountryCodeField = countryCode
	return o
}

func (o *Timezone) GetSoftDeletedAt() string {
	if o.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletedAt).ToDateTimeString()
}

func (o *Timezone) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletedAt)
}

func (o *Timezone) SetSoftDeletedAt(softDeletedAt string) *Timezone {
	if softDeletedAt == "" {
		o.SoftDeletedAt = time.Time{}
		return o
	}
	o.SoftDeletedAt = carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
	return o
}

func (o *Timezone) GlobalName() string {
	return o.GlobalNameField
}

func (o *Timezone) SetGlobalName(globalName string) *Timezone {
	o.GlobalNameField = globalName
	return o
}

func (o *Timezone) Offset() string {
	return o.OffsetField
}

func (o *Timezone) SetOffset(offset string) *Timezone {
	o.OffsetField = offset
	return o
}

func (o *Timezone) Status() string {
	return o.StatusField
}

func (o *Timezone) SetStatus(status string) *Timezone {
	o.StatusField = status
	return o
}

func (o *Timezone) Timezone() string {
	return o.TimezoneField
}

func (o *Timezone) SetTimezone(timezone string) *Timezone {
	o.TimezoneField = timezone
	return o
}

func (o *Timezone) UpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

func (o *Timezone) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

func (o *Timezone) SetUpdatedAt(updatedAt string) *Timezone {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

func (o *Timezone) ZoneName() string {
	return o.ZoneNameField
}

func (o *Timezone) SetZoneName(zoneName string) *Timezone {
	o.ZoneNameField = zoneName
	return o
}

func (o *Timezone) ID() string {
	return o.ShortID.ID
}

func (o *Timezone) SetID(id string) *Timezone {
	o.ShortID.ID = id
	return o
}

// == DATA METHODS ==========================================================

func (o *Timezone) Data() map[string]string {
	return map[string]string{
		COLUMN_ID:              o.ID(),
		COLUMN_STATUS:          o.Status(),
		COLUMN_TIMEZONE:        o.Timezone(),
		COLUMN_ZONE_NAME:       o.ZoneName(),
		COLUMN_GLOBAL_NAME:     o.GlobalName(),
		COLUMN_COUNTRY_CODE:    o.CountryCode(),
		COLUMN_OFFSET:          o.Offset(),
		COLUMN_CREATED_AT:      o.CreatedAt(),
		COLUMN_UPDATED_AT:      o.UpdatedAt(),
		COLUMN_SOFT_DELETED_AT: o.GetSoftDeletedAt(),
	}
}

func (o *Timezone) MarkAsNotDirty(columns ...string) {
	o.originalData = o.Data()
}
