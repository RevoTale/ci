import assert from 'node:assert/strict';
import { execFileSync } from 'node:child_process';
import { mkdtempSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import test, { type TestContext } from 'node:test';
import config from './release-config.ts';

test('release tags pin the checked commit and cannot move', verifyReleaseTag);

function verifyReleaseTag(t: TestContext) {
  const cwd = process.cwd();
  const previousSha = process.env.CI_COMMIT_SHA;
  const directory = mkdtempSync(join(tmpdir(), 'ci-release-test-'));
  t.after(function cleanup() {
    process.chdir(cwd);
    if (previousSha === undefined) delete process.env.CI_COMMIT_SHA;
    else process.env.CI_COMMIT_SHA = previousSha;
    rmSync(directory, { recursive: true, force: true });
  });

  const { checked, newer } = createHistory(directory);
  process.env.CI_COMMIT_SHA = checked;
  assert.equal(config.beforeRelease({ nextVersion: '1.2.3' }), false);
  assert.equal(remoteTag(), checked);
  assert.equal(config.beforeRelease({ nextVersion: '1.2.3' }), false);
  assert.equal(remoteTag(), checked);
  assertTagCannotMove(checked, newer);
  assertInvalidInputs(checked);
}

function createHistory(directory: string) {
  const remote = join(directory, 'remote.git');
  const work = join(directory, 'work');
  git('init', '--bare', remote);
  git('init', '--initial-branch=main', work);
  process.chdir(work);
  git('config', 'user.name', 'CI test');
  git('config', 'user.email', 'ci@example.invalid');
  git('remote', 'add', 'origin', remote);
  git('commit', '--allow-empty', '-m', 'checked commit');
  const checked = git('rev-parse', 'HEAD');
  git('commit', '--allow-empty', '-m', 'newer untested commit');
  const newer = git('rev-parse', 'HEAD');
  git('push', 'origin', 'main');
  return { checked, newer };
}

function assertTagCannotMove(checked: string, newer: string) {
  process.env.CI_COMMIT_SHA = newer;
  assert.throws(function rejectRetag() {
    config.beforeRelease({ nextVersion: '1.2.3' });
  });
  assert.equal(remoteTag(), checked);
}

function assertInvalidInputs(checked: string) {
  delete process.env.CI_COMMIT_SHA;
  assert.throws(function rejectMissingSha() {
    config.beforeRelease({ nextVersion: '1.2.3' });
  }, /CI_COMMIT_SHA/);
  process.env.CI_COMMIT_SHA = 'not-a-sha';
  assert.throws(function rejectMalformedSha() {
    config.beforeRelease({ nextVersion: '1.2.3' });
  }, /CI_COMMIT_SHA/);
  process.env.CI_COMMIT_SHA = checked;
  assert.throws(function rejectMissingVersion() {
    config.beforeRelease({});
  }, /version/);
}

function remoteTag(): string {
  return git('ls-remote', '--tags', 'origin', 'refs/tags/v1.2.3').split(/\s+/)[0];
}

function git(...args: string[]): string {
  return execFileSync('git', args, { encoding: 'utf8' }).trim();
}
