// Package sql provides a javascript module for performing SQL actions against relational databases
package sql

import (
	"crypto/tls"
	"crypto/x509"
	dbsql "database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/grafana/sobek"
	"os"
	"strings"
	// Blank imports required for initialization of drivers
	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/microsoft/go-mssqldb/azuread"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib/netext"
)

// TLSVersions is a map of TLS versions to their numeric values.
var TLSVersions map[string]uint16

func init() {
	// Initialize the TLS versions map.
	TLSVersions = map[string]uint16{
		netext.TLS_1_0: tls.VersionTLS10,
		netext.TLS_1_1: tls.VersionTLS11,
		netext.TLS_1_2: tls.VersionTLS12,
		netext.TLS_1_3: tls.VersionTLS13,
	}

	modules.Register("k6/x/sql", new(RootModule))
}

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create `k6/x/sql` module instances for each VU.
type RootModule struct{}

// SQL represents an instance of the SQL module for every VU.
type SQL struct {
	vu        modules.VU
	exports   *sobek.Object
	tlsConfig TLSConfig
}

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &SQL{}
)

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	runtime := vu.Runtime()

	moduleInstance := &SQL{
		vu:        vu,
		exports:   runtime.NewObject(),
		tlsConfig: TLSConfig{},
	}
	// Export constants to the JS code.
	moduleInstance.defineConstants()

	mustExport := func(name string, value interface{}) {
		if err := moduleInstance.exports.Set(name, value); err != nil {
			common.Throw(runtime, err)
		}
	}
	mustExport("loadTLS", moduleInstance.LoadTLS)
	mustExport("open", moduleInstance.Open)
	mustExport("query", moduleInstance.Query)

	return moduleInstance
}

// Exports implements the modules.Instance interface and returns the exports
// of the JS module.
func (sql *SQL) Exports() modules.Exports {
	return modules.Exports{
		Default: sql.exports,
	}
}

// KeyValue is a simple key-value pair.
type KeyValue map[string]interface{}

func contains(array []string, element string) bool {
	for _, item := range array {
		if item == element {
			return true
		}
	}
	return false
}

// defineConstants defines the constants that can be used in the JS code.
func (sql *SQL) defineConstants() {
	runtime := sql.vu.Runtime()
	mustAddProp := func(name string, val interface{}) {
		err := sql.exports.DefineDataProperty(
			name, runtime.ToValue(val), sobek.FLAG_FALSE, sobek.FLAG_FALSE, sobek.FLAG_TRUE,
		)
		if err != nil {
			common.Throw(runtime, err)
		}
	}

	// TLS versions
	mustAddProp("TLS_1_0", netext.TLS_1_0)
	mustAddProp("TLS_1_1", netext.TLS_1_1)
	mustAddProp("TLS_1_2", netext.TLS_1_2)
	mustAddProp("TLS_1_3", netext.TLS_1_3)
}

type TLSConfig struct {
	EnableTLS             bool   `json:"enableTLS"`
	InsecureSkipTLSverify bool   `json:"insecureSkipTLSverify"`
	MinVersion            string `json:"minVersion"`
	CAcertFile            string `json:"caCertFile"`
	ClientCertFile        string `json:"clientCertFile"`
	ClientKeyFile         string `json:"clientKeyFile"`
}

// LoadTLS loads the TLS configuration for the SQL module.
func (sql *SQL) LoadTLS(params map[string]interface{}) {
	runtime := sql.vu.Runtime()
	var tlsConfig *TLSConfig
	if b, err := json.Marshal(params); err != nil {
		common.Throw(runtime, err)
	} else {
		if err := json.Unmarshal(b, &tlsConfig); err != nil {
			common.Throw(runtime, err)
		}
	}
	sql.tlsConfig = *tlsConfig
}

// Open establishes a connection to the specified database type using
// the provided connection string.
func (sql *SQL) Open(database string, connectionString string) (*dbsql.DB, error) {
	supportedDatabases := []string{"mysql", "postgres", "sqlite3", "sqlserver", "azuresql", "clickhouse"}
	if !contains(supportedDatabases, database) {
		return nil, fmt.Errorf("database %s is not supported", database)
	}

	if database == "mysql" && sql.tlsConfig.EnableTLS {
		tlsConfigKey := "custom"
		if err := registerTLS(tlsConfigKey, sql.tlsConfig); err != nil {
			return nil, err
		}
		connectionString = prefixConnectionString(connectionString, tlsConfigKey)
	}

	db, err := dbsql.Open(database, connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// prefixConnectionString prefixes the connection string with the TLS configuration key.
func prefixConnectionString(connectionString string, tlsConfigKey string) string {
	var separator string
	if strings.Contains(connectionString, "?") {
		separator = "&"
	} else {
		separator = "?"
	}
	return fmt.Sprintf("%s%stls=%s", connectionString, separator, tlsConfigKey)
}

// registerTLS loads the ca-cert and registers the TLS configuration with the MySQL driver.
func registerTLS(tlsConfigKey string, tlsConfig TLSConfig) error {
	rootCAs := x509.NewCertPool()
	{
		pem, err := os.ReadFile(tlsConfig.CAcertFile)
		if err != nil {
			return err
		}
		if ok := rootCAs.AppendCertsFromPEM(pem); !ok {
			return fmt.Errorf("failed to append PEM")
		}
	}
	clientCerts := make([]tls.Certificate, 0, 1)
	{
		certs, err := tls.LoadX509KeyPair(tlsConfig.ClientCertFile, tlsConfig.ClientKeyFile)
		if err != nil {
			return err
		}
		clientCerts = append(clientCerts, certs)
	}

	mysqlTLSConfig := &tls.Config{
		RootCAs:            rootCAs,
		Certificates:       clientCerts,
		MinVersion:         TLSVersions[tlsConfig.MinVersion],
		InsecureSkipVerify: tlsConfig.InsecureSkipTLSverify,
	}
	if err := mysql.RegisterTLSConfig(tlsConfigKey, mysqlTLSConfig); err != nil {
		return err
	}
	return nil
}

// Query executes the provided query string against the database, while
// providing results as a slice of KeyValue instance(s) if available.
func (*SQL) Query(db *dbsql.DB, query string, args ...interface{}) ([]KeyValue, error) {
	rows, err := db.Query(query, args...)
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
