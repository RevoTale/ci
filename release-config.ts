import { execFileSync } from 'node:child_process';

const useVersionPrefixV = true;

function pushReleaseTag({ nextVersion }: { nextVersion?: string }): false {
  const commit = process.env.CI_COMMIT_SHA;

  if (!commit || !/^(?:[0-9a-f]{40}|[0-9a-f]{64})$/i.test(commit)) {
    throw new Error('Release requires an immutable CI_COMMIT_SHA');
  }
  if (!nextVersion) {
    throw new Error('Release version is missing');
  }

  const tag =
    useVersionPrefixV && !nextVersion.startsWith('v')
      ? `v${nextVersion}`
      : nextVersion;

  execFileSync('git', ['push', '--', 'origin', `${commit}:refs/tags/${tag}`], {
    stdio: 'inherit',
  });

  // The tag pipeline publishes artifacts before creating the forge release.
  return false;
}

export default {
  useVersionPrefixV,
  skipCommitsWithoutPullRequest: false,
  beforeRelease: pushReleaseTag,
};
