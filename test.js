/*

This is a k6 test script that imports the k6-plugin-sql.

*/

import { check } from 'k6';
import { open, query } from 'k6-plugin/sql';  // import sql plugin

export default function () {
    let db = open("sqlite3", "./test.sqlite");
    db.exec(`CREATE TABLE IF NOT EXISTS keyvalues (
        id integer PRIMARY KEY AUTOINCREMENT,
        key varchar NOT NULL,
        value varchar);`);
    db.exec("INSERT INTO keyvalues (key, value) VALUES('plugin-name', 'k6-plugin-sql');");
    let results = query(db, "SELECT * FROM keyvalues;");
    results.forEach(row => {
        check(row, {
            "correct key/value": r => r['key'] == 'plugin-name' && r['value'] == 'k6-plugin-sql'
        });
    });
    db.close();
}
