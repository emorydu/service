package dialect

import (
	"bytes"
	"errors"

	"github.com/go-sql-driver/mysql"
)

// Set of MySQL server error numbers this package maps onto the SDK's sentinel
// errors. The driver reports these numbers on *mysql.MySQLError but exports no
// named constant for them, so they are spelled out here.
// https://dev.mysql.com/doc/mysql-errors/8.4/en/server-error-reference.html
const (
	mySQLErrDupEntry    = 1062 // ER_DUP_ENTRY
	mySQLErrNoSuchTable = 1146 // ER_NO_SUCH_TABLE
)

// MySQLDialect is the Dialect for MySQL and MariaDB. It uses the LIMIT /
// OFFSET pagination form, which MySQL shares with SQLite and modern Postgres
// but which Postgres still accepts only in its own dialect.
type MySQLDialect struct{}

// Name implements Dialect.
func (MySQLDialect) Name() string {
	return "mysql"
}

// Paginate implements Dialect.
func (MySQLDialect) Paginate(buf *bytes.Buffer) {
	buf.WriteString(" LIMIT :rows_per_page OFFSET :offset")
}

// IsUndefinedTable implements Dialect.
func (MySQLDialect) IsUndefinedTable(err error) bool {
	myErr, ok := errors.AsType[*mysql.MySQLError](err)
	return ok && myErr.Number == mySQLErrNoSuchTable
}

// IsUniqueViolation implements Dialect.
func (MySQLDialect) IsUniqueViolation(err error) bool {
	myErr, ok := errors.AsType[*mysql.MySQLError](err)
	return ok && myErr.Number == mySQLErrDupEntry
}
