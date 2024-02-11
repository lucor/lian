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
	// Content holds the license content
	Content []byte
	// Name holds the license file name (i.e. LICENSE.md)
	Name string
	// Type holds the license type (i.e. MIT)
	Type string

	module.Version
}

func (l license) ModuleVersion() string {
	return l.Version.Path + "@" + l.Version.Version
}

func (l license) ModulePath() (string, error) {
	epath, err := module.EscapePath(l.Version.Path)
	if err != nil {
		return "", fmt.Errorf("invalid module path: %s", err)
	}
	return epath + "@" + l.Version.Version, nil
}

func (l license) LicensePath() (string, error) {
	modpath, err := l.ModulePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(modpath, l.Name), nil
}

var licenseNames = []string{
	"COPYING",
	"COPYING.MD",
	"COPYING.TXT",
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

func getLicenses(gomodcache string, info moduleInfo, patterns []string, excluded map[string]struct{}) ([]license, error) {
	licenses := []license{}
	for _, v := range info.Require {

		l := license{
			Version: v,
		}

		modpath, err := l.ModulePath()
		if err != nil {
			return nil, err
		}

		if _, ok := excluded[modpath]; ok {
			continue
		}

		licenseFiles, err := findLicenses(filepath.Join(gomodcache, modpath), patterns)
		if err != nil {
			return nil, fmt.Errorf("could not scan for licenses: %w", err)
		}

		if len(licenseFiles) == 0 {
			return nil, fmt.Errorf("license not found for %s in %s", v.String(), modpath)
		}

		for _, name := range licenseFiles {
			l := l
			l.Name = name
			licensePath, err := l.LicensePath()
			if err != nil {
				return nil, err
			}

			data, err := os.ReadFile(filepath.Join(gomodcache, licensePath))
			if err != nil {
				return nil, fmt.Errorf("error reading license file: %w", err)
			}

			coverage := licensecheck.Scan(data)
			for _, match := range coverage.Match {
				l.Content = data
				l.Type = match.ID
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
