import sql from 'k6/x/sql';

const db = sql.open("sqlite3", "./test.db");

export function setup() {
  db.exec(`CREATE TABLE IF NOT EXISTS keyvalues (
           id integer PRIMARY KEY AUTOINCREMENT,
           key varchar NOT NULL,
           value varchar);`);
}

export function teardown() {
  db.close();
}

export default function () {
  db.exec("INSERT INTO keyvalues (key, value) VALUES('plugin-name', 'k6-plugin-sql');");

  let results = sql.query(db, "SELECT * FROM keyvalues WHERE key = $1;", 'plugin-name');
  for (const row of results) {
    console.log(`key: ${row.key}, value: ${row.value}`);
  }
}
