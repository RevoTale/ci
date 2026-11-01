# RevoTale CI

Shared tools and CI conventions for RevoTale projects.
Workflows stay in each repository; images and dependency policies are reused.

[Workflow example](examples/woodpecker.yaml) ·
[Dependency updates](docs/integrate-dependency-updates.md) ·
[Releases](https://github.com/RevoTale/ci/releases)

## What you can reuse

| Component | Purpose |
| --- | --- |
| `ghcr.io/revotale/ci/base:v0.1.3` | Go and Task on Debian Trixie. |
| `ghcr.io/revotale/ci/buildx:v0.1.3` | Official Buildx plugin with required image metadata. |
| `ghcr.io/revotale/ci/updatecli:v0.1.3` | Updatecli with Go and Git for dependency updates. |
| `ghcr.io/revotale/ci/policies/dependencies:0.1.0` | Shared dependency discovery and GitHub PR configuration. |

Images support ARM64 and AMD64. Policy versions are independent of image versions.

## Connect a project

1. Choose the [Docker workflow](examples/woodpecker.yaml) or
   [Task workflow](examples/task-workflow.yaml).
2. Set explicit tool image versions and project-specific settings.
3. Provide `task lint` and `task test:unit`.
4. For Docker images, add `artifact-test` and `release` targets.
   For Task commands, keep only the capabilities your project implements.
5. Configure the agent and secrets below.

The YAML belongs to your project. There is no remote workflow include,
configuration extension or generator.

Add `task test:integration`, `task build`, `task test:artifact`,
`task publish` and `task release` only when the project needs them.
Native plugins can handle image builds and publication.

See [optional steps](docs/optional-steps.md) for dependencies,
environment setup and examples of removing unused capabilities.

## How the Docker workflow runs

```mermaid
flowchart TD
    A["PR · main push · manual · v* tag"] --> L["task lint"]
    A --> U["task test:unit"]
    L --> T["artifact-test target · native platform"]
    U --> T
    T --> E{"Event"}
    E -->|"PR / push / manual"| B["Build release target · native platform"]
    E -->|"v* tag"| P["Build and publish · ARM64 + AMD64"]
```

Checks use the builder's native platform. Only release publication builds
ARM64 and AMD64. Additional check platforms are an explicit project choice.
See [optional steps](docs/optional-steps.md) for integration tests and
the Task workflow.

Working-branch pushes do not duplicate PR checks. Each release tag runs
its own checks before publishing. Other tags are ignored.

Artifact checks run during the Dockerfile build. Both targets share
the runtime stage; matching layers are reused through Buildx cache.

Cache is shared between steps in one workflow. Persistent and
cross-project caching require separate configuration.

[Real Dockerfile](images/go-task/Dockerfile) ·
[Why tests live in Dockerfile targets](ARTIFACT_TESTING.md) ·
[Multiple images with a matrix](.woodpecker/image.yaml)

## Agent and credentials

Use a dedicated CI machine and require approval for untrusted PRs.
Allow the exact Buildx image used by the workflow: our wrapper needs
the same builder privileges as the official plugin.

| Secret | Used for | Allowed events |
| --- | --- | --- |
| `registry_username`, `registry_password` | Publishing images and policies | `tag` |
| `forge_token` | Release preparation | Default-branch `push` |
| `forge_token` | Dependency PRs, when enabled | `cron`, `manual` |

Registry credentials need push access. Forge credentials need repository
contents and pull-request write access. Do not expose them to PR checks.

For release preparation, use our [workflow](.woodpecker/release.yaml)
and [configuration](release-config.ts).

## Dependency updates

```mermaid
flowchart LR
    A["Daily cron / manual"] --> B["Updatecli runner + shared policy"]
    B --> C["Dependency PR"]
    C --> D["Project CI + review"]
```

Connect with two files: an Updatecli compose file and a Woodpecker workflow.
There is no automerge.

Follow the [integration guide](docs/integrate-dependency-updates.md)
for files, registry access, secrets and cron setup.
The current policy uses GitHub; do not assign the same dependencies
to both Updatecli and another updater.

## What is enforced

Our Buildx wrapper requires `title`, `description`, `authors` and `licenses`.
It checks that values are present, not that they are correct.
A workflow that bypasses the wrapper also bypasses these checks.

Keep image versions explicit in every step so dependency updates show
exactly which steps change.

## For contributors and AI agents

Read [AGENTS.md](AGENTS.md) and [DECISIONS.md](DECISIONS.md) before proposing changes.
For project integration, start with the linked examples and guides.

Build shared functionality here. Propose extraction into
[woodpecker-toolkit](https://github.com/RevoTale/woodpecker-toolkit)
only when it provides clear value outside RevoTale.
