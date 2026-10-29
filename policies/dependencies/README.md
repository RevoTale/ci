# Dependency policy

Shared GitHub dependency updates. Consumers follow the [integration guide](../../docs/integrate-dependency-updates.md).

Validate from the repository root with Woodpecker's `CI_REPO_OWNER`, `CI_REPO_NAME` and `CI_REPO_DEFAULT_BRANCH` available:

```sh
updatecli manifest validate --strict --config policies/dependencies/updatecli.yaml --values policies/dependencies/values.yaml
```

After validation, publish from the repository root using an existing OCI registry login:

```sh
updatecli manifest push --config policies/dependencies/updatecli.yaml --values policies/dependencies/values.yaml --policy policies/dependencies/Policy.yaml --tag ghcr.io/revotale/ci/policies/dependencies .
```

`Policy.yaml` supplies the policy version. Bump it when behavior changes; do not overwrite existing versions. Never include secrets in the bundle. Make the package public before consumers use it without registry credentials.

The runner is built and published by the existing image workflow on a repository release tag. Its version is separate from the policy version. Publish both before enabling a consumer's cron, including this repository's own updater.
