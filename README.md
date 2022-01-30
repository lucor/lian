# golicense

golicense is a simple tool that attempts to collect all the licenses of the dependencies of a Go binary in order to be used in redistribution.

golicense needs to be used with binaries compiled locally since it looks for licenses into the default location of the module cache (`$GOPATH/pkg/mod`).

golicense is not meant to be used for open source compliance

## Installation

```
$ go install github.com/lucor/golicense@latest
```

Note: requires Go >= 1.18

## Usage

```
Usage: golicense [OPTIONS] GO_BINARY

Collect all the licenses of the dependencies of a GO_BINARY built with module support.

Options:
  -h, --help		show this help message
  -o, --output <file>	write to file instead of stdout
  -v, --verbose		make the tool verbose
  -V, --version		show the version number
```
