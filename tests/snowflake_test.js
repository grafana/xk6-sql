import sql from 'k6/x/sql';

// The second argument is a Snowflake connection string, e.g.
// myuser[:mypassword]@myorganization-myaccount.snowflakecomputing.com/mydatabase?schema=myschema&warehouse=mywarehouse'
const db = sql.open('snowflake', '');

export function setup() {
  db.exec(`
    CREATE TABLE IF NOT EXISTS keyvalues (
      id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
      \`key\` VARCHAR(50) NOT NULL,
      value VARCHAR(50) NULL
    );
  `);
}

export function teardown() {
  db.close();
}

export default function () {
  db.exec("INSERT INTO keyvalues (`key`, value) VALUES('plugin-name', 'k6-plugin-sql');");

  let results = sql.query(db, "SELECT * FROM keyvalues WHERE `key` = ?;", 'plugin-name');
  for (const row of results) {
    console.log(`key: ${String.fromCharCode(...row.key.toString().split(','))}, value: ${String.fromCharCode(...row.value.toString().split(','))}`);
  }
}
