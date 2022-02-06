# golicense

golicense list information about the licenses of a Go module or binary and its dependencies.

Licenses are detected using the
[google/licensecheck](https://github.com/google/licensecheck) package that scans
source texts for known licenses into the [module
cache](https://go.dev/ref/mod#module-cache).

## Installation

```
$ go install github.com/lucor/golicense@latest
```

Note: requires Go >= 1.18

## Usage

```
Usage: golicense [OPTIONS] [PATH]

Options:
  -h, --help             show this help message
  -i, --include          include the licenses in the output
      --list-names       list the names of the license file can be detected and exit
      --list-licenses    list the licenses can be detected and exit
  -o, --output <file>    write to file instead of stdout
  -v, --verbose          make the tool verbose
  -V, --version          show the version number
```
