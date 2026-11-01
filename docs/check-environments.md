# Check environments

## Commands in workflow steps

Use `ci/base` when Go and Task are enough.

For additional tools, use a versioned project toolchain image or install
the tools explicitly before invoking tasks. A published toolchain avoids
repeating installation but needs its own build and update process.

Use Woodpecker services for databases. Integration tasks must wait for
readiness and use isolated test data.

The services example uses PostgreSQL and DSN names from the fiscalization
service. Adapt these names to the consumer's Taskfile.

## Commands in Dockerfile targets

Use targets when tool installation and checks should live together in Git
without publishing a separate toolchain image.

Each target invokes a Task command. Independent targets must not inherit
from one another, otherwise checks can run twice.

The example covers lint and unit tests. It assumes the tasks need only Go
and Task; add the project's tools to the tools stage.

A database inside a build needs explicit setup and cleanup. Do not assume
a Woodpecker service hostname is reachable from a Buildx RUN instruction.
Prefer the services example for database-backed checks unless the project
already owns a self-contained database test stage.

## Cache

Docker can reuse a successful RUN layer without executing the check again.
Use `no_cache: true` on a check step when a fresh execution is required;
this also disables reuse of its other build layers.

The example does not configure persistent cache or shared cache between
parallel Buildx steps. Configure cache explicitly when needed.

## Connect build and publication

Copy one checks example to `.woodpecker/validate.yaml`.

Remove lint and unit steps from the selected artifact workflow, together
with references to those steps. Add this at workflow level:

    depends_on: [validate]

Use compatible event filters in both workflows.

Keep artifact build, artifact testing and publication together.
Validation does not produce files consumed by another workflow.

Do not also invoke `task ci` if it repeats the checks already declared here.
