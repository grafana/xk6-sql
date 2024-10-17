// Package sql is the primary go package of the xk6-sql extension.
// Contains the registration of the xk6-sql extension as a k6 extension.
package sql

import (
	// Blank imports required for initialization of drivers.
	_ "github.com/grafana/xk6-sql/drivers/azuresql"
	_ "github.com/grafana/xk6-sql/drivers/clickhouse"
	_ "github.com/grafana/xk6-sql/drivers/mysql"
	_ "github.com/grafana/xk6-sql/drivers/postgres"
	_ "github.com/grafana/xk6-sql/drivers/sqlite3"
	_ "github.com/grafana/xk6-sql/drivers/sqlserver"

	"github.com/grafana/xk6-sql/sql"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register(sql.ImportPath, sql.New())
}
