package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

// Version can be used to set the version at link time
var Version string

type options struct {
	allowed          string
	download         bool
	listLicenses     bool
	listLicenseNames bool
	output           string
	verbose          bool
	version          bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: golicense [OPTIONS] [PATH]
List information about the licenses of a Go module or binary and its dependencies.
Additionally can check the detected licenses against an allowed list.
Default is to look for a go.mod file into the current directory.

Options:
  -a, --allowed          comma separated list of allowed licenses (i.e. MIT, BSD-3-Clause). Default to all
  -d, --download         download dependencies to local cache
  -h, --help             show this help message
      --list-names       list the names of the license file can be detected and exit
      --list-licenses    list the licenses can be detected and exit
  -o, --output <file>    write to file instead of stdout
  -v, --verbose          print additional info on stderr
  	  --version          show the version number
`)
	}

	var opts options

	flag.StringVar(&opts.allowed, "a", "", "")
	flag.StringVar(&opts.allowed, "allowed", "", "")
	flag.BoolVar(&opts.download, "d", false, "")
	flag.BoolVar(&opts.download, "download", false, "")
	flag.StringVar(&opts.output, "o", "", "")
	flag.StringVar(&opts.output, "output", "", "")
	flag.BoolVar(&opts.verbose, "v", false, "")
	flag.BoolVar(&opts.verbose, "verbose", false, "")
	flag.BoolVar(&opts.version, "version", false, "")
	flag.BoolVar(&opts.listLicenseNames, "list-names", false, "")
	flag.BoolVar(&opts.listLicenses, "list-licenses", false, "")
	flag.Parse()

	if opts.version {
		fmt.Println("golicense", version())
		os.Exit(0)
	}

	if opts.listLicenses {
		listLicenses()
		os.Exit(0)
	}

	if opts.listLicenseNames {
		listLicenseNames()
		os.Exit(0)
	}

	var allowed []string
	if opts.allowed != "" {
		for _, v := range strings.Split(opts.allowed, ",") {
			v := strings.TrimSpace(v)
			if v != "" {
				allowed = append(allowed, v)
			}
		}
	}

	path := "go.mod"
	if len(flag.Args()) > 0 {
		path = flag.Arg(0)
	}
	mi, err := getModuleInfo(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if opts.download {
		err = downloadModules(mi)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	licenses, err := getLicenses(getGoModCache(), mi, licenseNames)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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

	for _, l := range licenses {
		fmt.Fprintf(w, "* %s - (https://%s)\n", l.Version.String(), l.Version.Path)

		if !isAllowedLicense(l, allowed) {
			fmt.Fprintf(os.Stderr, "[✗] Not allowed license for %q. Found %q, want %q\n", l.Version.Path, l.Type, allowed)
			os.Exit(1)
		}

		if opts.verbose {
			fmt.Fprintf(os.Stderr, "Found: %s\n", l.Path)
			fmt.Fprintf(os.Stderr, "License: %s\n", l.Type)
			fmt.Fprintf(os.Stderr, "Source: https://%s\n", l.Version.Path)
			fmt.Fprintf(os.Stderr, "Link: https://pkg.go.dev/%s?tab=licenses\n", l.Version.String())
		}
		fmt.Fprintln(w)
		fmt.Fprintf(w, "%s", l.Content)
		fmt.Fprintln(w)
	}
}

func isAllowedLicense(l license, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, v := range allowed {
		if v == l.Type {
			return true
		}
	}
	return false
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
