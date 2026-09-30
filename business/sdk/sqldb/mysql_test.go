package sqldb

import (
	"strings"
	"testing"
	"time"

	"github.com/emorydu/service/business/sdk/sqldb/dialect"
)

func TestMySQLConfig(t *testing.T) {
	c := mySQLConfig(MySQLConfig{
		User:      "root",
		Password:  "postgres",
		Host:      "localhost",
		Name:      "sales",
		ParseTime: true,
		Timeout:   5 * time.Second,
		TLS:       "skip-verify",
		Params:    map[string]string{"charset": "utf8mb4"},
	})

	if c.User != "root" {
		t.Errorf("User: got %q, want %q", c.User, "root")
	}
	if c.Passwd != "postgres" {
		t.Errorf("Passwd: got %q, want %q", c.Passwd, "postgres")
	}
	if c.Net != "tcp" {
		t.Errorf("Net: got %q, want %q", c.Net, "tcp")
	}
	if c.Addr != "localhost" {
		t.Errorf("Addr: got %q, want %q", c.Addr, "localhost")
	}
	if c.DBName != "sales" {
		t.Errorf("DBName: got %q, want %q", c.DBName, "sales")
	}
	if !c.ParseTime {
		t.Error("ParseTime: got false, want true")
	}
	if c.Timeout != 5*time.Second {
		t.Errorf("Timeout: got %v, want %v", c.Timeout, 5*time.Second)
	}
	if c.TLSConfig != "skip-verify" {
		t.Errorf("TLSConfig: got %q, want %q", c.TLSConfig, "skip-verify")
	}
	if got := c.Params["charset"]; got != "utf8mb4" {
		t.Errorf("Params[charset]: got %q, want %q", got, "utf8mb4")
	}

	// UTC is what the Postgres path asks for with timezone=utc.
	if c.Loc != time.UTC {
		t.Errorf("Loc: got %v, want %v", c.Loc, time.UTC)
	}

	dsn := c.FormatDSN()
	for _, want := range []string{"parseTime=true", "tls=skip-verify", "charset=utf8mb4", "timeout=5s"} {
		if !strings.Contains(dsn, want) {
			t.Errorf("FormatDSN: %q does not contain %q", dsn, want)
		}
	}

	// FormatDSN only emits loc when it is not UTC, so its absence is the
	// proof that the connection reads times as UTC.
	if strings.Contains(dsn, "loc=") {
		t.Errorf("FormatDSN: %q should not pin a non-UTC location", dsn)
	}
}

func TestOpenMySQLRejectsUnknownTLSConfig(t *testing.T) {
	// A TLS name that was never registered is caught while opening rather
	// than at the first query.
	if _, err := OpenMySQL(MySQLConfig{Host: "localhost", Name: "sales", TLS: "not-registered"}); err == nil {
		t.Fatal("OpenMySQL: expected an error for an unregistered TLS config")
	}
}

func TestOpenMySQLDriverName(t *testing.T) {
	// sqlx does not dial here, so this pins the contract dialectFor depends
	// on without needing a server: a connection opened by OpenMySQL must
	// report the driver name "mysql".
	db, err := OpenMySQL(MySQLConfig{Host: "127.0.0.1:1", Name: "sales"})
	if err != nil {
		t.Fatalf("OpenMySQL: %v", err)
	}
	defer db.Close()

	if got := db.DriverName(); got != "mysql" {
		t.Errorf("DriverName: got %q, want %q", got, "mysql")
	}

	if got := dialect.For(db.DriverName()).Name(); got != "mysql" {
		t.Errorf("dialect.For: got %q, want %q", got, "mysql")
	}
}
