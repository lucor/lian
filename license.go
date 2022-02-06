package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/licensecheck"
	"golang.org/x/mod/module"
)

type license struct {
	// Name holds the license file name (i.e. LICENSE.md)
	Name string
	// ModuleInfo holds the module info (module + version)
	ModuleInfo string
	// Path holds the license path on the local host (GOMODCACHE/ModuleInfo)
	Path string
	// Type holds the license type (i.e. MIT)
	Type string
	// Content holds the license content
	Content []byte
}

var licenseNames = []string{
	"COPYING",
	"COPYING.MD",
	"COPYING.MD",
	"LICENSE",
	"LICENSE.MD",
	"LICENSE.TXT",
}

func getGoModCache() string {
	gomodcache := os.Getenv("GOMODCACHE")
	if gomodcache != "" {
		return gomodcache
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return filepath.Join(gopath, "pkg", "mod")
}

func getLicenses(gomodcache string, info moduleInfo, patterns []string) ([]license, error) {
	licenses := []license{}
	for _, v := range info.Require {
		epath, err := module.EscapePath(v.Path)
		if err != nil {
			return nil, fmt.Errorf("invalid module path: %s", err)
		}

		modpath := filepath.Join(gomodcache, epath+"@"+v.Version)

		licenseFiles, err := findLicenses(modpath, patterns)
		if err != nil {
			return nil, fmt.Errorf("could not scan for licenses: %w", err)
		}

		if len(licenseFiles) == 0 {
			return nil, fmt.Errorf("license not found for %s in %s", v.String(), modpath)
		}

		for _, licenseFile := range licenseFiles {
			licenseFilePath := filepath.Join(modpath, licenseFile)
			data, err := os.ReadFile(licenseFilePath)
			if err != nil {
				return nil, fmt.Errorf("error reading license file: %w", err)
			}

			coverage := licensecheck.Scan(data)
			for _, match := range coverage.Match {
				l := license{
					Name:       licenseFile,
					ModuleInfo: v.String(),
					Path:       licenseFilePath,
					Type:       match.ID,
					Content:    data,
				}
				licenses = append(licenses, l)
			}
		}
	}
	return licenses, nil
}

func findLicenses(path string, licenseFileNames []string) ([]string, error) {
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
		for _, v := range licenseFileNames {
			if upperFileName == v {
				licenses = append(licenses, entry.Name())
				break
			}
		}
	}
	return licenses, nil
}

func listLicenses() {
	m := map[string]bool{}
	licenses := []string{}
	for _, v := range licensecheck.BuiltinLicenses() {
		l := v.ID
		if _, ok := m[l]; ok {
			continue
		}
		m[l] = true
		licenses = append(licenses, l)
	}
	sort.Strings(licenses)
	for _, v := range licenses {
		fmt.Println(v)
	}
}

func listLicenseNames() {
	for _, v := range licenseNames {
		fmt.Println(v)
	}
}
