package buildx

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/urfave/cli/v3"
)

const testLicense = "MIT"

func TestMissingRequiredMetadata(t *testing.T) {
	t.Parallel()

	args := metadataArgs("Example application")
	for index := 0; index < len(args); index += 2 {
		incomplete := append([]string{}, args[:index]...)

		incomplete = append(incomplete, args[index+2:]...)
		if _, err := parseMetadata(t, incomplete); err == nil {
			t.Errorf("expected missing %s to be rejected", args[index])
		}
	}
}

func TestMetadataValuesPreserved(t *testing.T) {
	t.Parallel()

	for _, value := range []string{
		"Buildx smoke test, native platform",
		`Alice, Bob, Carol`,
		`A "quoted" value`,
		`C:\images\test`,
		"line one\nline two",
		"Unicode: Україна 🐳",
	} {
		t.Run(value, func(t *testing.T) {
			t.Parallel()

			labels, err := parseMetadata(t, append(metadataArgs(value), "--labels", "custom=preserved"))
			if err != nil {
				t.Fatal(err)
			}

			want := []string{
				"custom=preserved",
				"org.opencontainers.image.title=Example",
				"org.opencontainers.image.description=" + value,
				"org.opencontainers.image.authors=RevoTale",
				"org.opencontainers.image.licenses=MIT",
			}
			if !slices.Equal(labels, want) {
				t.Fatalf("labels = %q, want %q", labels, want)
			}
		})
	}
}

func TestEmptyMetadataRejected(t *testing.T) {
	t.Parallel()

	for _, value := range []string{"", " ", "\t\n"} {
		args := metadataArgs("Example application")
		for index := 1; index < len(args); index += 2 {
			invalid := slices.Clone(args)

			invalid[index] = value
			if _, err := parseMetadata(t, invalid); err == nil {
				t.Errorf("expected %s=%q to be rejected", args[index-1], value)
			}
		}
	}
}

func metadataArgs(description string) []string {
	return []string{"--title", "Example", "--description", description, "--authors", "RevoTale", "--licenses", testLicense}
}

func parseMetadata(t *testing.T, args []string) ([]string, error) {
	t.Helper()

	adapter := newWrapper()
	// Isolate parsing from CI environment variables; keep the real upstream flags and Before hook.
	for _, flag := range adapter.upstream.App.Flags {
		switch field := flag.(type) {
		case *cli.StringFlag:
			field.Sources = cli.EnvVars()
		case *cli.StringSliceFlag:
			field.Sources = cli.EnvVars()
		}
	}

	adapter.upstream.App.Action = acceptMetadata

	err := adapter.upstream.App.Run(t.Context(), append([]string{"buildx"}, args...))
	if err != nil {
		return nil, fmt.Errorf("parse metadata: %w", err)
	}

	return adapter.upstream.Settings.Build.Labels, nil
}

func acceptMetadata(_ context.Context, _ *cli.Command) error { return nil }

func TestMetadataFromEnvironment(t *testing.T) {
	t.Setenv("PLUGIN_TITLE", "CI smoke")
	t.Setenv("PLUGIN_DESCRIPTION", "Buildx smoke test, native platform")
	t.Setenv("PLUGIN_AUTHORS", "Alice, Bob")
	t.Setenv("PLUGIN_LICENSES", testLicense)
	t.Setenv("PLUGIN_LABELS", "custom=preserved")

	adapter := newWrapper()

	adapter.upstream.App.Action = acceptMetadata
	if err := adapter.upstream.App.Run(t.Context(), []string{"buildx"}); err != nil {
		t.Fatal(err)
	}

	want := []string{
		"custom=preserved",
		"org.opencontainers.image.title=CI smoke",
		"org.opencontainers.image.description=Buildx smoke test, native platform",
		"org.opencontainers.image.authors=Alice, Bob",
		"org.opencontainers.image.licenses=MIT",
	}
	if !slices.Equal(adapter.upstream.Settings.Build.Labels, want) {
		t.Fatalf("labels = %q, want %q", adapter.upstream.Settings.Build.Labels, want)
	}
}
