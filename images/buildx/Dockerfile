# Compile natively for each target architecture instead of running Go under QEMU.
FROM --platform=$BUILDPLATFORM golang:1.27.2-trixie AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN ["go", "mod", "download"]

COPY cmd ./cmd
COPY internal ./internal

ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH
RUN ["go", "build", "-trimpath", "-o", "/out/buildx", "./cmd/buildx"]

# Keep the official builder and daemon startup; add metadata handling.
FROM woodpeckerci/plugin-docker-buildx:6.1.2 AS runtime

COPY --from=build /out/buildx /usr/local/bin/ci-buildx
ENTRYPOINT ["/usr/local/bin/ci-buildx"]

# Check the installed entrypoint on the selected test platform.
FROM runtime AS artifact-test
RUN ["/usr/local/bin/ci-buildx", "--help"]

FROM runtime AS release
