# Go projects use this image to run Task commands without installing tools in CI.
FROM golang:1.27.2-trixie

WORKDIR /opt/tools
COPY go.mod go.sum ./

# Install every tool declared in this image's go.mod, independently of consumers.
ENV GOBIN=/usr/local/bin
RUN ["go", "install", "tool"]

WORKDIR /workspace
