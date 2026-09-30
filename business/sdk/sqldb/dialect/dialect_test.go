package dialect_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/emorydu/service/business/sdk/sqldb/dialect"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name    string
		dialect dialect.Dialect
		want    string
	}{
		{
			name:    "postgres",
			dialect: dialect.Postgres,
			want:    " OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY",
		},
		{
			name:    "mysql",
			dialect: dialect.MySQL,
			want:    " LIMIT :rows_per_page OFFSET :offset",
		},
		{
			name:    "sqlite",
			dialect: dialect.SQLite,
			want:    " LIMIT :rows_per_page OFFSET :offset",
		},
		{
			name:    "generic",
			dialect: dialect.Generic,
			want:    " LIMIT :rows_per_page OFFSET :offset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.dialect.Paginate(&buf)

			if got := buf.String(); got != tt.want {
				t.Errorf("Paginate: got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestName(t *testing.T) {
	tests := []struct {
		dialect dialect.Dialect
		want    string
	}{
		{dialect.Postgres, "postgres"},
		{dialect.MySQL, "mysql"},
		{dialect.SQLite, "sqlite"},
		{dialect.Generic, "generic"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.dialect.Name(); got != tt.want {
				t.Errorf("Name: got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFor(t *testing.T) {
	tests := []struct {
		driverName string
		want       string
	}{
		{"pgx", "postgres"},
		{"pgx/v5", "postgres"},
		{"postgres", "postgres"},
		{"mysql", "mysql"},
		{"sqlite3", "sqlite"},
		{"sqlite", "sqlite"},
		{"", "generic"},
		{"bogus", "generic"},
	}

	for _, tt := range tests {
		t.Run(tt.driverName, func(t *testing.T) {
			if got := dialect.For(tt.driverName).Name(); got != tt.want {
				t.Errorf("For(%q): got %q, want %q", tt.driverName, got, tt.want)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name            string
		dialect         dialect.Dialect
		err             error
		undefinedTable  bool
		uniqueViolation bool
	}{
		{
			name:           "mysql undefined table",
			dialect:        dialect.MySQL,
			err:            &mysql.MySQLError{Number: 1146},
			undefinedTable: true,
		},
		{
			name:            "mysql duplicate entry",
			dialect:         dialect.MySQL,
			err:             &mysql.MySQLError{Number: 1062},
			uniqueViolation: true,
		},
		{
			name:           "mysql wrapped by fmt",
			dialect:        dialect.MySQL,
			err:            fmt.Errorf("insert order: %w", &mysql.MySQLError{Number: 1146}),
			undefinedTable: true,
		},
		{
			name:            "mysql inside errors.Join",
			dialect:         dialect.MySQL,
			err:             errors.Join(errors.New("context"), &mysql.MySQLError{Number: 1062}),
			uniqueViolation: true,
		},
		{
			name:           "postgres undefined table",
			dialect:        dialect.Postgres,
			err:            &pgconn.PgError{Code: "42P01"},
			undefinedTable: true,
		},
		{
			name:            "postgres duplicate entry",
			dialect:         dialect.Postgres,
			err:             &pgconn.PgError{Code: "23505"},
			uniqueViolation: true,
		},
		{
			name:    "mysql error is not a postgres error",
			dialect: dialect.Postgres,
			err:     &mysql.MySQLError{Number: 1146},
		},
		{
			name:    "postgres error is not a mysql error",
			dialect: dialect.MySQL,
			err:     &pgconn.PgError{Code: "23505"},
		},
		{
			name:    "mysql error with unrelated number",
			dialect: dialect.MySQL,
			err:     &mysql.MySQLError{Number: 1064},
		},
		{
			name:    "sqlite claims neither",
			dialect: dialect.SQLite,
			err:     &mysql.MySQLError{Number: 1062},
		},
		{
			name:    "generic claims neither",
			dialect: dialect.Generic,
			err:     &mysql.MySQLError{Number: 1062},
		},
		{
			name:    "nil error",
			dialect: dialect.MySQL,
		},
		{
			name:    "unrelated error",
			dialect: dialect.MySQL,
			err:     errors.New("boom"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dialect.IsUndefinedTable(tt.err); got != tt.undefinedTable {
				t.Errorf("IsUndefinedTable: got %v, want %v", got, tt.undefinedTable)
			}
			if got := tt.dialect.IsUniqueViolation(tt.err); got != tt.uniqueViolation {
				t.Errorf("IsUniqueViolation: got %v, want %v", got, tt.uniqueViolation)
			}
		})
	}
}
