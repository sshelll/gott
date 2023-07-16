# Gott

<a href="https://996.icu"><img src="https://img.shields.io/badge/link-996.icu-red.svg" alt="996.icu" /></a>
[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)

Go test tool with terminal UI.

## Demo

<img src="/gif/demo.gif" alt="demo" width=50%>

## Install

go version >= 1.17, clone this repo and exec `go build .` should be ok，or you can exec the command below：

```sh
go install github.com/sshelll/gott@latest
```

go version <= 1.16, go to Release page and download the executable file(MacOS Only, I'm lazy...)

## Usage

**Use `gott -h` to get more details**

**Use `gott` instead of `go test`**

Or you can exec `gott -p` to get the test name, in this way you won't run `go test`

For example:

`go test -v` => `gott -v`

`go test -v -race` => `gott -v -race`

`go test -gcflags="all=-l -N"` => `gott -gcflags=\"all=-l -N\"`

`go test -gcflags=all=-l -race -coverprofile=coverage.out` => `gott -gcflags=all=-l -race -coverprofile=coverage.out`

`gott -p` => `you will get a go test func name`

**Use the script below to debug a test with dlv:**

```sh
#!/bin/zsh

fn=$(gott -p)

if [ ! $fn ]; then exit 0; fi

dlv test --build-flags=-test.run $fn
```

## QA

- Q1：Does this program recognize `github.com/stretchr/testify/suite` ？

  A1：Yes，but the entry func of 'suite' is limited, you can see how it works in the examples below:

  ```go
  type FooTestSuite struct {
    suite.Suite
  }

  // OK, allow 'new'
  func TestFoo1(t *testing.T) {
    suite.Run(t, new(FooTestSuite))
  }

  // OK, allow '&'
  func TestFoo2(t *testing.T) {
    suite.Run(t, &FooTestSuite{})
  }

  // not OK
  func TestFoo3(t *testing.T) {
    foo := new(FooTestSuite)
    suite.Run(t, foo)
  }

  // not OK
  func TestFoo4(t *testing.T) {
    m := make(map[int]interface{})
    m[1] = &FooTestSuite{}
    suite.Run(t, m[1])
  }
  ```

- Q2：What are the key mappings?

  A2：`↑` `↓` to move cursor，`/` to search (just strings.Contains()），`esc` to quit。
