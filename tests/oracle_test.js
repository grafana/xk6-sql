import sql from 'k6/x/sql';

// The second argument is a Oracle connection string, e.g.
// `user="myuser" password="mypass" connectString="127.0.0.1:1521/mydb"`
const db = sql.open('godror', '');

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

  let results = sql.query(db, "SELECT * FROM keyvalues WHERE `key` = :1", 'plugin-name');
  for (const row of results) {
    // Convert array of ASCII integers into strings. See https://github.com/grafana/xk6-sql/issues/12
    console.log(`key: ${String.fromCharCode(...row.key)}, value: ${String.fromCharCode(...row.value)}`);
  }
}
