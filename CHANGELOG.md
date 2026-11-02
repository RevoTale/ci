# Changelog

## [0.1.6](https://github.com/RevoTale/ci/releases/tag/v0.1.6) - 2026-11-02

### ❤️ Thanks to all contributors! ❤️

@Liakhovskyi Vladyslav

### Misc

- fix(deps): ignore only dynamic release image references ([2f1b2d9](https://github.com/RevoTale/ci/commit/2f1b2d98b6b1810ff7c7fd07b15fe83b803e8f80))
- perf(buildx): download only runtime module dependencies ([8764459](https://github.com/RevoTale/ci/commit/87644592cc0930948c950b612e10ec28e9318f5a))
- perf(ci): reuse released Task tooling for validation ([9023bc0](https://github.com/RevoTale/ci/commit/9023bc03102a90274a113b1167f08e20ac920e7b))
- fix(ci): quote tag values before workflow filtering ([ed0ce93](https://github.com/RevoTale/ci/commit/ed0ce93b5c59064b2ccda00ee09e79947442a959))
- ci: simplify release gates and verify reusable components ([9bbde2f](https://github.com/RevoTale/ci/commit/9bbde2fe874812ee365217e206a140a833b30781))
- docs(ci): add Task check environment examples ([53051d9](https://github.com/RevoTale/ci/commit/53051d920d2745791770893687d259a0e99b2ca4))
- ci: simplify consumer workflows and reuse event filters ([e402f15](https://github.com/RevoTale/ci/commit/e402f1577ae41f9f50293f92bf4d78071754f5db))

## [0.1.5](https://github.com/RevoTale/ci/releases/tag/v0.1.5) - 2026-11-01

### ❤️ Thanks to all contributors! ❤️

@Liakhovskyi Vladyslav

### Misc

- ci: keep checks native and document optional Task steps ([6283a6f](https://github.com/RevoTale/ci/commit/6283a6fe5f2b458316a6e4420c2bec5adf0e15fb))
- docs(readme): clarify setup and illustrate CI workflows ([2c2853c](https://github.com/RevoTale/ci/commit/2c2853cb84ca0196b8d60a41f540bd7c2741bcd5))

## [0.1.4](https://github.com/RevoTale/ci/releases/tag/v0.1.4) - 2026-10-30

### ❤️ Thanks to all contributors! ❤️

@Liakhovskyi Vladyslav

### Misc

- docs(updatecli): record published runner and policy versions ([6e5b537](https://github.com/RevoTale/ci/commit/6e5b5371a8c5135afa93ececad4cbdc5b0fb380b))

## [0.1.1](https://github.com/RevoTale/ci/releases/tag/v0.1.1) - 2026-10-29

### ❤️ Thanks to all contributors! ❤️

@Liakhovskyi Vladyslav, @l-you, @oauth

### Misc

- fix(release): publish shared Updatecli policy with images ([c107fac](https://github.com/RevoTale/ci/commit/c107fac6fef6c080dadcf0c42f6970b4f566e290))
- ci: avoid duplicate branch push and pull request checks ([0e81d4c](https://github.com/RevoTale/ci/commit/0e81d4cfe0978f5dd2ae7ac052faa15bd9b5c072))
- 🎉 Release 0.1.1 [[#3](https://github.com/RevoTale/ci/pull/3)]
- Merge branch 'origin/main' into 'next-release/main' ([84a0ecb](https://github.com/RevoTale/ci/commit/84a0ecb0a75c277272ab5f1f9c8a815147518bd2))
- feat(ci): add reusable Updatecli dependency policy and runner ([8b9e876](https://github.com/RevoTale/ci/commit/8b9e87620c2e2369beae1144636a278d5f23cfb0))
- docs(ci): record decisions and make step images explicit ([1af8082](https://github.com/RevoTale/ci/commit/1af80826be1ccc031f56fd4eced630f52a035f72))

## [0.1.0](https://github.com/RevoTale/ci/releases/tag/v0.1.0) - 2026-10-27

### ❤️ Thanks to all contributors! ❤️

@Liakhovskyi Vladyslav

### Misc

- docs(ci): add validated consumer workflow and setup contract ([7c080cc](https://github.com/RevoTale/ci/commit/7c080cc84a1365280180a9870426518fee80fff2))
- fix(release): repair native esbuild before running release plugin ([8db9d4e](https://github.com/RevoTale/ci/commit/8db9d4ed205cff27bd64988dcbd1c5139e8b53e5))
- fix(build): cross-compile CI tools without QEMU ([4279c9d](https://github.com/RevoTale/ci/commit/4279c9d4408da4db7f503f2dbc45bb14728741b6))
- fix(ci): use Alpine CLI image for workflow lint ([02d41ce](https://github.com/RevoTale/ci/commit/02d41cea3371ef124aa7f4b6788b99fc58950034))
- feat(ci): add Buildx plugin with required image metadata ([d06157e](https://github.com/RevoTale/ci/commit/d06157ef721d6e3baf6e72876b74d45844acc818))
- ci: test and publish versioned Go toolchain images ([cc357c6](https://github.com/RevoTale/ci/commit/cc357c6cd2869d388302297cd551540b4a046cb4))
- docs(agents): require approval before repository changes ([dada829](https://github.com/RevoTale/ci/commit/dada8297a4f187fdedca9dd49ec7aa2185c21066))
- build: manage toolchain dependencies through go.mod ([4feea89](https://github.com/RevoTale/ci/commit/4feea89d1b251cd3ba134e7e1ff2acccc3754959))
- build: add Go and Task toolchain image ([cee1ef3](https://github.com/RevoTale/ci/commit/cee1ef3df0046ec748b498e6892d601cca52ec63))
- ci: validate workflows with official Woodpecker CLI ([0a06eae](https://github.com/RevoTale/ci/commit/0a06eaed7a507c2cfc1a32d0378672a084980e67))
- build(devcontainer): use Go Trixie and Task feature ([3ae1121](https://github.com/RevoTale/ci/commit/3ae1121354731f6c751e51158d5cca39759edef8))
- chore(license): set personal MIT copyright ([0da1c4a](https://github.com/RevoTale/ci/commit/0da1c4a290c76e8d90ad95a612f0e5ab5f578765))
