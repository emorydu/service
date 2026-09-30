package dialect

import (
	"bytes"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Set of PostgreSQL SQLSTATE codes this package maps onto the SDK's
// sentinel errors.
// https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	pgUniqueViolation = "23505"
	pgUndefinedTable  = "42P01"
)

// PostgresDialect is the Dialect for PostgreSQL. It uses the SQL:2008 standard
// OFFSET / FETCH NEXT pagination form, which is what the current stores
// in this repository emit today.
type PostgresDialect struct{}

// Name implements Dialect.
func (PostgresDialect) Name() string {
	return "postgres"
}

// Paginate implements Dialect.
func (PostgresDialect) Paginate(buf *bytes.Buffer) {
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
}

// IsUndefinedTable implements Dialect.
func (PostgresDialect) IsUndefinedTable(err error) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	return ok && pgErr.Code == pgUndefinedTable
}

// IsUniqueViolation implements Dialect.
func (PostgresDialect) IsUniqueViolation(err error) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	return ok && pgErr.Code == pgUniqueViolation
}
