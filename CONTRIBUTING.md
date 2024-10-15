# How to contribute to xk6-sql

Thank you for your interest in contributing to **xk6-sql**!

Before you begin, make sure to familiarize yourself with the [Code of Conduct](CODE_OF_CONDUCT.md). If you've previously contributed to other open source project, you may recognize it as the classic [Contributor Covenant](https://contributor-covenant.org/).

## Prerequisites

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git
- If you're using SQLite, a build toolchain for your system that includes `gcc` or
  another C compiler. On Debian and derivatives install the `build-essential`
  package. On Windows you can use [tdm-gcc](https://jmeubank.github.io/tdm-gcc/).
  Make sure that `gcc` is in your `PATH`.
- The tools listed in the [tools] section. It is advisable to first install the [cdo] tool, which can be used to easily perform the tasks described here. The [cdo] tool can most conveniently be installed using the [eget] tool.

```bash
eget szkiba/cdo
```

The [cdo] tool can then be used to perform the tasks described in the following sections.

Help about tasks:

```
cdo
```

[cdo]: (https://github.com/szkiba/cdo)
[eget]: https://github.com/zyedidia/eget

## tools - Install the required tools

Contributing will require the use of some tools, which can be installed most easily with a well-configured [eget] tool.

```bash
eget szkiba/mdcode
eget -t 1.57.2 golangci/golangci-lint
go install go.k6.io/xk6/cmd/xk6@latest
```

### lint - Run the linter

The `golangci-lint` tool is used for static analysis of the source code.
It is advisable to run it before committing the changes.

```bash
golangci-lint run
```

### test - Run the tests

```bash
go test -count 1 -race -coverprofile=coverage.out ./...
```

[test]: <#test---run-the-tests>

### coverage - View the test coverage report

Requires
: [test]

```bash
go tool cover -html=coverage.out
```

### build - Build a custom k6 with the extension

```bash
CGO_ENABLED=1 xk6 build --with github.com/grafana/xk6-sql=.
```

[build]: <#build---build-a-custom-k6-with-the-extension>

### clean - Delete the build directory

```bash
rm -rf build
rm -f ./k6
rm -f ./intg_test.db
```

### format - Applies Go formatting to code

```bash
go fmt ./...
```

### all - Clean build

Requires
: clean, format, test, build

### makefile - Generate Makefile

```bash
cdo --makefile Makefile
```
