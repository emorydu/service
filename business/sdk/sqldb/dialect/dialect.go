// Package dialect provides cross-cutting helpers that vary between SQL
// engines but are otherwise reusable by every domain store. A store wires
// in one implementation (Postgres, MySQL, SQLite, ...) and delegates the
// small set of engine-specific decisions to it.
//
// The set of decisions encapsulated here is intentionally narrow. Only add
// to the interface when a real, second engine forces a difference.
//
// This package is also the only one below business/sdk/sqldb that names a
// concrete driver's error type, so the engine-specific SQL fragments and the
// engine-specific error codes are read side by side rather than scattered
// across the query helpers.
package dialect

import "bytes"

// The dialects are stateless, so callers share one value per engine rather
// than constructing their own.
var (
	Postgres Dialect = PostgresDialect{}
	MySQL    Dialect = MySQLDialect{}
	SQLite   Dialect = SQLiteDialect{}

	// Generic is the Dialect for a driver this package does not recognize.
	// It claims no engine-specific error and paginates with the form the
	// widest range of engines accepts, so an unknown driver keeps the
	// behavior it had before error translation existed instead of failing.
	Generic Dialect = GenericDialect{}
)

// Dialect describes the engine-specific behavior a store needs in order
// to compose portable SQL from shared fragments and to interpret the errors
// an engine reports for the two constraint failures a store maps onto a
// domain error.
type Dialect interface {

	// Name reports a short identifier for the engine (for logging only).
	Name() string

	// Paginate appends a pagination clause to buf. The clause must consume
	// the named bind variables ":offset" and ":rows_per_page" already
	// supplied by the caller in the parameter map.
	Paginate(buf *bytes.Buffer)

	// IsUndefinedTable reports whether err is the engine's error for a query
	// against a table that does not exist.
	IsUndefinedTable(err error) bool

	// IsUniqueViolation reports whether err is the engine's error for a
	// violation of a unique constraint.
	IsUniqueViolation(err error) bool
}

// For returns the Dialect for the driver name database/sql recorded when the
// connection was opened, for example "pgx" or "mysql". An unrecognized name
// yields Generic rather than an error: For is consulted while an error is
// already on its way back to the caller, so it must not introduce a second
// failure path of its own.
func For(driverName string) Dialect {
	switch driverName {
	case "pgx", "pgx/v5", "postgres":
		return Postgres
	case "mysql":
		return MySQL
	case "sqlite3", "sqlite":
		return SQLite
	default:
		return Generic
	}
}
