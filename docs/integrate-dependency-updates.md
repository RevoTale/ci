# Updatecli as a companion tool

We use only official Updatecli releases, plugins and configuration. RevoTale/ci
does not build an Updatecli runner, publish a shared Updatecli policy or maintain
custom update logic. The previous dependency workflow and integration examples
have been removed.

## Why I am making this change

This decision and the reasoning below come directly from Vladyslav.

> I am removing our custom Updatecli setup because I want to move RevoTale
> entirely to a Go-only stack. For private projects, we use only Go. For public
> projects, pnpm is only there to build static files. PHP-related work is just a
> side activity of my job, so it is not a reason to build our dependency tooling
> around it.
>
> For the stack we actually use for the public internet, official Updatecli is
> enough, including getting timely security fixes. That is why I am removing
> everything we built for our custom Updatecli setup and keeping Updatecli as a
> companion tool, using only its official plugins and configuration.

## Official integration

Configure Updatecli directly in the repository that needs dependency updates.
Use its built-in plugins and documented configuration. Do not restore our
custom runner, shell targets, policy distribution or ecosystem-specific helpers.

Start with the official [installation instructions](https://www.updatecli.io/docs/prologue/installation/),
[configuration reference](https://www.updatecli.io/docs/core/configuration/) and
[autodiscovery reference](https://www.updatecli.io/docs/core/autodiscovery/).
For Go projects, see the official [Go crawler](https://www.updatecli.io/docs/plugins/autodiscovery/golang/).
Use the documented runtime requirements of the selected plugins.

Existing consumers must replace references to `ghcr.io/revotale/ci/updatecli`
and `ghcr.io/revotale/ci/policies/dependencies` before relying on this change.
Removing this integration does not migrate other repositories, delete published
artifacts or configure a replacement scheduled job.
