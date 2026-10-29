# Integrate dependency updates

Two files connect a GitHub repository to our shared Updatecli policy. It opens one dependency PR, without automerge.

**Not published yet:** policy `0.1.0` and runner `v0.2.0` below are the planned first releases. Verify both exist before enabling this workflow.

## Connect a repository

1. Check existing Renovate, Dependabot or Updatecli configuration. Do not assign the same dependencies to two updaters.
2. Copy [updatecli-compose.yaml](../examples/dependencies/updatecli-compose.yaml) to the repository root.
3. Copy [dependencies.yaml](../examples/dependencies/.woodpecker/dependencies.yaml) to `.woodpecker/dependencies.yaml`. If the repo uses a single `.woodpecker.yaml`, first consolidate its workflows into `.woodpecker/`; do not leave competing config paths.
4. Provide the Woodpecker secret `forge_token` for **cron and manual events**, with repository contents and pull-request write permissions. A public repository still needs authentication to push branches and create PRs.
5. Create the Woodpecker cron `dependencies`, schedule `@daily`, on the default branch. YAML filters the event; it does not create the cron.
6. Run manually on the default branch. Verify the resulting PR and its normal CI before merging.

The repository owner, name and default branch come from Woodpecker. No per-project substitutions are needed. Manual runs on the default branch also run this updater; do not enable its write credentials for PR events.

## What gets updated

Dockerfile base images, Woodpecker step images, Go versions and modules, and policy references in Updatecli compose files. The runner includes Go so Updatecli can run `go mod tidy` where `go.sum` exists.

This policy does not cover every dependency file: npm, Composer, Docker Compose and other ecosystems need an explicit extension. It uses GitHub SCM and PR actions; changing only the token does not switch it to GitLab or Forgejo.

## Project exceptions

Keep the shared policy reference. Override its native crawler settings instead of copying its implementation:

```yaml
policies:
  - policy: ghcr.io/revotale/ci/policies/dependencies:0.1.0
    valuesinline:
      crawlers:
        golang:
          ignore:
            - modules:
                example.com/legacy: ""
```

Keep image and policy versions explicit. Policy updates arrive through the `updatecli` crawler; runner updates through the `woodpecker` crawler.

## For AI agents

- Read this guide, the two examples and the target repo's instructions before editing.
- Verify published versions; do not treat example versions as proof of publication.
- Preserve existing workflows and project exceptions. Show the proposed diff for approval.
- Use the shared policy. Do not add helper scripts, a custom plugin or a second updater.
- Keep `forge_token` required for PR creation. Do not add a preview mode or automerge.
- Validate the workflow, then report files, secrets, cron, publication and runtime verification separately. Writing YAML does not complete server setup.

[Updatecli policies](https://www.updatecli.io/docs/core/shareandreuse/) · [Compose overrides](https://www.updatecli.io/docs/core/compose/) · [Woodpecker cron](https://woodpecker-ci.org/docs/usage/cron)
