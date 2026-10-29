# Decisions

Current agreed decisions. Implementation may still be pending.

| We do | Why |
| --- | --- |
| Show changes and get approval before editing. Commit and push only when requested. | We review the actual changes together. |
| Keep workflow YAML in each repository. | Execution stays visible in Git without a server-side configuration extension or generator. |
| Prefer official plugins, images and native configuration. | Keep dependencies and custom code small. No Bash helper scripts. |
| Build shared functionality in `RevoTale/ci`. Extract to `woodpecker-toolkit` only when it helps other people too. | Sharing code between our own repos alone does not justify another dependency. |
| Write the full image and version in each step, without image anchors. | See what each step uses, change versions independently when bugs occur, and see exactly which steps a dependency bump affects in the Git diff. |
| Use Taskfile for project commands. Require `task lint` and `task test:unit`. | Local development and CI use the same commands. |
| Add `test:integration`, `build`, `test:artifact`, `publish` and `release` tasks when the project declares that capability. | Keep a consistent contract without empty tasks. Native plugins can handle image builds and publication. |
| Prefer Debian Trixie for our toolchain. Allow Alpine where appropriate and scratch or distroless for suitable runtime images. | Keep environments consistent without forcing a full distribution into every image. |
| Manage Task through Go module tools and install its binary into the toolchain image. | Keep the tool version in the module and make `task` directly available. |
| Run independent checks in parallel. Make image workflows depend on validation. | Save time without letting builds bypass required checks. |
| Test the assembled environment in an `artifact-test` Dockerfile target. Build `release` from the same runtime stage. | Reuse matching build layers without a temporary registry or custom image-loading scripts. See [the reasoning](ARTIFACT_TESTING.md). |
| Test artifacts on ARM64 by default. Make testing additional architectures optional. | One platform is enough for the default tool checks. |
| Publish ARM64 and AMD64 images under one release tag. Cross-compile Go where possible. | Consumers pull the matching platform without requiring emulation for compilation. |
| Use explicit Buildx cache export and import between steps. | Reuse layers within the workflow. Persistent or cross-project cache remains an operator choice. |
| Declare each image's context, Dockerfile and full registry path. | Support multiple images without requiring an `IMAGE_PREFIX` convention. |
| Require `title`, `description`, `authors` and `licenses` in our Buildx wrapper. | Missing metadata fails early. Reuse plugin-provided labels instead of duplicating them. |
| Declare the actual author in image metadata. Use `Liakhovskyi Vladyslav` for these images. | Authorship should be explicit, not inferred from the organization or CI actor. |
| Use Go for necessary custom logic, with named functions, strict golangci-lint, gofmt and go fix checks. | Keep the small amount of custom code readable and checked. |
| Validate PR, push and manual runs without publishing. Publish releases from tags after checks pass. | Publication has an explicit trigger and uses the same validation gates. |
| Use `v`-prefixed release tags by default. Let consumers change the tag filter. | Keep our convention without imposing it on every project. |
| Use the official ready-release-go plugin for release preparation. | Reuse release PR and changelog automation. |
| Add daily Updatecli dependency PRs. | Keep updates visible and let CI check them before merging. |
| Keep forge and registry credentials in Woodpecker secrets with access limited to the events that need them. | Keep credentials out of Git and untrusted PR checks; allow platforms to change. |
| Require approval before running untrusted contributions. Keep privileged builders on a separate CI machine. | Approval controls execution; the separate machine limits exposure of the homelab. Infrastructure setup stays separate from this repo. |
| Mark AI involvement as `AI: advisory`. | Keep attribution consistent with our working agreement. |
