package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/module"
)

func TestGetModuleInfoWithReplace(t *testing.T) {
	// Create a temporary go.mod file with replace directive
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")

	goModContent := `module test/replace

go 1.18

require (
    github.com/example/original v1.0.0
    github.com/example/normal v1.0.0
)

replace github.com/example/original => github.com/example/replacement v2.0.0
`

	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	require.NoError(t, err, "Failed to create test go.mod")

	// Test getModuleInfo
	mi, err := getModuleInfo(goModPath)
	require.NoError(t, err, "getModuleInfo should not fail")

	// Verify module info
	expectedModule := module.Version{Path: "test/replace"}
	assert.Equal(t, expectedModule, mi.Module, "Module should match expected")

	// Verify that replace directive is applied
	require.Len(t, mi.Require, 2, "Should have exactly 2 requirements")

	// Build a map for easier verification
	reqMap := make(map[string]string)
	for _, req := range mi.Require {
		reqMap[req.Path] = req.Version
	}

	// Verify the replacement happened
	assert.Equal(t, "v2.0.0", reqMap["github.com/example/replacement"], "Replacement module should have correct version")
	assert.Equal(t, "v1.0.0", reqMap["github.com/example/normal"], "Normal module should have correct version")
	assert.NotContains(t, reqMap, "github.com/example/original", "Original module should be replaced, not present")
}

func TestGetModuleInfoWithoutReplace(t *testing.T) {
	// Create a temporary go.mod file without replace directive
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")

	goModContent := `module test/normal

go 1.18

require (
    github.com/example/dep1 v1.0.0
    github.com/example/dep2 v1.2.0
)
`

	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	require.NoError(t, err, "Failed to create test go.mod")

	// Test getModuleInfo
	mi, err := getModuleInfo(goModPath)
	require.NoError(t, err, "getModuleInfo should not fail")

	// Verify module info
	expectedModule := module.Version{Path: "test/normal"}
	assert.Equal(t, expectedModule, mi.Module, "Module should match expected")

	// Verify requirements
	require.Len(t, mi.Require, 2, "Should have exactly 2 requirements")

	// Build a map for easier verification
	reqMap := make(map[string]string)
	for _, req := range mi.Require {
		reqMap[req.Path] = req.Version
	}

	expectedReqs := map[string]string{
		"github.com/example/dep1": "v1.0.0",
		"github.com/example/dep2": "v1.2.0",
	}

	assert.Equal(t, expectedReqs, reqMap, "Requirements should match expected")
}

func TestExclusionFormat(t *testing.T) {
	// Test how ModulePath() formats the path for exclusion comparison
	l := license{
		Version: module.Version{
			Path:    "github.com/example/module",
			Version: "v1.0.0",
		},
	}

	modpath, err := l.ModulePath()
	assert.NoError(t, err)

	t.Logf("ModulePath() returns: %q", modpath)
	t.Logf("Expected exclusion format: %q", "github.com/example/module@v1.0.0")

	// Test if they match
	assert.Equal(t, "github.com/example/module@v1.0.0", modpath)
}

func TestExclusionWithSpecialCharacters(t *testing.T) {
	// Test module with special characters that might need escaping
	l := license{
		Version: module.Version{
			Path:    "github.com/example/module-with-dash",
			Version: "v1.0.0",
		},
	}

	modpath, err := l.ModulePath()
	assert.NoError(t, err)

	t.Logf("ModulePath() with dash: %q", modpath)
}

func TestExclusionWithReplace(t *testing.T) {
	// Test the actual exclusion logic with replace directives
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")

	goModContent := `module test/exclusion

go 1.18

require (
    github.com/example/original v1.0.0
    github.com/example/normal v1.0.0
)

replace github.com/example/original => github.com/example/replacement v2.0.0
`

	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	require.NoError(t, err)

	// Get module info with replacement
	mi, err := getModuleInfo(goModPath)
	require.NoError(t, err)

	// Test exclusion of the replacement module
	excluded := map[string]struct{}{
		"github.com/example/replacement@v2.0.0": {},
	}

	// Check if exclusion would work for the replacement
	for _, v := range mi.Require {
		l := license{Version: v}
		modpath, err := l.ModulePath()
		require.NoError(t, err)

		t.Logf("Module: %s, ModulePath: %s", v.String(), modpath)

		if v.Path == "github.com/example/replacement" {
			_, isExcluded := excluded[modpath]
			assert.True(t, isExcluded, "Replacement module should be excluded")
		}
	}
}

func TestEndToEndExclusion(t *testing.T) {
	// Test exclusion parsing from command line format
	excludedFlag := "github.com/example/replacement@v2.0.0,github.com/example/other@v1.0.0"

	// Parse exclusions like main.go does
	excluded := splitCommaSeparatedFlag(excludedFlag)
	excludedExist := map[string]struct{}{}
	for _, e := range excluded {
		excludedExist[e] = struct{}{}
	}

	// Test if our module would be excluded
	l := license{
		Version: module.Version{
			Path:    "github.com/example/replacement",
			Version: "v2.0.0",
		},
	}

	modpath, err := l.ModulePath()
	require.NoError(t, err)

	t.Logf("Excluded list: %v", excluded)
	t.Logf("ModulePath: %s", modpath)

	_, isExcluded := excludedExist[modpath]
	assert.True(t, isExcluded, "Module should be excluded")
}
