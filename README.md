# golicense

golicense report information about the licenses of a Go binary or module and its dependencies.
Additionally can check the detected licenses against an allowed list.
Default is to look for a go.mod file into the current directory.

Licenses are detected using the
[google/licensecheck](https://github.com/google/licensecheck) package that scans
source texts for known licenses into the [module
cache](https://go.dev/ref/mod#module-cache).

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
  -a, --allowed          list of allowed licenses separated by comma (i.e. MIT, BSD-3-Clause). Default to all
  -d, --download         download dependencies to local cache
  -h, --help             show this help message
      --list-names       list the names of the license file can be detected and exit
      --list-licenses    list the licenses can be detected and exit
  -o, --output <file>    write to file instead of stdout
  -v, --verbose          make the tool verbose
      --version          show the version number
```

### License check for a Go module

```
golicense -a "MIT,BSD-3-CLAUSE" /path/to/go.mod > /dev/null
```

### Dump all licenses to a file

```
golicense -o LICENSE-THIRD-PARTY
```
