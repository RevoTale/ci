# Run CI locally

Use the official Woodpecker CLI and a running Docker daemon. Authentication to
the Woodpecker server is not needed. Select the Docker backend explicitly so
workflow commands run in their configured images, not on the host.

Run from this repository's root. On macOS with Docker contexts, point the CLI
at the active daemon:

```sh
export WOODPECKER_BACKEND_DOCKER_HOST="$(docker context inspect --format '{{.Endpoints.docker.Host}}')"
```

## PR checks

```sh
woodpecker-cli --disable-update-check exec \
  --local --backend-engine docker --timeout 20m \
  --repo RevoTale/ci --repo-default-branch main \
  --pipeline-event pull_request \
  --commit-sha "$(git rev-parse HEAD)" \
  --plugins-privileged woodpeckerci/plugin-docker-buildx:6.1.2@sha256:d26b69cbe748b8a219c707e7889d5ffb3f9214bc3c480c30fe857ff85954267c \
  .woodpecker
```

This runs validation and both image matrix entries. The PR event excludes
registry publication and the release-PR workflow. Do not supply production
secrets for this check.

## Consumer smoke

Build the current sources into local-only image tags:

```sh
docker buildx build --load --target release \
  --file images/buildx/Dockerfile \
  --tag ghcr.io/revotale/ci/buildx:vlocal-metadata .
docker buildx build --load --target release \
  --file images/go-task/Dockerfile \
  --tag ghcr.io/revotale/ci/base:vlocal-metadata images/go-task
```

Reuse the release smoke workflow. Disable registry pulls for those local tags
and remove its dependency on the image workflow, which was checked separately:

```sh
sed '/^depends_on: \[image\]$/d; s/pull: true/pull: false/' \
  .woodpecker/release.yaml > /tmp/ci-local-consumer-smoke.yaml

woodpecker-cli --disable-update-check exec \
  --local --backend-engine docker --timeout 10m \
  --repo RevoTale/ci --pipeline-event tag \
  --commit-ref refs/tags/vlocal-metadata \
  --commit-sha "$(git rev-parse HEAD)" \
  --plugins-privileged ghcr.io/revotale/ci/buildx:vlocal-metadata \
  /tmp/ci-local-consumer-smoke.yaml
```

The synthetic tag is pipeline metadata; this does not create a Git tag or
publish an image. The build uses the wrapper's normal entrypoint and isolated
Docker daemon. The next step inspects the exported OCI layout with
`task test:integration`, including the description containing a comma.

Local success does not verify server scheduling, server-side permissions,
secret filters or registry publication. Published consumers still need a
successful release and explicit image digests; do not use the local tags in
their workflows.
