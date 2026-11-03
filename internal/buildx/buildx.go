// Package buildx adds OCI metadata settings to the official Woodpecker Buildx plugin.
package buildx

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	plugin "codeberg.org/woodpecker-plugins/go-plugin"
	"github.com/urfave/cli/v3"
)

type metadataField struct {
	name  string
	help  string
	value string
}

type wrapper struct {
	metadata []metadataField
}

// Run adds required OCI metadata settings to the official Buildx plugin.
func Run() {
	p := newWrapper()

	plugin.New(plugin.Options{
		Name:    "docker-buildx",
		Flags:   p.metadataFlags(),
		Execute: p.runBuildx,
	}).Run()
}

func newWrapper() *wrapper {
	return &wrapper{
		metadata: []metadataField{
			{name: "title", help: "Name of the image"},
			{name: "description", help: "What the image runs or provides"},
			{name: "authors", help: "Author names for this image"},
			{name: "licenses", help: "License for this image, for example MIT"},
		},
	}
}

func (p *wrapper) metadataFlags() []cli.Flag {
	flags := make([]cli.Flag, 0, len(p.metadata))

	for i := range p.metadata {
		field := &p.metadata[i]

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

func (p *wrapper) runBuildx(_ context.Context) error {
	labels := imageLabels(p.metadata, os.Getenv("PLUGIN_LABELS"))

	if err := os.Setenv("PLUGIN_LABELS", labels); err != nil {
		return fmt.Errorf("set image labels: %w", err)
	}

	return startOfficialBuildx()
}

func imageLabels(fields []metadataField, existing string) string {
	labels := existing

	for _, field := range fields {
		if field.value == "" {
			continue
		}

		if labels != "" {
			labels += ","
		}

		// Buildx separates labels with commas; escape commas inside values.
		value := strings.ReplaceAll(field.value, ",", `\,`)
		labels += "org.opencontainers.image." + field.name + "=" + value
	}

	return labels
}

func startOfficialBuildx() error {
	// Preserve the official plugin's Docker daemon startup.
	const entrypoint = "/usr/local/bin/dockerd-entrypoint.sh"

	// #nosec G204 -- Executable and argv are fixed; inherited environment carries plugin settings.
	if err := syscall.Exec(
		entrypoint,
		[]string{entrypoint, "plugin-docker-buildx"},
		os.Environ(),
	); err != nil {
		return fmt.Errorf("start official Buildx plugin: %w", err)
	}

	return nil
}
