package sqldb

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQLConfig is the required properties to open a MySQL database. It is
// deliberately separate from Config, whose Schema and DisableTLS fields only
// have meaning for Postgres.
//
// Schema has no MySQL counterpart: in MySQL a schema *is* a database, and no
// connection parameter selects a namespace inside one. A store ported from
// Postgres that relied on search_path must either move that value into Name
// or qualify its table names explicitly.
//
// TLS likewise does not share Config's semantics. The zero value here means
// no TLS at all, whereas Config's zero value asks for an encrypted but
// unverified connection. To approximate sslmode=require, set TLS to
// "skip-verify"; to verify the server certificate, set it to "true".
type MySQLConfig struct {
	User         string
	Password     string
	Host         string // host:port; the driver appends :3306 when omitted
	Name         string // the database name
	MaxIdleConns int
	MaxOpenConns int
	ParseTime    bool // required for DATE and DATETIME to scan into time.Time
	Timeout      time.Duration
	TLS          string // "", "true", "false", "skip-verify", "preferred", or a name registered with mysql.RegisterTLSConfig
	Params       map[string]string
}

// OpenMySQL knows how to open a MySQL database connection based on the
// configuration.
func OpenMySQL(cfg MySQLConfig) (*sqlx.DB, error) {

	// Going through the connector rather than a DSN string keeps the password
	// out of a concatenated string and validates the configuration up front,
	// so a misspelled TLS name fails here instead of at the first query.
	connector, err := mysql.NewConnector(mySQLConfig(cfg))
	if err != nil {
		return nil, err
	}

	db := sqlx.NewDb(sql.OpenDB(connector), "mysql")
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

// mySQLConfig translates the SDK's configuration into the driver's, which is
// the only other place the MySQL driver is named.
func mySQLConfig(cfg MySQLConfig) *mysql.Config {
	c := mysql.NewConfig()
	c.User = cfg.User
	c.Passwd = cfg.Password
	c.Net = "tcp"
	c.Addr = cfg.Host
	c.DBName = cfg.Name
	c.TLSConfig = cfg.TLS
	c.Timeout = cfg.Timeout
	c.Params = cfg.Params
	c.ParseTime = cfg.ParseTime

	// NewConfig already sets Loc to time.UTC, which is what the Postgres path
	// asks for with timezone=utc. FormatDSN omits the loc parameter while it
	// stays UTC, so the alignment is implicit rather than written out.

	return c
}
