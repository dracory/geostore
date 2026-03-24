# Migration to GORM

## Status

**Proposed** | Date: 2025-02-02 | Author: Development Team

## Executive Summary

This proposal recommends migrating GeoStore's internal database storage layer from goqu to GORM (Go Object Relational Mapper) v2. This is a **non-breaking, implementation-only change** that maintains full API compatibility while modernizing the data access layer. The migration reduces boilerplate code, improves maintainability, and expands database support to enterprise systems (SQL Server, Oracle, CockroachDB).

**Expected Impact:**
- 40-50% reduction in database-related code
- Support for 7+ database engines (currently 2)
- Improved developer experience with type-safe queries
- Zero breaking changes to public API

## Overview

GeoStore currently uses goqu for query building and manual row mapping. This proposal suggests replacing the internal implementation with GORM while preserving the existing `StoreInterface` contract. GORM will be strictly an implementation detail - library consumers will experience no API changes.

## Constraints

- **Interface Preservation**: The `StoreInterface` and all public methods remain unchanged
- **Implementation Only**: GORM is strictly an internal implementation detail for database access
- **API Compatibility**: No breaking changes to existing code using GeoStore

## Current State Analysis

### Existing Architecture

GeoStore uses a goqu-based database abstraction layer with the following components:

- **Query Builder**: goqu library for SQL query construction
- **Manual Mapping**: Custom row scanning using `database.SelectToMapString`
- **Schema Management**: SQL-based migrations in `sql_create_table.go`
- **Dialect Support**: SQLite and PostgreSQL via goqu dialects
- **Entity Models**: Country, State, and Timezone structs
- **Data Layer**: Direct SQL execution with manual error handling

### Pain Points

1. **Boilerplate Code**: Repetitive query building and row mapping logic
2. **Limited Database Support**: Only SQLite and PostgreSQL officially supported
3. **Manual Schema Management**: SQL migrations require manual synchronization
4. **Type Safety Gaps**: Runtime errors for column mismatches
5. **Testing Complexity**: Difficult to mock database layer
6. **Maintenance Burden**: Small community means fewer bug fixes and edge case coverage

## Proposed Solution

### Architecture Changes

1. **ORM Layer**: Replace goqu with GORM v2 as the internal database abstraction
2. **Dual Model Pattern**: Introduce internal GORM models alongside existing public entities
3. **Schema Management**: Migrate to GORM's AutoMigrate for table creation and updates
4. **Query API**: Replace goqu query builders with GORM's chainable query interface
5. **Relationship Mapping**: Use GORM associations for entity relationships
6. **Transaction Support**: Leverage GORM's built-in transaction management

### Design Principles

- **Zero Breaking Changes**: Public API remains 100% compatible
- **Internal Only**: GORM types never exposed in public interfaces
- **Gradual Migration**: Implement incrementally, method by method
- **Backward Compatibility**: Support existing database schemas during transition
- **Performance Parity**: Maintain or improve current query performance

## Technical Implementation

### Dual Model Architecture

The migration introduces internal GORM models that coexist with public entity structs:

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

### Key Design Decisions

- **Internal GORM Models**: `gormCountry`, `gormState`, `gormTimezone` with struct tags (never exported)
- **Public Entities Unchanged**: `Country`, `State`, `Timezone` maintain existing structure
- **Conversion Layer**: `NewCountryFromGorm()` family of constructors bridge internal/external types
- **Method Signatures Preserved**: All `StoreInterface` methods keep identical signatures
- **Transparent Migration**: Library consumers see no changes
- **Table Name Mapping**: GORM models use existing table names via `TableName()` method

### Migration Strategy

**Phase 1: Foundation (Week 1)**
- Add GORM dependency to `go.mod`
- Create internal GORM model structs with appropriate tags
- Implement conversion functions between GORM models and public entities
- Add GORM DB instance to `storeImplementation` struct
- Ensure existing goqu code continues to work

**Phase 2: Core Operations (Week 2-3)**
- Migrate Create operations (Country, State, Timezone)
- Migrate Find operations (ByID, ByIso2, ByIso3, etc.)
- Migrate List operations with pagination and filtering
- Migrate Update operations
- Migrate Delete operations (soft delete support)
- Run parallel testing: goqu vs GORM results

**Phase 3: Advanced Features (Week 4)**
- Migrate seeding functions (`seedCountriesIfTableIsEmpty`, etc.)
- Replace custom SQL migrations with GORM AutoMigrate
- Implement query options using GORM scopes
- Add transaction support for batch operations
- Performance optimization and query tuning

**Phase 4: Cleanup & Documentation (Week 5-6)**
- Remove goqu dependencies
- Update tests to use GORM patterns
- Benchmark and optimize critical paths
- Update internal documentation
- Final validation and release preparation

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

## Risk Assessment & Mitigation

| Risk | Severity | Probability | Mitigation Strategy |
|------|----------|-------------|---------------------|
| **Data Loss** | Critical | Low | • Full database backups before migration<br>• Dry-run migrations in staging<br>• Rollback plan with goqu fallback<br>• Schema validation tests |
| **Performance Regression** | High | Medium | • Benchmark suite (before/after)<br>• Query analysis and optimization<br>• Prepared statement caching<br>• Index optimization |
| **API Breaking Changes** | Critical | Very Low | • Comprehensive integration tests<br>• Interface contract validation<br>• Semantic versioning enforcement<br>• Consumer compatibility matrix |
| **Migration Bugs** | High | Medium | • Incremental migration (method-by-method)<br>• Parallel execution (goqu + GORM)<br>• Extensive test coverage (>90%)<br>• Canary deployments |
| **Team Learning Curve** | Medium | High | • GORM training sessions<br>• Code review guidelines<br>• Internal documentation<br>• Pair programming |
| **Dependency Risk** | Medium | Low | • Pin GORM version<br>• Monitor security advisories<br>• Vendor dependencies if needed |
| **Rollback Complexity** | Medium | Low | • Feature flag for GORM/goqu toggle<br>• Maintain goqu code during transition<br>• Database schema compatibility |

