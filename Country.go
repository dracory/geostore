package geostore

import (
	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// == CLASS =================================================================

type Country struct {
	dataobject.DataObject
}

// == CONSTRUCTORS ==========================================================

func NewCountry() *Country {
	country := &Country{}
	country.SetID(uid.HumanUid())
	country.SetStatus(COUNTRY_STATUS_ACTIVE)
	country.SetName("")
	country.SetIsoCode2("")
	country.SetIsoCode3("")
	country.SetContinent("")
	country.SetPhonePrefix("")
	country.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	country.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	country.SetDeletedAt(sb.NULL_DATETIME)
	return country
}

func NewCountryFromExistingData(data map[string]string) *Country {
	country := &Country{}
	country.Hydrate(data)
	return country
}

// == SETTERS AND GETTERS ===================================================

func (o *Country) Continent() string {
	return o.Get(COLUMN_CONTINENT)
}

func (o *Country) SetContinent(continent string) *Country {
	o.Set(COLUMN_CONTINENT, continent)
	return o
}

func (o *Country) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *Country) SetCreatedAt(createdAt string) *Country {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *Country) DeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *Country) SetDeletedAt(deletedAt string) *Country {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *Country) IsoCode2() string {
	return o.Get(COLUMN_ISO2_CODE)
}

func (o *Country) SetIsoCode2(isoCode2 string) *Country {
	o.Set(COLUMN_ISO2_CODE, isoCode2)
	return o
}

func (o *Country) IsoCode3() string {
	return o.Get(COLUMN_ISO3_CODE)
}

func (o *Country) SetIsoCode3(isoCode3 string) *Country {
	o.Set(COLUMN_ISO3_CODE, isoCode3)
	return o
}

func (o *Country) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *Country) SetName(name string) *Country {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *Country) PhonePrefix() string {
	return o.Get(COLUMN_PHONE_PREFIX)
}

func (o *Country) SetPhonePrefix(phonePrefix string) *Country {
	o.Set(COLUMN_PHONE_PREFIX, phonePrefix)
	return o
}

func (o *Country) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *Country) SetStatus(status string) *Country {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *Country) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *Country) SetUpdatedAt(updatedAt string) *Country {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
