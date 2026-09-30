package dialect

import "bytes"

// SQLiteDialect is the Dialect for SQLite. It uses the LIMIT / OFFSET
// pagination form, which is also what MySQL and modern Postgres accept. It
// exists as a second, deliberately different dialect so the engine-specific
// seam in the store layer is visible and exercised.
type SQLiteDialect struct{}

// Name implements Dialect.
func (SQLiteDialect) Name() string {
	return "sqlite"
}

// Paginate implements Dialect.
func (SQLiteDialect) Paginate(buf *bytes.Buffer) {
	buf.WriteString(" LIMIT :rows_per_page OFFSET :offset")
}

// IsUndefinedTable implements Dialect.
//
// SQLite reports a missing table as a generic SQLITE_ERROR whose message
// varies by driver and version, and no SQLite driver is vendored yet, so no
// mapping is claimed here.
func (SQLiteDialect) IsUndefinedTable(error) bool {
	return false
}

// IsUniqueViolation implements Dialect. See IsUndefinedTable: the extended
// result codes that would identify a unique violation are driver-specific.
func (SQLiteDialect) IsUniqueViolation(error) bool {
	return false
}
