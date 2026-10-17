# Go projects use this image to run Task commands without installing tools in CI.
FROM golang:1.27.2-trixie AS runtime

WORKDIR /opt/tools
COPY go.mod go.sum ./

# Install the tools and versions declared in this image's module.
ENV GOBIN=/usr/local/bin
RUN ["go", "install", "tool"]

WORKDIR /workspace

# Test the installed tools without adding test layers to the release image.
FROM runtime AS artifact-test
RUN ["go", "version"]
RUN ["task", "--version"]

FROM runtime AS release
