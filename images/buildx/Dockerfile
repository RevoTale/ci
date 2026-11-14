# Compile natively for each target architecture instead of running Go under QEMU.
FROM --platform=$BUILDPLATFORM golang:1.27.3-trixie@sha256:af1bfc2614ebb1dbd31af514b78740fd678850e76498762b425ecee4c5c6ea80 AS build

WORKDIR /src
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal

ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH
RUN ["go", "build", "-trimpath", "-o", "/out/buildx", "./cmd/buildx"]

# Keep the official builder and daemon startup; add metadata handling.
FROM woodpeckerci/plugin-docker-buildx:6.1.2@sha256:d26b69cbe748b8a219c707e7889d5ffb3f9214bc3c480c30fe857ff85954267c AS runtime

COPY --from=build /out/buildx /bin/plugin-docker-buildx
ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "plugin-docker-buildx"]

# Check the installed entrypoint on the selected test platform.
FROM runtime AS artifact-test
RUN ["/bin/plugin-docker-buildx", "--help"]

FROM runtime AS release
