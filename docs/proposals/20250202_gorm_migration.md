# Migration to GORM

## Status: Proposed

## Overview

This proposal suggests migrating GeoStore's internal database storage layer to use GORM (Go Object Relational Mapper) as an implementation detail. The existing `StoreInterface` and public APIs will remain unchanged - GORM will only replace the goqu-based implementation behind the scenes.

## Constraints

- **Interface Preservation**: The `StoreInterface` and all public methods remain unchanged
- **Implementation Only**: GORM is strictly an internal implementation detail for database access
- **API Compatibility**: No breaking changes to existing code using GeoStore

## Current Implementation

Currently, GeoStore uses a goqu-based database abstraction layer:

- **Goqu Query Builder**: SQL queries constructed using goqu library
- **Manual Mapping**: Hand-written code for scanning rows into structs using `database.SelectToMapString`
- **Custom Migrations**: SQL-based migration scripts in `sql_create_table.go`
- **Dialect Handling**: Goqu handles SQLite and PostgreSQL differences
- **Entity Models**: Country, State, and Timezone entities with custom data objects

## Proposed Changes

1. **Adopt GORM v2**: Migrate from goqu to GORM as the primary ORM layer
2. **Model Definitions**: Define Country, State, and Timezone structs with GORM tags for automatic table mapping
3. **Auto-Migrations**: Replace custom SQL migrations with GORM's AutoMigrate for schema management
4. **Query Builder**: Replace goqu queries with GORM's fluent query API
5. **Relationship Handling**: Leverage GORM's associations for Country-State-Timezone relationships

## Implementation Details

The migration will replace goqu query building with GORM operations **inside the existing implementation**, without changing any interfaces:

```go
// gormCountry is the GORM model (internal, with tags)
type gormCountry struct {
    ID          string `gorm:"primaryKey"`
    Name        string `gorm:"not null"`
    Iso2Code    string `gorm:"uniqueIndex"`
    Iso3Code    string `gorm:"uniqueIndex"`
    Status      string `gorm:"default:active"`
    CreatedAt   string
    UpdatedAt   string
    DeletedAt   *string `gorm:"index"`
}

// Country struct remains unchanged (public interface)
type Country struct {
    // existing fields and methods remain the same
}

// NewCountryFromGorm constructor converts GORM model to Country
func NewCountryFromGorm(gc *gormCountry) *Country {
    return &Country{
        // map fields from gormCountry to Country
    }
}

// Internal method implementation changes only
func (s *storeImplementation) CountryFindByIso2(ctx context.Context, iso2Code string) (*Country, error) {
    // BEFORE: goqu query
    // sqlStr, params, errSql := goqu.Dialect(s.dbDriverName).
    //     From(s.countryTableName).
    //     Where(goqu.C(COLUMN_ISO2_CODE).Eq(iso2Code)).
    //     ToSQL()
    
    // AFTER: GORM query (same return type, same interface)
    var gc gormCountry
    result := s.gormDB.WithContext(ctx).Where("iso2_code = ?", iso2Code).First(&gc)
    if result.Error != nil {
        return nil, result.Error
    }
    return NewCountryFromGorm(&gc), nil
}
```

Key points:
- `gormCountry`, `gormState`, `gormTimezone` are internal GORM models with struct tags
- `Country`, `State`, `Timezone` structs remain unchanged as public entities
- `NewCountryFromGorm`, `NewStateFromGorm`, `NewTimezoneFromGorm` constructors bridge the types
- All public methods (`CountryCreate`, `CountryFindByIso2`, etc.) keep identical signatures
- `StoreInterface` methods remain unchanged; only their internal implementation uses GORM
- Migration is transparent to library consumers

## Pros (Benefits)

| Benefit | Description |
|---------|-------------|
| **Reduced Boilerplate** | Eliminates goqu query building and manual row mapping |
| **Type Safety** | Compile-time checks for database operations |
| **Cross-Database Support** | Built-in SQLite, PostgreSQL, MySQL, SQL Server, Oracle, CockroachDB, TiDB support |
| **Migration Management** | Automated schema migrations with AutoMigrate |
| **Query Building** | Fluent, chainable API for complex queries |
| **Relationships** | Easy handling of Country-State-Timezone associations |
| **Hooks** | Before/after create/update/delete callbacks |
| **Community & Testing** | 80K+ stars, extensively battle-tested in production systems |
| **Enterprise Adoption** | SQL Server/Oracle support opens corporate markets |
| **Testing** | Easier to mock database layer for unit tests |
| **Connection Pooling** | Built-in connection management |
| **Future-Proofing** | Large ecosystem, active maintenance, rapid bug fixes |

## Cons (Drawbacks)

| Drawback | Description |
|----------|-------------|
| **Performance Overhead** | ORM abstraction adds latency compared to goqu |
| **Learning Curve** | Team needs to learn GORM-specific patterns and conventions |
| **Migration Risk** | Existing Country/State/Timezone data must be preserved during transition |
| **Complex Queries** | Custom goqu queries may still be needed for complex operations |
| **Dependency** | Adds external dependency (GORM + database drivers) |
| **Generated SQL** | Less control over exact SQL generated |
| **Memory Usage** | Reflection-heavy operations may use more memory |
| **Magic Behavior** | Implicit behaviors (callbacks, automatic timestamps) can be confusing |

## Risks and Mitigations

- **Data Loss Risk**: Migration could corrupt existing data. Mitigation: Full backup and dry-run migrations.
- **Performance Regression**: ORM overhead may slow operations. Mitigation: Benchmark before/after, optimize critical paths.
- **Breaking Changes**: API surface may change. Mitigation: Maintain backward compatibility layer.
- **Learning Curve**: Team unfamiliarity. Mitigation: Documentation and training sessions.
- **Long-term Maintenance Risk**: Staying with goqu poses higher risk due to limited community testing and fewer edge cases discovered.

## Effort Estimation

- **Research & Planning**: 3-5 days
- **Core Migration**: 2-3 weeks
- **Testing & Validation**: 1-2 weeks
- **Documentation**: 3-5 days
- **Total**: ~4-6 weeks

## Conclusion

Migrating to GORM would modernize GeoStore's database layer, reducing maintenance burden and improving developer experience. The migration from goqu to GORM would provide better relationship handling, automated migrations, and extensive database support while maintaining the existing public API.

**Key considerations:**
- **Community & Testing**: GORM's 80K+ stars and extensive production testing provide significantly more reliability than goqu's smaller community
- **Database Coverage**: Support for SQL Server, Oracle, CockroachDB, TiDB opens enterprise markets and future-proofs the library
- **Risk Balance**: While migration has short-term risks, staying with goqu carries higher long-term maintenance risks due to limited battle-testing

The team should weigh the short-term migration effort against the long-term benefits of broader database support, extensive community testing, and reduced maintenance burden. For a geographical data store where reliability and broad compatibility are critical, GORM's advantages likely outweigh the migration costs.
