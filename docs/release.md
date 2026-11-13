# Releases

These workflows maintain this repository. Consumers start from `examples/`;
do not copy our whole `.woodpecker/` directory.

For checks without a working CI server, see [local Woodpecker execution](local-verification.md).

- `validate.yaml`: shared code and configuration checks.
- `image.yaml`: build the base and Buildx images through a matrix.
- `release-pr.yaml`: prepare release PRs and create forge releases on main pushes.
- `release.yaml`: verify the published base and Buildx images with the consumer smoke test.

## Flow

- Main pipelines validate the commit and prepare a release PR.
- After the release PR merges and main checks pass, ready-release-go creates the tag and forge release from the release branch.
- Existing Git tags are never moved.
- The tag pipeline checks and publishes the base and Buildx images.
- Consumer smoke verifies the published wrapper and base image.

Registry publication is not atomic across packages. A failed pipeline can
leave a visible forge release with only some artifacts available.

## Retry

- Retry the failed tag pipeline, not the main release-PR pipeline.
- Woodpecker restarts the whole pipeline. Images are rebuilt and their registry
  tags may point to new digests; this is not a byte-identical retry.
- Do not restart successful releases. Fix released code with a new version.
- Consumers pin image digests so a tag retry cannot change their dependencies.

A successful tag pipeline confirms artifacts are ready for consumers.

Smoke uses the current release tag because those images are its test subjects.
Consumer workflows keep explicit versions and digests.

## Local verification, 2026-11-07

- `task lint`, `task test:unit`, `task fmt` and `task fix` passed in the Dev Container; no Go source changes were required.
- Both Go modules passed `go mod verify`.
- The Buildx wrapper and Task compiled for Linux ARM64 and AMD64. Native ARM64 help/version checks passed; AMD64 binaries were not executed.
- The base and Buildx images passed local `artifact-test` and `release` Docker builds on Linux ARM64. No images were pushed.

Full AMD64 Docker image builds and the release consumer integration test were not verified in that pass.

## Metadata fix verification, 2026-11-13

The wrapper passes required OCI metadata directly to the official Buildx Go
implementation instead of encoding it into `PLUGIN_LABELS`. Missing, empty and
whitespace-only values fail before building.

- Lint, unit tests, module verification, formatting and Go fix checks passed in the Dev Container.
- Regression tests cover commas, quotes, backslashes, newlines, Unicode, existing labels and all four required fields. Restoring the old comma-splitting behavior makes them fail.
- The ARM64 `artifact-test` and `release` images built locally. The release image ran the official daemon startup, built the consumer fixture and exported its OCI layout. `task test:integration` passed against that output.
- The AMD64 `artifact-test` image also built and ran its packaged help check through local emulation. The full consumer smoke ran on ARM64.
- The packaged executable rejected a whitespace-only description before building.
- Woodpecker CLI 3.19.0 executed the consumer smoke through the Docker backend with local Buildx and base images. Both the build and the subsequent OCI integration test passed. Only registry pulling and the dependency on the separately tested image workflow were removed from the temporary smoke configuration.
- The same CLI executed `.woodpecker` locally for a PR event at commit `64c9767`: both image matrix entries passed `artifact-test` and dry-run release builds, followed by workflow lint, formatting/fix checks, lint and unit tests. The command exited with status 0; no publication or release-PR step ran.

These are local checks. The next published release still needs a successful tag
pipeline before consumer examples can pin its digests.

## Custom Updatecli removal, 2026-11-09

Removed the custom runner, policy publication, dependency workflow and their
validation and smoke steps. The remaining base and Buildx release flow is intact.
`task lint`, `task test:unit` and strict lint of the four remaining repository
workflows passed. These local checks do not establish a successful release.
