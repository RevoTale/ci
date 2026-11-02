# Compile natively for each target architecture instead of running Go under QEMU.
FROM --platform=$BUILDPLATFORM golang:1.27.2-trixie@sha256:47667934874913a449b96faa7ef4777fa285d9ec29126ee80ed61c45b7cccf0e AS build

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
FROM woodpeckerci/plugin-docker-buildx:6.1.2@sha256:d26b69cbe748b8a219c707e7889d5ffb3f9214bc3c480c30fe857ff85954267c AS runtime

COPY --from=build /out/buildx /usr/local/bin/ci-buildx
ENTRYPOINT ["/usr/local/bin/ci-buildx"]

# Check the installed entrypoint on the selected test platform.
FROM runtime AS artifact-test
RUN ["/usr/local/bin/ci-buildx", "--help"]

FROM runtime AS release
