package buildx

import (
	"context"
	"fmt"
	"testing"

	"github.com/urfave/cli/v3"
)

const testLicense = "MIT"

func TestOptionalMetadata(t *testing.T) {
	t.Parallel()

	if err := parseMetadata(false, nil); err != nil {
		t.Fatalf("optional metadata rejected: %v", err)
	}
}

func TestMissingRequiredMetadata(t *testing.T) {
	t.Parallel()

	args := []string{
		"--title", "Example",
		"--description", "Example application",
		"--licenses", testLicense,
	}

	if err := parseMetadata(true, args); err == nil {
		t.Fatal("expected missing authors to be rejected")
	}
}

func TestProvidedRequiredMetadata(t *testing.T) {
	t.Parallel()

	args := []string{
		"--title", "Example",
		"--description", "Example application",
		"--authors", "RevoTale",
		"--licenses", testLicense,
	}

	if err := parseMetadata(true, args); err != nil {
		t.Fatalf("provided metadata rejected: %v", err)
	}
}

func TestImageLabelsPreserveExistingLabels(t *testing.T) {
	t.Parallel()

	fields := []metadataField{{name: "licenses", value: testLicense}}
	got := imageLabels(fields, "custom=value")
	want := "custom=value,org.opencontainers.image.licenses=MIT"

	if got != want {
		t.Fatalf("labels = %q, want %q", got, want)
	}
}

func TestImageLabelsEscapeCommas(t *testing.T) {
	t.Parallel()

	fields := []metadataField{{name: "authors", value: "Alice, Bob"}}
	got := imageLabels(fields, "")
	want := `org.opencontainers.image.authors=Alice\, Bob`

	if got != want {
		t.Fatalf("labels = %q, want %q", got, want)
	}
}

func TestImageLabelsSkipEmptyOptionalMetadata(t *testing.T) {
	t.Parallel()

	fields := []metadataField{{name: "authors"}}

	if got := imageLabels(fields, "custom=value"); got != "custom=value" {
		t.Fatalf("unexpected labels: %q", got)
	}
}

func parseMetadata(required bool, args []string) error {
	p := newWrapper()
	flags := p.metadataFlags(required)

	// Test CLI inputs independently of the parent process's environment.
	for _, flag := range flags {
		if field, ok := flag.(*cli.StringFlag); ok {
			field.Sources = cli.EnvVars()
		}
	}

	command := cli.Command{
		Name:   "buildx",
		Flags:  flags,
		Action: acceptMetadata,
	}

	if err := command.Run(context.Background(), append([]string{"buildx"}, args...)); err != nil {
		return fmt.Errorf("parse metadata: %w", err)
	}

	return nil
}

func acceptMetadata(_ context.Context, _ *cli.Command) error {
	return nil
}
