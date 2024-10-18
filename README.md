# xk6-sql

**Use SQL databases from k6 tests.**

xk6-sql is a [Grafana k6 extension](https://grafana.com/docs/k6/latest/extensions/) that enables the use of SQL databases in [k6](https://grafana.com/docs/k6/latest/) tests.

## Usage

To use the xk6-sql API, the `k6/x/sql` module and the driver module corresponding to the database type should be imported. In the example below, `k6/x/sql/driver/ramsql` is the RamSQL database driver module.

The driver module exports a driver ID. This driver identifier should be used to identify the database driver to be used in the API functions.

**example**

```javascript file=examples/example.js
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
```

<details>
<summary><b>output</b></summary>

```bash file=examples/example.txt

         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: examples/example.js
        output: -

     scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
              * default: 1 iterations for each of 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)

time="2024-10-18T09:06:52+02:00" level=info msg="name: extension-name, value: xk6-foo" source=console

     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=496.46µs min=496.46µs med=496.46µs max=496.46µs p(90)=496.46µs p(95)=496.46µs
     iterations...........: 1   550.030197/s


running (00m00.0s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [ 100% ] 1 VUs  00m00.0s/10m0s  1/1 iters, 1 per VU
```

</details>

## Build

The [xk6](https://github.com/grafana/xk6) build tool can be used to build a k6 that will include **xk6-sql** extension and database drivers.

> [!IMPORTANT]
> In the command line bellow, **xk6-sql-driver-ramsql** is just an example, it should be replaced with the database driver extension you want to use.
> For example use `--with github.com/grafana/xk6-sql-driver-mysql` to access MySQL databases.

```bash
xk6 build --with github.com/grafana/xk6-sql@latest --with github.com/grafana/xk6-sql-driver-ramsql
```

For more build options and how to use xk6, check out the [xk6 documentation](https://github.com/grafana/xk6).

Supported RDBMSs includes: `mysql`, `postgres`, `sqlite3`, `sqlserver`, `azuresql`, `clickhouse`.

Check the [xk6-sql-driver GitHub topic](https://github.com/topics/xk6-sql-driver) to discover database driver extensions.

## Drivers

To use the xk6-sql extension, one or more database driver extensions should also be embedded. Database driver extension names typically start with the prefix `xk6-sql-driver-` followed by the name of the database, for example `xk6-sql-driver-mysql` is the name of the MySQL database driver extension.

For easier discovery, the `xk6-sql-driver` topic is included in the database driver extensions repository. The [xk6-sql-driver GitHub topic search](https://github.com/topics/xk6-sql-driver) therefore lists the available driver extensions.

### Create driver

Check the [grafana/xk6-sql-driver-ramsql](https://github.com/grafana/xk6-sql-driver-ramsql) template repository to create a new driver extension. This is a working driver extension with instructions in its README for customization.

[Postgres driver extension](https://github.com/grafana/xk6-sql-driver-postgres) and [MySQL driver extension](https://github.com/grafana/xk6-sql-driver-mysql) are also good examples.

## Feedback

If you find the **xk6-sql** extension useful, please star the repo. The number of stars will affect the time allocated for maintenance.

[![Stargazers over time](https://starchart.cc/grafana/xk6-sql.svg?variant=adaptive)](https://starchart.cc/grafana/xk6-sql)

## Contribute

If you want to contribute or help with the development of **xk6-sql**, start by reading [CONTRIBUTING.md](CONTRIBUTING.md). 
