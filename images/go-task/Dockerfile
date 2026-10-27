# Compile the module's tools natively for the target architecture without QEMU.
FROM --platform=$BUILDPLATFORM golang:1.27.2-trixie AS build

WORKDIR /opt/tools
COPY go.mod go.sum ./
RUN ["go", "mod", "download"]

ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH
RUN ["go", "build", "-trimpath", "-o", "/out/", "tool"]

# Go projects use this image to run Task commands without installing tools in CI.
FROM golang:1.27.2-trixie AS runtime

COPY --from=build /out/ /usr/local/bin/
WORKDIR /workspace

# Test the installed tools without adding test layers to the release image.
FROM runtime AS artifact-test
RUN ["go", "version"]
RUN ["task", "--version"]

FROM runtime AS release
