// Package sqlite3 contains SQLite3 driver registration for xk6-sql.
package sqlite3

import (
	"github.com/grafana/xk6-sql/sql"

	// Blank import required for initialization of driver.
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	sql.RegisterModule("sqlite3")
}
