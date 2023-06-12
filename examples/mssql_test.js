import sql from 'k6/x/sql';

// The second argument is a MS SQL connection string, e.g.
// Server=127.0.0.1;Database=myDB;User Id=myUser;Password=myPassword;
const db = sql.open('sqlserver', '');

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

  let results = sql.query(db, "SELECT * FROM keyvalues WHERE [key] = @p1;", 'plugin-name');
  for (const row of results) {
    console.log(`key: ${row.key}, value: ${row.value}`);
  }
}
