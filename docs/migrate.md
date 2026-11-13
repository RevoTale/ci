# Migrate a repository to Woodpecker

Start with the target repository URL. Prepare a migration PR in that repository.

1. **Inventory the current CI.** Read the target's agent rules, workflows (including existing Woodpecker YAML), Taskfiles, Dockerfiles and referenced scripts. Trace each job to the command it runs. Record triggers, services, artifacts, publishing destinations and required checks.
2. **Pick one starting workflow.** Use the [README selector](../README.md#choose-a-workflow). Adapt existing Woodpecker jobs in place. Map real checks to `task lint`, `task test:unit` and any needed [optional steps](optional-steps.md). Check Task dependencies so each check runs once. For multiple images, share validation and use a matrix for the builds.
3. **Use the project's actual settings.** Choose a [check environment](check-environments.md) that contains the required tools. Take registry paths, release tag patterns and supported platforms from existing publication settings. For OCI labels, verify the author, read the license and describe the actual image. Flag unknown values in the PR.
4. **Compare both pipelines.** Run Woodpecker strict lint and compare old and new PR checks on the same commit with publishing disabled for the new pipeline during comparison. In the PR description, show the mapping below. Explain any removed check or changed trigger and link the runs. Mark tag, cron or deployment paths pending until exercised.

| Existing job and trigger | Woodpecker step, command and trigger | Verification |
| --- | --- | --- |
| Unit tests on PRs | `unit` runs `task test:unit` on PRs | Link both runs for the same commit |

**Before switching**, list settings outside Git that still need action: Woodpecker activation, required status checks, untrusted PR approval, scheduled jobs and secret names with allowed events. Keep publish credentials out of PR checks. Switch required checks after the replacement passes. Hand each artifact or deployment to one publisher before retiring its old workflow.

For Buildx, add the exact image reference used by the workflow, including its
version and digest, to the server's privileged-plugin allowlist. A bare tag entry
does not authorize a different digest-pinned reference. Keep repository-wide
privileged-step and host-volume permissions disabled.

Use a release whose published consumer smoke passed. A GitHub release or an
existing registry tag alone does not establish this. Local verification can
prepare the migration while the server is unavailable; it does not verify the
server allowlist, secrets or publication.
