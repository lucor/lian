# golicense

golicense reports information about the licenses of a Go module or binary and its dependencies.
## Use cases

golicense aims to help in the following use cases:

- report all the dependencies, their versions, and license type along with the URL on pkg.go.dev
- dump all licenses to comply with package distribution
- check against a set of allowed licenses

## Example

`golicense` in action with itself

<div align="center">
    <img alt="golicense example" src="example.gif" />
</div>

## How it works

It is designed to work without connecting to third-party services.

The licenses are detected using the
[google/licensecheck](https://github.com/google/licensecheck) library that will scans
source texts for known licenses directly from the [module cache](https://go.dev/ref/mod#module-cache).

The module cache usually is already warmed if the module has been already built locally.
If the dependencies are not present the `-d, --download` option can be specified and golicense will automatically download the dependencies using the `go mod download` command.

## Installation

```
$ go install github.com/lucor/golicense@latest
```

Note: requires Go >= 1.18

## Download

Pre-built binaries can be downloaded from the [releases](https://github.com/lucor/golicense/releases) page

## Usage

```
Usage: golicense [OPTIONS] [PATH]

Options:
  -a, --allowed          comma separated list of allowed licenses (i.e. MIT, BSD-3-Clause). Default to all
  -d, --download         download dependencies to local cache
      --dump             dump all licenses
  -h, --help             show this help message
      --list-names       list the names of the license file can be detected and exit
      --list-licenses    list the licenses can be detected and exit
  -o, --output <file>    write to file instead of stdout
  	  --version          show the version number
```

### License check for a Go module

```
golicense --allowed "MIT,BSD-3-CLAUSE" /path/to/go.mod
```

### Dump all licenses to a file for a Go binary

```
golicense --dump -o LICENSE-THIRD-PARTY /path/to/go_binary
```
