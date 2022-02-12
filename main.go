package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"text/tabwriter"
)

// Version can be used to set the version at link time
var Version string

type options struct {
	allowed          string
	download         bool
	dump             bool
	listLicenses     bool
	listLicenseNames bool
	output           string
	version          bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: lian [OPTIONS] [PATH]
lian is a license analyzer for Go binaries and modules. 
Default is to search for a go.mod file into the current directory.

Options:
  -a, --allowed          comma separated list of allowed licenses (i.e. MIT, BSD-3-Clause). Default to all
  -d, --download         download dependencies to local cache
      --dump             dump all licenses
  -h, --help             show this help message
      --list-names       list the names of the license file can be detected and exit
      --list-licenses    list the licenses can be detected and exit
  -o, --output <file>    write to file instead of stdout
  	  --version          show the version number
`)
	}

	var opts options

	flag.StringVar(&opts.allowed, "a", "", "")
	flag.StringVar(&opts.allowed, "allowed", "", "")
	flag.BoolVar(&opts.download, "d", false, "")
	flag.BoolVar(&opts.download, "download", false, "")
	flag.BoolVar(&opts.dump, "dump", false, "")
	flag.StringVar(&opts.output, "o", "", "")
	flag.StringVar(&opts.output, "output", "", "")
	flag.BoolVar(&opts.version, "version", false, "")
	flag.BoolVar(&opts.listLicenseNames, "list-names", false, "")
	flag.BoolVar(&opts.listLicenses, "list-licenses", false, "")
	flag.Parse()

	if opts.version {
		fmt.Println("lian", version())
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

	gomodcache := getGoModCache()
	licenses, err := getLicenses(gomodcache, mi, licenseNames)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var w io.Writer
	w = os.Stdout
	if opts.output != "" {
		f, err := os.Create(opts.output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
		w = f
	}

	if opts.dump {
		dump(w, licenses)
		os.Exit(0)
	}

	err = report(w, licenses, opts)
	if err != nil {
		os.Exit(1)
	}
}

func dump(w io.Writer, licenses []license) {
	// add the Go Programming Language license
	fmt.Fprintln(w, golangLicense)
	for _, l := range licenses {
		fmt.Fprintf(w, "* %s - (https://%s)\n", l.Version.String(), l.Version.Path)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "%s", l.Content)
		fmt.Fprintln(w)
	}
}

func report(w io.Writer, licenses []license, opts options) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.Debug)
	defer tw.Flush()

	// the table header
	th := "License\tDependency\tFile\tpkg.go.dev URL"

	var allowed []string
	if opts.allowed != "" {
		for _, v := range strings.Split(opts.allowed, ",") {
			v := strings.TrimSpace(v)
			if v != "" {
				allowed = append(allowed, v)
			}
		}
	}
	if len(allowed) > 0 {
		th = "Allowed\t" + th
	}

	// print the table header
	fmt.Fprintln(tw, th)

	var err error
	for _, l := range licenses {
		row := fmt.Sprintf("%s\t%s\t%s\thttps://pkg.go.dev/%s?tab=licenses", l.Type, l.Version.String(), l.Name, l.Version.String())
		if len(allowed) == 0 {
			// no allowed rule, print the data and continue
			fmt.Fprintln(tw, row)
			continue
		}

		// check for allowed license and add result to the report row
		isAllowedColumn := "Yes\t"
		if !isAllowedLicense(l, allowed) {
			err = fmt.Errorf("license not allowed: license %q - dependency %q", l.Type, l.Path)
			isAllowedColumn = "No\t"
		}
		row = isAllowedColumn + row
		fmt.Fprintln(tw, row)
	}

	return err
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
