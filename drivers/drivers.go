// Package drivers contains SQL driver initializations.
package drivers

import (
	"github.com/grafana/xk6-sql/sql"

	// Blank imports required for initialization of drivers
	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/microsoft/go-mssqldb/azuread"
)

func init() {
	sql.RegisterDriver("mysql")
	sql.RegisterDriver("postgres")
	sql.RegisterDriver("sqlite3")
	sql.RegisterDriver("sqlserver")
	sql.RegisterDriver("azuresql")
	sql.RegisterDriver("clickhouse")
}
