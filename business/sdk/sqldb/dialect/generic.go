package dialect

import "bytes"

// GenericDialect is the Dialect for a database/sql driver name this package
// does not recognize. It never claims an error is a known constraint failure,
// so the caller falls back to the raw driver error, and it paginates with
// LIMIT / OFFSET, the form the widest range of engines accepts.
type GenericDialect struct{}

// Name implements Dialect.
func (GenericDialect) Name() string {
	return "generic"
}

// Paginate implements Dialect.
func (GenericDialect) Paginate(buf *bytes.Buffer) {
	buf.WriteString(" LIMIT :rows_per_page OFFSET :offset")
}

// IsUndefinedTable implements Dialect.
func (GenericDialect) IsUndefinedTable(error) bool {
	return false
}

// IsUniqueViolation implements Dialect.
func (GenericDialect) IsUniqueViolation(error) bool {
	return false
}
