package sql

import (
	"testing"

	"github.com/grafana/sobek"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) { //nolint: paralleltest
	mod, ok := New().NewModuleInstance(nil).(*module)

	require.True(t, ok)

	driver := RegisterDriver("ramsql")
	require.NotNil(t, driver)

	db, err := mod.Open(driver, "", nil)

	require.NoError(t, err)
	require.NotNil(t, db)

	_, err = mod.Open(sobek.New().ToValue("foo"), "testdb", nil) // not a Symbol

	require.Error(t, err)

	_, err = mod.Open(sobek.NewSymbol("ramsql"), "testdb", nil) // not a registered Symbol

	require.Error(t, err)
}
