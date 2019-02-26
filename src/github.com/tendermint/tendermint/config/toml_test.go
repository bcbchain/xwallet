package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ensureFiles(t *testing.T, rootDir string, files ...string) {
	for _, f := range files {
		p := rootify(rootDir, f)
		_, err := os.Stat(p)
		assert.Nil(t, err, p)
	}
}

func TestEnsureRoot(t *testing.T) {
	require := require.New(t)

	tmpDir, err := ioutil.TempDir("", "config-test")
	require.Nil(err)
	defer os.RemoveAll(tmpDir)

	EnsureRoot(tmpDir)

	data, err := ioutil.ReadFile(filepath.Join(tmpDir, defaultConfigFilePath))
	require.Nil(err)

	if !checkConfig(string(data)) {
		t.Fatalf("config file missing some information")
	}

	ensureFiles(t, tmpDir, "data")
}

func TestEnsureTestRoot(t *testing.T) {
	require := require.New(t)

	testName := "ensureTestRoot"

	cfg := ResetTestRoot(testName)
	rootDir := cfg.RootDir

	data, err := ioutil.ReadFile(filepath.Join(rootDir, defaultConfigFilePath))
	require.Nil(err)

	if !checkConfig(string(data)) {
		t.Fatalf("config file missing some information")
	}

	baseConfig := DefaultBaseConfig()
	ensureFiles(t, rootDir, defaultDataDir, baseConfig.Genesis, baseConfig.PrivValidator)
}

func checkConfig(configFile string) bool {
	var valid bool

	var elems = []string{
		"moniker",
		"seeds",
		"proxy_app",
		"fast_sync",
		"create_empty_blocks",
		"peer",
		"timeout",
		"broadcast",
		"send",
		"addr",
		"wal",
		"propose",
		"max",
		"genesis",
	}
	for _, e := range elems {
		if !strings.Contains(configFile, e) {
			valid = false
		} else {
			valid = true
		}
	}
	return valid
}
