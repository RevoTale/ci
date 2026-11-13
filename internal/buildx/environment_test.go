package buildx

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPluginEnvironmentFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "plugin.env")
	t.Setenv("PLUGIN_ENV_FILE", file)
	t.Setenv("PLUGIN_TITLE", "original")

	if err := os.WriteFile(file, []byte("PLUGIN_TITLE='From file, intact'\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := loadPluginEnvironment(); err != nil {
		t.Fatal(err)
	}

	if got := os.Getenv("PLUGIN_TITLE"); got != "From file, intact" {
		t.Fatalf("title = %q", got)
	}
}

func TestMissingPluginEnvironmentFileRejected(t *testing.T) {
	t.Setenv("PLUGIN_ENV_FILE", filepath.Join(t.TempDir(), "missing.env"))

	if err := loadPluginEnvironment(); !errors.Is(err, errEnvironment) {
		t.Fatalf("error = %v, want environment-file error", err)
	}
}

func TestInvalidPluginEnvironmentFileRejected(t *testing.T) {
	file := filepath.Join(t.TempDir(), "plugin.env")
	t.Setenv("PLUGIN_ENV_FILE", file)
	t.Setenv("CI_METADATA_TEST_VALUE", "unchanged")

	if err := os.WriteFile(file, []byte("CI_METADATA_TEST_VALUE='unterminated"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := loadPluginEnvironment()
	if !errors.Is(err, errEnvironment) {
		t.Fatalf("error = %v, want environment-file error", err)
	}

	if strings.Contains(err.Error(), "unterminated") {
		t.Fatal("environment-file error exposes its contents")
	}

	if os.Getenv("CI_METADATA_TEST_VALUE") != "unchanged" {
		t.Fatal("invalid environment file changed the environment")
	}
}
