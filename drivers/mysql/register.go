// Package mysql contains MySQL driver registration for xk6-sql.
package mysql

import (
	"github.com/grafana/xk6-sql/sql"

	// Blank import required for initialization of driver.
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	sql.RegisterModule("mysql")
}
