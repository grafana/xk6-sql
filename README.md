# k6-plugin-sql

This is a [k6](https://github.com/k6io/k6) extension using the
[xk6](https://github.com/k6io/xk6) system.

| :exclamation: This is a proof of concept, isn't supported by the k6 team, and may break in the future. USE AT YOUR OWN RISK! |
|------|

## Build

To build a `k6` binary with this plugin, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- If you're using SQLite, a build toolchain for your system that includes `gcc` or another C compiler. On Debian and derivatives install the `build-essential` package.
- Git

Then:

1. Install `xk6`:
  ```shell
  go install github.com/k6io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```shell
  xk6 build master \
    --with github.com/imiric/xk6-sql
  ```

  If you're using SQLite, build with CGO enabled:
  ```shell
  CGO_ENABLED=1 xk6 build master \
    --with github.com/imiric/xk6-sql
  ```

  If you're building on Windows and need SQLite support, you'll also need to setup a
  C compiler like [tdm-gcc](https://jmeubank.github.io/tdm-gcc/). Ensure that `gcc`
  is in your `PATH` and run the following two commands:
  ```shell
  set CGO_ENABLED=1
  xk6 build master --with github.com/imiric/xk6-sql
  ```


## Example

```javascript
// script.js
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

  let results = sql.query(db, "SELECT * FROM keyvalues;");
  for (const row of results) {
    console.log(`key: ${row.key}, value: ${row.value}`);
  }
}
```

Result output:

```shell
$ ./k6 run script.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: /tmp/script.js
     output: -

  scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
           * default: 1 iterations for each of 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)

INFO[0000] key: plugin-name, value: k6-plugin-sql        source=console

running (00m00.1s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  00m00.0s/10m0s  1/1 iters, 1 per VU

    █ setup

    █ teardown

    data_received........: 0 B 0 B/s
    data_sent............: 0 B 0 B/s
    iteration_duration...: avg=9.22ms min=19.39µs med=8.86ms max=18.8ms p(90)=16.81ms p(95)=17.8ms
    iterations...........: 1   15.292228/s
```
