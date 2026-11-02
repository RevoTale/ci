# Releases

These workflows maintain this repository. Consumers start from `examples/`;
do not copy our whole `.woodpecker/` directory.

- `validate.yaml`: shared code and configuration checks.
- `image.yaml`: build the three images through a matrix.
- `release-pr.yaml`: prepare release PRs and tags on main pushes.
- `release.yaml`: publish policy, verify consumers, then announce the release.
- `dependencies.yaml`: daily dependency PRs.

## Flow

- Main pipelines validate the commit and prepare a release PR.
- Merging the release PR creates a tag at the checked CI commit.
- Existing Git tags are never moved.
- The tag pipeline checks and publishes images and the dependency policy.
- Consumer smoke verifies the published wrapper, base image and Updatecli runner.
- Policy access is checked without registry credentials; this checks resolution, not equality with local files.
- Bump `policies/dependencies/Policy.yaml` when policy manifests or values change. Updatecli skips already published versions.
- The forge release is created only after these checks pass.

Registry publication is not atomic across packages. A failed pipeline can
leave some artifacts published without a forge release.

## Retry

- Retry the failed tag pipeline, not the main release-PR pipeline.
- Woodpecker restarts the whole pipeline. Images are rebuilt and their registry
  tags may point to new digests; this is not a byte-identical retry.
- Do not restart successful releases. Fix released code with a new version.
- Consumers pin image digests so a tag retry cannot change their dependencies.

The main prepare-release step is not proof of publication.
The publish-release step in the tag workflow is the completion signal.

Smoke uses the current release tag because those images are its test subjects.
Consumer workflows keep explicit versions and digests.
