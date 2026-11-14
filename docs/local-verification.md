# Run CI locally

Use the official Woodpecker CLI and a running Docker daemon. Authentication to
the Woodpecker server is not needed. Select the Docker backend explicitly so
workflow commands run in their configured images, not on the host.

Run from this repository's root. On macOS with Docker contexts, point the CLI
at the active daemon:

```sh
export WOODPECKER_BACKEND_DOCKER_HOST="$(docker context inspect --format '{{.Endpoints.docker.Host}}')"
```

CLI 3.19.0 exposes `--repo`, but its metadata loader reads `repo-name` instead.
Supply repository variables explicitly with the native `--env` flags below;
otherwise `${CI_REPO_NAME}` expands to an empty title locally. This does not
require changing the workflow or weakening required metadata.
See the pinned [flag definitions](https://github.com/woodpecker-ci/woodpecker/blob/v3.19.0/cli/exec/flags.go)
and [metadata loader](https://github.com/woodpecker-ci/woodpecker/blob/v3.19.0/cli/exec/metadata.go).

## PR checks

```sh
woodpecker-cli --disable-update-check exec \
  --local --backend-engine docker --timeout 20m \
  --repo RevoTale/ci --repo-default-branch main \
  --env CI_REPO=RevoTale/ci --env CI_REPO_NAME=ci --env CI_REPO_OWNER=RevoTale \
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
  --tag ghcr.io/revotale/ci/buildx:v0.0.0-local .
docker buildx build --load --target release \
  --file images/go-task/Dockerfile \
  --tag ghcr.io/revotale/ci/base:v0.0.0-local images/go-task
```

Reuse the release smoke workflow. Disable registry pulls for those local tags
and remove its dependency on the image workflow, which was checked separately:

```sh
sed '/^depends_on: \[image\]$/d; s/pull: true/pull: false/' \
  .woodpecker/release.yaml > /tmp/ci-local-consumer-smoke.yaml

woodpecker-cli --disable-update-check exec \
  --local --backend-engine docker --timeout 10m \
  --repo RevoTale/ci --pipeline-event tag \
  --env CI_REPO=RevoTale/ci --env CI_REPO_NAME=ci --env CI_REPO_OWNER=RevoTale \
  --commit-ref refs/tags/v0.0.0-local \
  --commit-sha "$(git rev-parse HEAD)" \
  --plugins-privileged ghcr.io/revotale/ci/buildx:v0.0.0-local \
  /tmp/ci-local-consumer-smoke.yaml
```

The synthetic tag is pipeline metadata; this does not create a Git tag or
publish an image. The build uses the wrapper's normal entrypoint and isolated
Docker daemon. The next step inspects the exported OCI layout with
`task test:integration`, including the description containing a comma.
It also checks the creation timestamp, commit revision and release version.
Commit and tag metadata are required; use the CLI command above instead of
running the integration task without its workflow environment.

Local success does not verify server scheduling, server-side permissions,
secret filters or registry publication. Published consumers still need a
successful release and explicit image digests; do not use the local tags in
their workflows.