### Rollback Plan

If critical issues arise during migration:

1. **Immediate Rollback** (< 1 hour)
   - Revert to previous release
   - Restore database from backup if schema changed
   - Enable goqu code path via feature flag

2. **Partial Rollback** (< 4 hours)
   - Disable GORM for specific operations
   - Route traffic to goqu implementation
   - Investigate and fix issues in isolated environment

3. **Schema Rollback** (< 24 hours)
   - Run reverse migrations if schema changed
   - Validate data integrity
   - Resume operations with goqu

## Resource Planning

### Effort Estimation

| Phase | Duration | Team Size | Effort (Person-Days) |
|-------|----------|-----------|----------------------|
| Research & Planning | 3-5 days | 2 developers | 6-10 days |
| Foundation Setup | 1 week | 1 developer | 5 days |
| Core Migration | 2-3 weeks | 2 developers | 20-30 days |
| Testing & Validation | 1-2 weeks | 2 developers + 1 QA | 15-25 days |
| Documentation | 3-5 days | 1 developer | 3-5 days |
| **Total** | **5-7 weeks** | **2-3 people** | **49-75 days** |

### Prerequisites

- GORM v2 documentation review
- Test environment with SQLite and PostgreSQL
- Backup and restore procedures validated
- Benchmark baseline established
- Stakeholder approval

### Success Metrics

1. **Functional Correctness**
   - 100% of existing tests pass
   - Zero regression bugs in production
   - All query results match goqu implementation

2. **Performance**
   - Query latency within 10% of baseline
   - Memory usage within 15% of baseline
   - No degradation in throughput (queries/sec)

3. **Code Quality**
   - 30-50% reduction in database-related LOC
   - Test coverage maintained at >85%
   - Zero new linter warnings

4. **Compatibility**
   - Successful testing on SQLite, PostgreSQL
   - Optional: Validation on MySQL, SQL Server
   - Backward compatibility with existing schemas

## Alternatives Considered

### Option 1: Stay with goqu (Status Quo)
**Pros:** No migration effort, zero risk, team familiarity  
**Cons:** Limited database support, higher maintenance burden, smaller community  
**Verdict:** ❌ Rejected - Long-term technical debt outweighs short-term stability

### Option 2: Raw SQL with sqlx
**Pros:** Maximum performance, full control, minimal abstraction  
**Cons:** More boilerplate than GORM, manual migrations, no relationship handling  
**Verdict:** ❌ Rejected - Increases code complexity without sufficient benefits

### Option 3: Migrate to GORM (Recommended)
**Pros:** Reduced boilerplate, broad database support, active community, type safety  
**Cons:** Migration effort, learning curve, slight performance overhead  
**Verdict:** ✅ **Recommended** - Best balance of maintainability and functionality

### Option 4: Hybrid Approach (GORM + Raw SQL)
**Pros:** GORM for CRUD, raw SQL for complex queries  
**Cons:** Increased complexity, two paradigms to maintain  
**Verdict:** ⚠️ Consider as fallback if GORM performance is insufficient

## Decision Criteria

| Criterion | Weight | goqu | GORM | sqlx |
|-----------|--------|------|------|------|
| Code Maintainability | 25% | 6/10 | 9/10 | 5/10 |
| Database Support | 20% | 5/10 | 10/10 | 8/10 |
| Performance | 20% | 8/10 | 7/10 | 9/10 |
| Community & Support | 15% | 4/10 | 10/10 | 7/10 |
| Migration Effort | 10% | 10/10 | 5/10 | 4/10 |
| Type Safety | 10% | 6/10 | 9/10 | 7/10 |
| **Weighted Score** | | **6.4** | **8.4** | **6.9** |

## Recommendation

**Proceed with GORM migration** with the following conditions:

1. ✅ **Approve** if team commits to 5-7 week timeline
2. ✅ **Approve** if performance benchmarks show <10% regression
3. ✅ **Approve** if rollback plan is validated in staging
4. ⚠️ **Defer** if critical production issues require immediate attention
5. ❌ **Reject** if team lacks bandwidth for proper testing

## Conclusion

Migrating to GORM modernizes GeoStore's database layer while maintaining full API compatibility. The migration reduces technical debt, expands database support to enterprise systems, and improves long-term maintainability.

### Strategic Value

- **Community Strength**: GORM's 80K+ GitHub stars and extensive production usage provide superior reliability and faster bug fixes compared to goqu's smaller ecosystem
- **Market Expansion**: SQL Server and Oracle support enables enterprise adoption and government contracts
- **Future-Proofing**: Active development and broad database support ensure long-term viability
- **Developer Experience**: Reduced boilerplate and type-safe queries improve productivity

### Risk-Benefit Analysis

While the migration carries short-term implementation risk, maintaining goqu poses higher long-term risks:
- Limited community testing means undiscovered edge cases
- Narrow database support restricts market opportunities  
- Higher maintenance burden as team expertise declines

**For a geographical data library where reliability, compatibility, and maintainability are paramount, GORM's advantages justify the migration investment.**

---

## Next Steps

1. **Stakeholder Review** (Week 0): Present proposal to team, gather feedback
2. **Approval Decision** (Week 0): Go/No-Go decision from technical leadership
3. **Kickoff** (Week 1): If approved, begin Phase 1 implementation
4. **Status Updates**: Weekly progress reports during migration
5. **Post-Migration Review**: Retrospective after completion to capture learnings
