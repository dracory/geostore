package geostore

import (
	"time"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// == CLASS =================================================================

type Country struct {
	orm.ShortID

	StatusField      string `db:"status"`
	Iso2CodeField    string `db:"iso2_code"`
	Iso3CodeField    string `db:"iso3_code"`
	NameField        string `db:"name"`
	ContinentField   string `db:"continent"`
	PhonePrefixField string `db:"phone_prefix"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate

	originalData map[string]string
}

// == CONSTRUCTORS ==========================================================

func NewCountry() *Country {
	country := &Country{}
	country.SetID(neatuid.GenerateShortID())
	country.SetStatus(COUNTRY_STATUS_ACTIVE)
	country.SetName("")
	country.SetIsoCode2("")
	country.SetIsoCode3("")
	country.SetContinent("")
	country.SetPhonePrefix("")
	country.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	country.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	country.SetSoftDeletedAt(MAX_DATETIME)
	country.MarkAsNotDirty()
	return country
}

func NewCountryFromExistingData(data map[string]string) *Country {
	country := &Country{}
	country.SetID(data[COLUMN_ID])
	country.SetStatus(data[COLUMN_STATUS])
	country.SetIsoCode2(data[COLUMN_ISO2_CODE])
	country.SetIsoCode3(data[COLUMN_ISO3_CODE])
	country.SetName(data[COLUMN_NAME])
	country.SetContinent(data[COLUMN_CONTINENT])
	country.SetPhonePrefix(data[COLUMN_PHONE_PREFIX])
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		country.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_UPDATED_AT]; ok {
		country.SetUpdatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		country.SetSoftDeletedAt(v)
	}
	country.MarkAsNotDirty()
	return country
}

// == SETTERS AND GETTERS ===================================================

func (o *Country) Continent() string {
	return o.ContinentField
}

func (o *Country) SetContinent(continent string) *Country {
	o.ContinentField = continent
	return o
}

func (o *Country) CreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

func (o *Country) CreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

func (o *Country) SetCreatedAt(createdAt string) *Country {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

func (o *Country) GetSoftDeletedAt() string {
	if o.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletedAt).ToDateTimeString()
}

func (o *Country) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletedAt)
}

func (o *Country) SetSoftDeletedAt(softDeletedAt string) *Country {
	if softDeletedAt == "" {
		o.SoftDeletedAt = time.Time{}
		return o
	}
	o.SoftDeletedAt = carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
	return o
}

func (o *Country) IsoCode2() string {
	return o.Iso2CodeField
}

func (o *Country) SetIsoCode2(isoCode2 string) *Country {
	o.Iso2CodeField = isoCode2
	return o
}

func (o *Country) IsoCode3() string {
	return o.Iso3CodeField
}

func (o *Country) SetIsoCode3(isoCode3 string) *Country {
	o.Iso3CodeField = isoCode3
	return o
}

func (o *Country) Name() string {
	return o.NameField
}

func (o *Country) SetName(name string) *Country {
	o.NameField = name
	return o
}

func (o *Country) PhonePrefix() string {
	return o.PhonePrefixField
}

func (o *Country) SetPhonePrefix(phonePrefix string) *Country {
	o.PhonePrefixField = phonePrefix
	return o
}

func (o *Country) Status() string {
	return o.StatusField
}

func (o *Country) SetStatus(status string) *Country {
	o.StatusField = status
	return o
}

func (o *Country) UpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

func (o *Country) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

func (o *Country) SetUpdatedAt(updatedAt string) *Country {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

func (o *Country) ID() string {
	return o.ShortID.ID
}

func (o *Country) SetID(id string) *Country {
	o.ShortID.ID = id
	return o
}

// == DATA METHODS ==========================================================

func (o *Country) Data() map[string]string {
	return map[string]string{
		COLUMN_ID:              o.ID(),
		COLUMN_STATUS:          o.Status(),
		COLUMN_ISO2_CODE:       o.IsoCode2(),
		COLUMN_ISO3_CODE:       o.IsoCode3(),
		COLUMN_NAME:            o.Name(),
		COLUMN_CONTINENT:       o.Continent(),
		COLUMN_PHONE_PREFIX:    o.PhonePrefix(),
		COLUMN_CREATED_AT:      o.CreatedAt(),
		COLUMN_UPDATED_AT:      o.UpdatedAt(),
		COLUMN_SOFT_DELETED_AT: o.GetSoftDeletedAt(),
	}
}

func (o *Country) MarkAsNotDirty(columns ...string) {
	o.originalData = o.Data()
}
