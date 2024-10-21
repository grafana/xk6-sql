// Package sql provides a javascript module for performing SQL actions against relational databases.
package sql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"
)

// ImportPath contains module's JavaScript import path.
const ImportPath = "k6/x/sql"

// New creates a new instance of the extension's JavaScript module.
func New() modules.Module {
	return new(rootModule)
}

// rootModule is the global module object type. It is instantiated once per test
// run and will be used to create `k6/x/sql` module instances for each VU.
type rootModule struct{}

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*rootModule) NewModuleInstance(_ modules.VU) modules.Instance {
	instance := &module{}

	instance.exports.Default = instance
	instance.exports.Named = map[string]interface{}{
		"open": instance.Open,
	}

	return instance
}

// module represents an instance of the JavaScript module for every VU.
type module struct {
	exports modules.Exports
}

// Exports is representation of ESM exports of a module.
func (mod *module) Exports() modules.Exports {
	return mod.exports
}

// KeyValue is a simple key-value pair.
type KeyValue map[string]interface{}

// open establishes a connection to the specified database type using
// the provided connection string.
func (mod *module) Open(driverID sobek.Value, connectionString string) (*Database, error) {
	driverSym, ok := driverID.(*sobek.Symbol)
	if !ok {
		return nil, fmt.Errorf("%w: invalid driver parameter type", errUnsupportedDatabase)
	}

	registered, database := lookupDriver(driverSym)
	if !registered {
		return nil, fmt.Errorf("%w: %s", errUnsupportedDatabase, database)
	}

	db, err := sql.Open(database, connectionString)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

// Database is a database handle representing a pool of zero or more underlying connections.
type Database struct {
	db *sql.DB
}

// Query executes a query that returns rows, typically a SELECT.
func (dbase *Database) Query(query string, args ...interface{}) ([]KeyValue, error) {
	rows, err := dbase.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	result := make([]KeyValue, 0)

	for rows.Next() {
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		data := make(KeyValue, len(cols))
		for i, colName := range cols {
			data[colName] = *valuePtrs[i].(*interface{}) //nolint:forcetypeassert
		}
		result = append(result, data)
	}

	return result, nil
}

// Exec a query without returning any rows.
func (dbase *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return dbase.db.Exec(query, args...)
}

// Close the database and prevents new queries from starting.
func (dbase *Database) Close() error {
	return dbase.db.Close()
}

var errUnsupportedDatabase = errors.New("unsupported database")
