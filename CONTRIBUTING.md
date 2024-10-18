# Contributing Guidelines

Thank you for your interest in contributing to **xk6-sql**!

Before you begin, make sure to familiarize yourself with the [Code of Conduct](CODE_OF_CONDUCT.md). If you've previously contributed to other open source project, you may recognize it as the classic [Contributor Covenant](https://contributor-covenant.org/).

If you want to chat with the team or the community, you can [join our community forums](https://community.grafana.com/c/grafana-k6).

### Filing issues

Don't be afraid to file issues! Nobody can fix a bug we don't know exists, or add a feature we didn't think of.

The worst that can happen is that someone closes it and points you in the right direction.

That said, "how do I..."-type questions are often more suited for [community forums](https://community.grafana.com/c/grafana-k6).

### Contributing code

If you'd like to contribute code, this is the basic procedure.

1. Find an issue you'd like to fix. If there is none already, or you'd like to add a feature, please open one, and we can talk about how to do it. Out of respect for your time, please start a discussion regarding any bigger contributions either in a GitHub Issue, in the community forums **before** you get started on the implementation.
  
   Remember, there's more to software development than code; if it's not properly planned, stuff gets messy real fast.

2. Create a fork and open a feature branch - `feature/my-cool-feature` is the classic way to name these, but it really doesn't matter.

3. Create a pull request!

4. We will discuss implementation details until everyone is happy, then a maintainer will merge it.

## Prerequisites

Prerequisites are listed in the [tools] section in addition to the [go toolchain](https://go101.org/article/go-toolchain.html) and [git](https://git-scm.com/) CLI.

### Special cases

If the database driver uses shared library, a build toolchain for your system that includes `gcc` or another C compiler. On Debian and derivatives install the `build-essential` package. On Windows you can use [tdm-gcc](https://jmeubank.github.io/tdm-gcc/). Make sure that `gcc` is in your `PATH`.

The `Makefile` is generated from the task list defined in the `CONTRIBUTING.md` file using the [cdo] tool. If the contribution is made to the task list, the `Makefile` must be regenerated, which is why the [cdo] tool is needed. The [cdo] tool can most conveniently be installed using the [eget] tool.

```bash
eget szkiba/cdo
```

[cdo]: https://github.com/szkiba/cdo
[eget]: https://github.com/zyedidia/eget

## Tasks

The tasks defined here can be executed manually or conveniently using the make or [cdo] tool.

**Help about tasks**

The command below lists the possible tasks.

using make:

```bash
make
```

using [cdo]:

```bash
cdo
```

**Execute task**

Tasks are executed by passing the name of the task as a parameter.

using make:

```bash
make taskname
```

using [cdo]:

```bash
cdo taskname
```

### tools - Install the required tools

Contributing will require the use of some tools, which can be installed most easily with a well-configured [eget] tool.

```bash
eget szkiba/mdcode
eget -t 1.57.2 golangci/golangci-lint
go install go.k6.io/xk6/cmd/xk6@latest
```

[tools]: (#tools---install-the-required-tools)
[xk6]: https://github.com/grafana/xk6

### lint - Run the linter

The `golangci-lint` tool is used for static analysis of the source code.
It is advisable to run it before committing the changes.

```bash
golangci-lint run
```

### test - Run the tests

The `go test` command is used to run the tests and generate the coverage report.

```bash
go test -count 1 -race -coverprofile=coverage.out -timeout 60s ./...
```

[test]: <#test---run-the-tests>

### coverage - View the test coverage report

The go `cover` tool should be used to display the coverage report in the browser.

Requires
: [test]

```bash
go tool cover -html=coverage.out
```

### build - Build custom k6 with extension

The [xk6] tool is used to build the k6. In addition to the xk6-sql extension, the database driver extension you want to use should also be specified.

```bash
xk6 build --with github.com/grafana/xk6-sql=. --with github.com/grafana/xk6-sql-driver-ramsql
```

[build]: <#build---build-custom-k6-with-extension>

### example - Run the example

Run the example embedded in `README.md` and save the output.

```bash
./k6 run examples/example.js  > examples/example.txt 2>&1
```

### readme - Update README.md

Update the example code and its output in `README.md`.

```bash
mdcode update
```

### clean - Clean the working directory

Delete the work files created in the work directory (also included in .gitignore).

```bash
rm -f ./k6
```

[clean]: (#clean---clean-the-working-directory)

### all - Clean build

Performs the most important tasks. It can be used to check whether the CI workflow will run successfully.

Requires
: [clean], [format], [test], [build]

### format - Format the go source codes

```bash
go fmt ./...
```

[format]: (#format---format-the-go-source-codes)

### makefile - Generate the Makefile

```bash
cdo --makefile Makefile
```
