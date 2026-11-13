// Package buildx adds required OCI metadata to the official Woodpecker Buildx plugin.
package buildx

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	upstream "codeberg.org/woodpecker-plugins/plugin-docker-buildx/plugin"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

var errEmptyMetadata = errors.New("metadata must not be empty")
var errEnvironment = errors.New("cannot load plugin environment file")

type metadataField struct {
	name  string
	help  string
	value string
}

type wrapper struct {
	upstream *upstream.Plugin
	metadata []metadataField
}

// Run validates OCI metadata and runs the official Buildx implementation.
func Run() {
	if err := loadPluginEnvironment(); err != nil {
		log.Fatal(err)
	}

	newWrapper().upstream.Run()
}

func newWrapper() *wrapper {
	adapter := &wrapper{metadata: []metadataField{
		{name: "title", help: "Name of the image"},
		{name: "description", help: "What the image runs or provides"},
		{name: "authors", help: "Author names for this image"},
		{name: "licenses", help: "License for this image, for example MIT"},
	}}
	adapter.upstream = upstream.New(nil, "6.1.2")
	adapter.upstream.App.Flags = append(adapter.upstream.App.Flags, adapter.metadataFlags()...)
	adapter.upstream.App.Before = adapter.applyMetadata

	return adapter
}

func (adapter *wrapper) metadataFlags() []cli.Flag {
	flags := make([]cli.Flag, 0, len(adapter.metadata))
	for i := range adapter.metadata {
		field := &adapter.metadata[i]
		flags = append(flags, &cli.StringFlag{
			Name:        field.name,
			Usage:       field.help + " (Woodpecker: settings." + field.name + ")",
			Required:    true,
			Sources:     cli.EnvVars("PLUGIN_" + strings.ToUpper(field.name)),
			Destination: &field.value,
		})
	}

	return flags
}

func (adapter *wrapper) applyMetadata(ctx context.Context, _ *cli.Command) (context.Context, error) {
	for _, field := range adapter.metadata {
		if strings.TrimSpace(field.value) == "" {
			return ctx, fmt.Errorf("%w: %s", errEmptyMetadata, field.name)
		}

		adapter.upstream.Settings.Build.Labels = append(adapter.upstream.Settings.Build.Labels,
			"org.opencontainers.image."+field.name+"="+field.value)
	}

	return ctx, nil
}

// Preserve the official command's environment-file support before parsing flags.
func loadPluginEnvironment() error {
	files := []string{}
	if _, err := os.Stat("/run/drone/env"); err == nil {
		files = append(files, "/run/drone/env")
	}

	if file := os.Getenv("PLUGIN_ENV_FILE"); file != "" {
		files = append(files, file)
	}

	for _, file := range files {
		if err := godotenv.Overload(file); err != nil {
			return fmt.Errorf("%w: %s", errEnvironment, file)
		}
	}

	return nil
}
