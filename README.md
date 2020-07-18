# k6-plugin-template

Plugin example script:

```javascript
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
```

Result output:

```bash
$ ./build.sh && ./k6 run --iterations 1 --plugin=sql.so test.js

          /\      |‾‾|  /‾‾/  /‾/
     /\  /  \     |  |_/  /  / /
    /  \/    \    |      |  /  ‾‾\  
   /          \   |  |‾\  \ | (_) |
  / __________ \  |__|  \__\ \___/ .io

  execution: local
    plugins: SQL
     output: -
     script: test.js

    duration: -,  iterations: 1
         vus: 1,

  execution: local
     script: test.js
     output: -

  scenarios: (100.00%) 1 executors, 1 max VUs, 10m30s max duration (incl. graceful stop):
           * default: 1 iterations shared among 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)


running (00m00.1s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  00m00.0s/10m0s  1/1 shared iters


    ✓ correct key/value

    checks...............: 100.00% ✓ 1 ✗ 0
    data_received........: 0 B     0 B/s
    data_sent............: 0 B     0 B/s
    iteration_duration...: avg=21.53ms min=21.53ms med=21.53ms max=21.53ms p(90)=21.53ms p(95)=21.53ms
    iterations...........: 1       11.747758/s
```
