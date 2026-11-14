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

## Local verification, 2026-11-13

The wrapper now passes required OCI metadata directly to upstream Buildx,
without comma-separated env encoding. Missing, empty and whitespace-only
values fail before building.

- Dev Container: lint, unit tests, module verification and format/fix checks passed. Regression tests cover special characters, required fields and environment-file failures.
- Woodpecker CLI 3.19.0 with Docker: the PR workflow passed both image matrices, YAML validation, formatting, lint and unit tests. No publication step ran.
- Consumer smoke: the local wrapper built and exported an OCI image; integration checks passed for required labels, the comma-containing description and upstream provenance labels.
- File-artifact example: a temporary consumer checkout ran `examples/task-workflow.yaml` with the local base image. It built the real wrapper binary once, verified its checksum and startup in the artifact step, and skipped publication for the PR event.
- Docker example: a temporary consumer checkout ran `examples/woodpecker.yaml` with local images and the real Buildx Dockerfile. Both build targets passed through the wrapper, with a comma-containing description and cache export/import. Repository metadata was supplied through `--env` to work around the CLI 3.19.0 `--repo` bug; publication was skipped.
- Architecture coverage: ARM64 artifact and release builds plus full consumer smoke; AMD64 artifact build and executable startup under emulation.

Use [the local commands](local-verification.md) to reproduce these checks.
They do not establish registry publication or server-side permissions.

## Deployment handoff

1. Merge the implementation PR and let the normal release flow prepare its release PR.
2. Before the release tag is created, allow that exact Buildx image tag in the server's privileged-plugin list. Keep repository-wide privilege and host-volume permissions disabled.
3. Merge the release PR. Let ready-release-go create the tag, then verify publication of both architectures and the published consumer smoke.
4. Update README and consumer examples to the verified release versions and multi-platform index digests. Add the exact digest-pinned Buildx references to the server allowlist before switching consumers.

Keep existing release tags unchanged. Do not substitute local-only image tags
or unverified digests in consumer examples.
