import sql from 'k6/x/sql';

const db = sql.open("sqlserver", "./test.db");

export function setup() {
  db.exec(`IF object_id('keyvalues') is null
            CREATE TABLE keyvalues (
            [id] INT IDENTITY PRIMARY KEY,
            [key] varchar(50) NOT NULL,
            [value] varchar(50));`);
}

export function teardown() {
  db.close();
}

export default function () {
  db.exec("INSERT INTO keyvalues ([key], [value]) VALUES('plugin-name', 'k6-plugin-sql');");

  let results = sql.query(db, "SELECT * FROM keyvalues;");
  for (const row of results) {
    console.log(`key: ${row.key}, value: ${row.value}`);
  }
}
