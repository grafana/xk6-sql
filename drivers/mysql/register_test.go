package mysql

import (
	"context"
	_ "embed"
	"runtime"
	"testing"

	"github.com/grafana/xk6-sql/sqltest"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

//go:embed testdata/script.js
var script string

func TestIntegration(t *testing.T) { //nolint:paralleltest
	if testing.Short() {
		t.Skip()
	}

	if runtime.GOOS != "linux" {
		t.Skip("Works only on Linux (Testcontainers)")
	}

	ctx := context.Background()

	ctr, err := mysql.Run(ctx, "mysql:8.0.36")

	require.NoError(t, err)
	defer func() { require.NoError(t, ctr.Terminate(ctx)) }()

	conn, err := ctr.ConnectionString(ctx)

	require.NoError(t, err)

	sqltest.RunScript(t, "mysql", conn, script)
}
