package main

import (
	"debug/buildinfo"
	"errors"
	"fmt"
	"os"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

type moduleInfo struct {
	Module  module.Version
	Require []module.Version
}

// getModuleInfo extracts the module info from a go.mod or Go binary file
func getModuleInfo(path string) (moduleInfo, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return moduleInfo{}, fmt.Errorf("%s: no such go.mod file or Go binary", path)
		}
		return moduleInfo{}, err
	}

	if fi.IsDir() {
		return moduleInfo{}, fmt.Errorf("path must be a go.mod file or a Go binary built with module enabled")
	}

	if fi.Name() == "go.mod" {
		data, err := os.ReadFile(path)
		if err != nil {
			return moduleInfo{}, fmt.Errorf("could not read the go.mod file: %w", err)
		}

		// parse the go.mod content
		mf, err := modfile.Parse("", data, nil)
		if err != nil {
			return moduleInfo{}, fmt.Errorf("could not parse the go.mod file: %w", err)
		}

		mi := moduleInfo{
			Module: mf.Module.Mod,
		}

		for _, v := range mf.Require {
			mi.Require = append(mi.Require, v.Mod)
		}
		return mi, nil
	}

	bi, err := buildinfo.ReadFile(path)
	if err != nil {
		return moduleInfo{}, err
	}

	mi := moduleInfo{
		Module: module.Version{
			Path:    bi.Main.Path,
			Version: bi.Main.Version,
		},
	}

	for _, dep := range bi.Deps {
		m := module.Version{
			Path:    dep.Path,
			Version: dep.Version,
		}
		mi.Require = append(mi.Require, m)
	}
	return mi, nil
}
