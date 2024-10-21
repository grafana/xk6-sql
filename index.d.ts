/**
 * **Use SQL databases from k6 tests.**
 *
 * xk6-sql is a [Grafana k6 extension](https://grafana.com/docs/k6/latest/extensions/) that enables
 * the use of SQL databases in [k6](https://grafana.com/docs/k6/latest/) tests.
 *
 * @example
 *  ```ts file=examples/example.js
 *  import sql from "k6/x/sql";
 *
 *  // the actual database driver should be used instead of ramsql
 *  import driver from "k6/x/sql/driver/ramsql";
 *
 *  const db = sql.open(driver, "roster_db");
 *
 *  export function setup() {
 *    db.exec(`
 *      CREATE TABLE IF NOT EXISTS roster
 *        (
 *          id INTEGER PRIMARY KEY AUTOINCREMENT,
 *          given_name VARCHAR NOT NULL,
 *          family_name VARCHAR NOT NULL
 *        );
 *    `);
 *  }
 *
 *  export function teardown() {
 *    db.close();
 *  }
 *
 *  export default function () {
 *    let result = db.exec(`
 *      INSERT INTO roster
 *        (given_name, family_name)
 *      VALUES
 *        ('Peter', 'Pan'),
 *        ('Wendy', 'Darling'),
 *        ('Tinker', 'Bell'),
 *        ('James', 'Hook');
 *    `);
 *    console.log(`${result.rowsAffected()} rows inserted`);
 *
 *    let rows = db.query("SELECT * FROM roster WHERE given_name = $1;", "Peter");
 *    for (const row of rows) {
 *      console.log(`${row.family_name}, ${row.given_name}`);
 *    }
 *  }
 *  ```
 *
 * @module sql
 */
export as namespace sql;

/**
 * Open a database specified by database driver identifier Symbol and a driver-specific data source name,
 * usually consisting of at least a database name and connection information.
 *
 * @param dirverID driver identification symbol, the default export of the driver module
 * @param dataSourceName driver-specific data source name, like a database name
 *
 * @example
 *  ```ts file=examples/example.js
 *  import sql from "k6/x/sql";
 *
 *  // the actual database driver should be used instead of ramsql
 *  import driver from "k6/x/sql/driver/ramsql";
 *
 *  const db = sql.open(driver, "roster_db");
 * ```
 */
export function open(dirverID: Symbol, dataSourceName: String): Database;

/**
 * Database is a database handle representing a pool of zero or more underlying connections.
 *
 * @example
 *  ```ts file=examples/example.js
 *  import sql from "k6/x/sql";
 *
 *  // the actual database driver should be used instead of ramsql
 *  import driver from "k6/x/sql/driver/ramsql";
 *
 *  const db = sql.open(driver, "roster_db");
 * ```
 */
export interface Database {
  /**
   * Close the database and prevents new queries from starting.
   *
   * Close waits for all queries that have started processing on the server to finish.
   *
   * @example
   *  ```ts file=examples/example.js
   *  import sql from "k6/x/sql";
   *
   *  // the actual database driver should be used instead of ramsql
   *  import driver from "k6/x/sql/driver/ramsql";
   *
   *  const db = sql.open(driver, "roster_db");
   *
   *  export function teardown() {
   *    db.close();
   *  }
   * ```
   */
  close(): void;

  /**
   * Execute a query without returning any rows.
   *
   * @param query the query to execute
   * @param args placeholder parameters in the query
   * @returns summary of the executed SQL commands
   * @example
   *  ```ts file=examples/example.js
   *  import sql from "k6/x/sql";
   *
   *  // the actual database driver should be used instead of ramsql
   *  import driver from "k6/x/sql/driver/ramsql";
   *
   *  const db = sql.open(driver, "roster_db");
   *
   *  export function setup() {
   *    db.exec(`
   *      CREATE TABLE IF NOT EXISTS roster
   *        (
   *          id INTEGER PRIMARY KEY AUTOINCREMENT,
   *          given_name VARCHAR NOT NULL,
   *          family_name VARCHAR NOT NULL
   *        );
   *    `);
   *
   *    let result = db.exec(`
   *      INSERT INTO roster
   *        (given_name, family_name)
   *      VALUES
   *        ('Peter', 'Pan'),
   *        ('Wendy', 'Darling'),
   *        ('Tinker', 'Bell'),
   *        ('James', 'Hook');
   *    `);
   *    console.log(`${result.rowsAffected()} rows inserted`);
   *  }
   * ```
   */
  exec(query: string, ...args: any[]): Result;
  /**
   * Query executes a query that returns rows, typically a SELECT.
   * @param query the query to execute
   * @param args placeholder parameters in the query
   * @returns rows of the query result
   * @example
   *  ```ts file=examples/example.js
   *  import sql from "k6/x/sql";
   *
   *  // the actual database driver should be used instead of ramsql
   *  import driver from "k6/x/sql/driver/ramsql";
   *
   *  const db = sql.open(driver, "roster_db");
   *
   *  export default function () {
   *    let rows = db.query("SELECT * FROM roster WHERE given_name = $1;", "Peter");
   *    for (const row of results) {
   *      console.log(`${row.family_name}, ${row.given_name}`);
   *    }
   *  }
   * ```
   */
  query(query: string, ...args: any[]): Row[];
}

/**
 * An object containing a single row of the query result.
 */
export interface Row {
  /**
   * Each column has a property with the same name as the column name.
   * The value of the property contains the value of the given column in the current row.
   *
   * The value of the property is automatically mapped to the corresponding JavaScript type.
   *
   * @example
   *  ```ts file=examples/example.js
   *  import sql from "k6/x/sql";
   *
   *  // the actual database driver should be used instead of ramsql
   *  import driver from "k6/x/sql/driver/ramsql";
   *
   *  const db = sql.open(driver, "roster_db");
   *
   *  export default function () {
   *    let rows = db.query("SELECT * FROM roster WHERE given_name = $1;", "Peter");
   *    for (const row of results) {
   *      console.log(`${row.family_name}, ${row.given_name}`);
   *    }
   *  }
   * ```
   */
  [key: string]: unknown;
}

/**
 * A Result summarizes an executed SQL command.
 * @example
 *  ```ts file=examples/example.js
 *  import sql from "k6/x/sql";
 *
 *  // the actual database driver should be used instead of ramsql
 *  import driver from "k6/x/sql/driver/ramsql";
 *
 *  const db = sql.open(driver, "roster_db");
 *
 *  export function setup() {
 *    let result = db.exec(`
 *      INSERT INTO roster
 *        (given_name, family_name)
 *      VALUES
 *        ('Peter', 'Pan'),
 *        ('Wendy', 'Darling'),
 *        ('Tinker', 'Bell'),
 *        ('James', 'Hook');
 *    `);
 *    console.log(`${result.rowsAffected()} rows inserted`);
 *  }
 * ```
 */
export interface Result {
  /**
   * Returns the integer generated by the database
   * in response to a command. Typically this will be from an
   * "auto increment" column when inserting a new row. Not all
   * databases support this feature, and the syntax of such
   * statements varies.
   */
  lastInsertId(): number;
  /**
   * Returns the number of rows affected by an
   * update, insert, or delete. Not every database or database
   * driver may support this.
   */
  rowsAffected(): number;
}
