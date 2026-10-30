# Dependency policy

Shared GitHub dependency updates. Consumers follow the [integration guide](../../docs/integrate-dependency-updates.md).

Validate from the repository root with Woodpecker's `CI_REPO_OWNER`, `CI_REPO_NAME` and `CI_REPO_DEFAULT_BRANCH` available:

```sh
updatecli manifest validate --strict --config policies/dependencies/updatecli.d --values policies/dependencies/values.yaml
```

Release tags publish the runner and policy after image checks pass. `Policy.yaml` controls the policy version; existing versions are not overwritten.

Bump the policy version when its behavior changes. Never include secrets in the bundle. Make the package public before consumers use it without registry credentials.
