# Go projects use this image to run Task commands without installing tools in CI.
FROM golang:1.27.2-trixie

# Install a fixed Task version in PATH, independently of the project's go.mod.
RUN ["go", "install", "github.com/go-task/task/v3/cmd/task@v3.54.0"]
