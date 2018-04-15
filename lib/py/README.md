
# Skycoin Python client library

[![Build Status](https://travis-ci.org/skycoin/skycoin.svg)](https://travis-ci.org/skycoin/skycoin)
[![GoDoc](https://godoc.org/github.com/skycoin/skycoin?status.svg)](https://godoc.org/github.com/skycoin/skycoin)
[![Go Report Card](https://goreportcard.com/badge/github.com/skycoin/skycoin)](https://goreportcard.com/report/github.com/skycoin/skycoin)

Skycoin Python client library (a.k.a pyskycoin) provides access to Skycoin Core
internal and API functions for implementing third-party applications.

## API Interface

The interface is defined in the [pyskycoin API docs](docs/libpy/).

## Building

```sh
$ make build-libpy
```

## Testing

In order to test the C client libraries follow these steps

- Install [py.test]()
  * locally by executing `make install-deps-libpy` command
  * or by [installing py.test system-wide]()
- Run `make test-libpy` command

## Binary distribution

The following files will be generated

