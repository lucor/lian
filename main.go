package main

import (
	"debug/buildinfo"
	"flag"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"golang.org/x/mod/module"
)

// Version can be used to set the version at link time
var Version string

type options struct {
	output  string
	verbose bool
	version bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: golicense [OPTIONS] GO_BINARY
Collect all the licenses of the dependencies of a GO_BINARY built with module support.

Options:
  -h, --help		show this help message
  -o, --output <file>	write to file instead of stdout
  -v, --verbose		make the tool verbose
  -V, --version		show the version number
`)
	}

	var opts options

	flag.StringVar(&opts.output, "o", "", "")
	flag.StringVar(&opts.output, "output", "", "")
	flag.BoolVar(&opts.verbose, "v", false, "")
	flag.BoolVar(&opts.verbose, "verbose", false, "")
	flag.BoolVar(&opts.version, "V", false, "")
	flag.BoolVar(&opts.version, "version", false, "")
	flag.Parse()

	if opts.version {
		fmt.Println("golicense", version())
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(0)
	}

	w := os.Stdout
	if opts.output != "" {
		f, err := os.Create(opts.output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
		w = f
	}

	gobinary := flag.Arg(0)
	info, err := buildinfo.ReadFile(gobinary)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	modRootPath := filepath.Join(gopath, "pkg", "mod")
	for _, v := range info.Deps {
		epath, err := module.EscapePath(v.Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid module path", err)
			continue
		}
		mod := v.Path + "@" + v.Version
		modpath := filepath.Join(modRootPath, epath+"@"+v.Version)
		if opts.verbose {
			fmt.Fprintf(os.Stderr, "dependency: %s\n", mod)
			fmt.Fprintf(os.Stderr, "local path: %s\n", modpath)
			fmt.Fprintf(os.Stderr, "pkg.go.dev link: https://pkg.go.dev/%s?tab=licenses\n", mod)
		}

		licenses, err := findLicenses(modpath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to look for licenses", err)
			continue
		}

		if len(licenses) == 0 {
			fmt.Fprintf(os.Stderr, "license not found at %s\n", modpath)
			continue
		}

		if opts.verbose {
			fmt.Fprintf(os.Stderr, "found: %v\n", licenses)
		}

		// write license to writer
		fmt.Fprintf(w, "## %s\n", mod)
		for _, license := range licenses {
			licenseFile := filepath.Join(modpath, license)
			b, err := os.ReadFile(licenseFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading license from %s. Error %s\n", licenseFile, err)
			}
			fmt.Fprintf(w, "Source: %s\n\n", filepath.Join(mod, license))
			fmt.Fprintf(w, "%s\n", b)
		}
	}
}

var licensePatterns = []string{
	"LICENSE",
	"COPYING",
}

func findLicenses(path string) ([]string, error) {
	licenses := []string{}
	entries, err := os.ReadDir(path)
	if err != nil {
		return licenses, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		upperFileName := strings.ToUpper(entry.Name())
		for _, v := range licensePatterns {
			if strings.HasPrefix(upperFileName, v) {
				licenses = append(licenses, entry.Name())
				break
			}
		}
	}
	return licenses, nil
}

func version() string {
	if Version != "" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}
	return "(unknown)"
}
