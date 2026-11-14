//go:build integration

package consumer

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

type descriptor struct {
	Annotations map[string]string `json:"annotations"`
	MediaType   string            `json:"mediaType"`
	Digest      string            `json:"digest"`
}

type imageIndex struct {
	Manifests []descriptor `json:"manifests"`
}

type imageManifest struct {
	Config descriptor `json:"config"`
}

type imageConfiguration struct {
	Config struct {
		Labels map[string]string `json:"labels"`
	} `json:"config"`
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}

func TestPublishedWrapperLabels(t *testing.T) {
	t.Parallel()

	root := openImageLayout(t)
	defer closeImageLayout(t, root)

	labels := nativeImageLabels(t, root)
	want := map[string]string{
		"org.opencontainers.image.title":       "CI smoke",
		"org.opencontainers.image.description": "Buildx smoke test, native platform",
		"org.opencontainers.image.authors":     "Liakhovskyi Vladyslav",
		"org.opencontainers.image.licenses":    "MIT",
		"com.revotale.smoke":                   "preserved",
	}

	for key, value := range want {
		if labels[key] != value {
			t.Errorf("label %q = %q, want %q", key, labels[key], value)
		}
	}
}

func TestPublishedWrapperProvenance(t *testing.T) {
	t.Parallel()

	root := openImageLayout(t)
	defer closeImageLayout(t, root)

	labels := nativeImageLabels(t, root)
	if _, err := time.Parse(time.RFC3339, labels["org.opencontainers.image.created"]); err != nil {
		t.Fatalf("invalid created label: %v", err)
	}

	for name, value := range map[string]string{
		"revision": os.Getenv("CI_COMMIT_SHA"),
		"version":  strings.TrimPrefix(os.Getenv("CI_COMMIT_TAG"), "v"),
	} {
		if value == "" {
			t.Fatalf("missing expected %s; run the smoke workflow with commit and tag metadata", name)
		}

		if got := labels["org.opencontainers.image."+name]; got != value {
			t.Errorf("label %s = %q, want %q", name, got, value)
		}
	}
}

func nativeImageLabels(t *testing.T, root *os.Root) map[string]string {
	t.Helper()

	var index imageIndex
	decodeImageJSON(t, root, "index.json", &index)

	var labels map[string]string

	for _, entry := range imageDescriptors(t, root, index.Manifests) {
		var manifest imageManifest
		decodeImageJSON(t, root, blobPath(entry.Digest), &manifest)

		var config imageConfiguration
		decodeImageJSON(t, root, blobPath(manifest.Config.Digest), &config)

		if config.OS != runtime.GOOS || config.Architecture != runtime.GOARCH {
			continue
		}

		if labels != nil {
			t.Fatal("multiple native image configurations")
		}

		labels = config.Config.Labels
	}

	if labels == nil {
		t.Fatal("native image configuration with labels not found")
	}

	return labels
}

func imageDescriptors(t *testing.T, root *os.Root, entries []descriptor) []descriptor {
	t.Helper()

	var result []descriptor

	for _, entry := range entries {
		if entry.Annotations["vnd.docker.reference.type"] == "attestation-manifest" {
			continue
		}

		switch entry.MediaType {
		case "application/vnd.oci.image.index.v1+json":
			var index imageIndex
			decodeImageJSON(t, root, blobPath(entry.Digest), &index)
			result = append(result, imageDescriptors(t, root, index.Manifests)...)
		case "application/vnd.oci.image.manifest.v1+json":
			result = append(result, entry)
		}
	}

	return result
}

func decodeImageJSON(t *testing.T, root *os.Root, name string, result any) {
	t.Helper()

	data, err := root.ReadFile(name)
	if err != nil {
		t.Fatalf("read OCI entry %q: %v", name, err)
	}

	if err := json.Unmarshal(data, result); err != nil {
		t.Fatalf("decode OCI entry %q: %v", name, err)
	}
}

func blobPath(digest string) string {
	return "blobs/" + strings.ReplaceAll(digest, ":", "/")
}

func openImageLayout(t *testing.T) *os.Root {
	t.Helper()

	root, err := os.OpenRoot(os.Getenv("CI_SMOKE_IMAGE"))
	if err != nil {
		t.Fatalf("open OCI image layout: %v", err)
	}

	return root
}

func closeImageLayout(t *testing.T, root *os.Root) {
	t.Helper()

	if err := root.Close(); err != nil {
		t.Errorf("close OCI image layout: %v", err)
	}
}
