package sql_test

import (
	"context"
	"io"
	"testing"

	"github.com/grafana/sobek"
	"github.com/grafana/xk6-sql/sql"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/metrics"

	_ "github.com/proullon/ramsql/driver"
)

// TestIntegration performs an integration test creating an ramsql database
func TestIntegration(t *testing.T) {
	t.Parallel()

	rt := setupTestEnv(t)

	_, err := rt.RunString(`
const db = sql.open(driver, "testdb");

db.exec("CREATE TABLE test_table (id integer PRIMARY KEY AUTOINCREMENT, name varchar NOT NULL, value varchar);")

for (let i = 0; i < 5; i++) {
	db.exec("INSERT INTO test_table (name, value) VALUES ('name-" + i + "', 'value-" + i + "');")
}

let all_rows = sql.query(db, "SELECT * FROM test_table;")
if (all_rows.length != 5) {
	throw new Error("Expected all five rows to be returned; got " + all_rows.length)
}

let one_row = sql.query(db, "SELECT * FROM test_table WHERE name = $1;", "name-2");
if (one_row.length != 1) {
	throw new Error("Expected single row to be returned; got " + one_row.length)
}

let no_rows = sql.query(db, "SELECT * FROM test_table WHERE name = $1;", 'bogus-name');
if (no_rows.length != 0) {
	throw new Error("Expected no rows to be returned; got " + no_rows.length)
}

db.close()
`)
	require.NoError(t, err)
}

func setupTestEnv(t *testing.T) *sobek.Runtime {
	t.Helper()

	rt := sobek.New()
	rt.SetFieldNameMapper(common.FieldNameMapper{})

	testLog := logrus.New()
	testLog.AddHook(&testutils.SimpleLogrusHook{
		HookedLevels: []logrus.Level{logrus.WarnLevel},
	})
	testLog.SetOutput(io.Discard)

	state := &lib.State{
		Options: lib.Options{
			SystemTags: metrics.NewSystemTagSet(metrics.TagVU),
		},
		Logger: testLog,
		Tags:   lib.NewVUStateTags(metrics.NewRegistry().RootTagSet()),
	}

	root := sql.New()
	m := root.NewModuleInstance(
		&modulestest.VU{
			RuntimeField: rt,
			InitEnvField: &common.InitEnvironment{},
			CtxField:     context.Background(),
			StateField:   state,
		},
	)

	require.NoError(t, rt.Set("sql", m.Exports().Default))

	root = sql.RegisterModule("ramsql")
	m = root.NewModuleInstance(
		&modulestest.VU{
			RuntimeField: rt,
			InitEnvField: &common.InitEnvironment{},
			CtxField:     context.Background(),
			StateField:   state,
		},
	)

	require.NoError(t, rt.Set("driver", m.Exports().Default))

	return rt
}
