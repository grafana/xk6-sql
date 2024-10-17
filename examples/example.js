import sql from "k6/x/sql";

// ramsql is hypothetical, the actual driver name should be used instead.
import driver from "k6/x/sql/driver/ramsql";

const db = sql.open(driver, "test_db");

export function setup() {
  db.exec(`CREATE TABLE IF NOT EXISTS namevalue (
             id INTEGER PRIMARY KEY AUTOINCREMENT,
             name VARCHAR NOT NULL,
             value VARCHAR
           );`);
}

export function teardown() {
  db.close();
}

export default function () {
  db.exec("INSERT INTO namevalue (name, value) VALUES('extension-name', 'xk6-foo');");

  let results = sql.query(db, "SELECT * FROM namevalue WHERE name = $1;", "extension-name");
  for (const row of results) {
    console.log(`name: ${row.name}, value: ${row.value}`);
  }
}
