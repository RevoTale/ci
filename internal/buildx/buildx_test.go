package buildx

import (
	"context"
	"fmt"
	"testing"

	"github.com/urfave/cli/v3"
)

const testLicense = "MIT"

func TestMissingRequiredMetadata(t *testing.T) {
	t.Parallel()

	args := []string{
		"--title", "Example",
		"--description", "Example application",
		"--authors", "RevoTale",
		"--licenses", testLicense,
	}

	for index := 0; index < len(args); index += 2 {
		incomplete := append([]string{}, args[:index]...)
		incomplete = append(incomplete, args[index+2:]...)

		if err := parseMetadata(incomplete); err == nil {
			t.Errorf("expected missing %s to be rejected", args[index])
		}
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

	if err := parseMetadata(args); err != nil {
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

func TestImageLabelsSkipExplicitEmptyMetadata(t *testing.T) {
	t.Parallel()

	fields := []metadataField{{name: "authors"}}

	if got := imageLabels(fields, "custom=value"); got != "custom=value" {
		t.Fatalf("unexpected labels: %q", got)
	}
}

func parseMetadata(args []string) error {
	p := newWrapper()
	flags := p.metadataFlags()

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
