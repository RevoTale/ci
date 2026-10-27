# RevoTale CI

Shared Go/Task toolchain and image-build policy for RevoTale projects.

- `ghcr.io/revotale/ci/base:v0.1.0`: Go and Task on Debian Trixie.
- `ghcr.io/revotale/ci/buildx:v0.1.0`: official Buildx 6.1.2 with required `title`, `description`, `authors`, and `licenses` settings.

## Use

Copy [the workflow](examples/woodpecker.yaml) into your project's `.woodpecker.yaml`. Set its image repository and metadata. Keep image versions explicit so Renovate can update them.

The images are reused directly. The YAML is a starting point owned by each project, not a remotely included workflow. No configuration extension or generator is required.

Provide `task lint` and `task test:unit`. Add `task test:integration`, `task build`, `task test:artifact`, `task publish`, or `task release` only when that capability belongs to the project. Native plugins handle container builds and publication.

The Dockerfile needs `artifact-test` and `release` targets sharing the same runtime stage. See [our real Dockerfile](images/go-task/Dockerfile) and [why we use this pattern](ARTIFACT_TESTING.md).

## Setup once

- Configure the agent to permit the **exact Buildx plugin image** used in the workflow. Our wrapper still starts Docker and needs the same builder privileges as the official plugin. Use an isolated CI machine; don't expose a production Docker socket to jobs.
- Provide `registry_username` and `registry_password` for **tag events only**, scoped to the repositories that publish. Give the registry account push access to their image repositories.
- For release PRs, use [our release workflow](.woodpecker/release.yaml) and [config](release-config.ts). Provide `forge_token` for **push events only**, with repository contents and pull-request write access. On GitHub, private-repository reads also need access.
- Require approval for untrusted pull requests in Woodpecker. PR checks receive no publication or release credentials.
- Keep registries and forge-specific credentials in Woodpecker. Change the registry/repository values and secret contents when moving platforms.

## Behavior

PR, push and manual runs validate and build without publishing. A `v*` tag validates, tests and publishes that exact tag. Change the tag filter if a project uses another convention. Artifact tests run on ARM64 by default; release images include ARM64 and AMD64.

The cache is a workspace directory shared by steps in one workflow. It is not a persistent or cross-project cache. Separate workflows do not share it.

Required metadata makes omissions visible before a build starts; it does not validate authorship or license correctness. Projects own those values. Plugins cannot enforce policy against workflows that bypass them.

For multiple images, use [our matrix](.woodpecker/image.yaml). Each entry declares its context, Dockerfile and complete registry path.

Develop shared functionality here. Extract it into [woodpecker-toolkit](https://github.com/RevoTale/woodpecker-toolkit) only when it provides clear value outside RevoTale.
