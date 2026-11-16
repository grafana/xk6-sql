package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"go.k6.io/k6/v2/js/modulestest"
)

func TestConvertValue(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		typeName string
		in       any
		want     any
	}{
		{"text", "VARCHAR", []byte("hello"), "hello"},
		{"signed integer", "BIGINT", []byte("926"), int64(926)},
		{"negative integer", "INT", []byte("-7"), int64(-7)},
		{"unsigned integer", "UNSIGNED BIGINT", []byte("18446744073709551615"), uint64(18446744073709551615)},
		{"double", "DOUBLE", []byte("2.5"), 2.5},
		{"decimal stays exact", "DECIMAL", []byte("12.340"), "12.340"},
		{"unparsable integer", "INT", []byte("abc"), "abc"},
		{"blob", "BLOB", []byte{0x00, 0xff}, []byte{0x00, 0xff}},
		{"varbinary", "VARBINARY", []byte{0x01}, []byte{0x01}},
		{"bytea", "BYTEA", []byte{0x01}, []byte{0x01}},
		{"null", "VARCHAR", nil, nil},
		{"native integer", "BIGINT", int64(5), int64(5)},
		{"native string", "VARCHAR", "x", "x"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, convertValue(tc.typeName, tc.in))
		})
	}
}

// bytesDriver is a minimal driver that returns every value as []byte,
// the way the MySQL driver does for text protocol results.
type bytesDriver struct{}

func (bytesDriver) Open(string) (driver.Conn, error) { return bytesConn{}, nil }

type bytesConn struct{}

func (bytesConn) Prepare(string) (driver.Stmt, error) { return bytesStmt{}, nil }
func (bytesConn) Close() error                        { return nil }
func (bytesConn) Begin() (driver.Tx, error)           { return nil, driver.ErrSkip }

type bytesStmt struct{}

func (bytesStmt) Close() error                               { return nil }
func (bytesStmt) NumInput() int                              { return -1 }
func (bytesStmt) Exec([]driver.Value) (driver.Result, error) { return nil, driver.ErrSkip }
func (bytesStmt) Query([]driver.Value) (driver.Rows, error)  { return &bytesRows{}, nil }

type bytesRows struct{ done bool }

func (*bytesRows) Columns() []string { return []string{"count", "name", "data"} }
func (*bytesRows) Close() error      { return nil }

func (r *bytesRows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}

	r.done = true
	dest[0], dest[1], dest[2] = []byte("926"), []byte("alice"), []byte{0x00, 0x01}

	return nil
}

func (*bytesRows) ColumnTypeDatabaseTypeName(index int) string {
	return []string{"BIGINT", "VARCHAR", "BLOB"}[index]
}

func TestQueryConvertsRawBytes(t *testing.T) {
	t.Parallel()

	sql.Register("xk6sqlbytes", bytesDriver{})

	rt := modulestest.NewRuntime(t)
	mod, ok := New().NewModuleInstance(rt.VU).(*module)
	require.True(t, ok)

	db, err := mod.Open(RegisterDriver("xk6sqlbytes"), "", nil)
	require.NoError(t, err)

	db.ctx = context.Background

	rows, err := db.Query("SELECT 1")
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, KeyValue{"count": int64(926), "name": "alice", "data": []byte{0x00, 0x01}}, rows[0])
}
