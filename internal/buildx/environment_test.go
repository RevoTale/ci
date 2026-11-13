package buildx

import (
	"errors"
	"os"
	"path/filepath"
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
