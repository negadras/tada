## [1.0.2](https://github.com/negadras/tada/compare/1.0.1...1.0.2) (2025-06-04)

### 🐛 Bug Fixes

* update formula to 1.0.1 and ensure script is executable ([#11](https://github.com/negadras/tada/issues/11)) ([94618c5](https://github.com/negadras/tada/commit/94618c58a9a85c4b32a6f371d8800e07bbc8cdc4))

## [1.0.1](https://github.com/negadras/tada/compare/v1.0.0...1.0.1) (2025-06-04)

### 🐛 Bug Fixes

* configure semantic-release to use non-prefixed tags ([#10](https://github.com/negadras/tada/issues/10)) ([2ec8462](https://github.com/negadras/tada/commit/2ec8462034f2d96d3fbd4996e72c112a03e67448))

## 1.0.0 (2025-06-04)

### ⚠ BREAKING CHANGES

* Manual release process is replaced with automated semantic versioning.
Tags must now be created through the automated workflow based on conventional commits.

### ✨ Features

* add automated semantic release workflow ([#9](https://github.com/negadras/tada/issues/9)) ([d3b26ef](https://github.com/negadras/tada/commit/d3b26ef25ed601b44884cb64dab4b6b4ccfc198d))
* **homebrew:** add Homebrew formula for installation ([#3](https://github.com/negadras/tada/issues/3)) ([0a4b430](https://github.com/negadras/tada/commit/0a4b4302f08334e2608fbc4161ae68670f5c31cf))
* **quote:** implement random quote command ([708e361](https://github.com/negadras/tada/commit/708e361e7648da5ec9e81e2b1cc92a61e6b5956b))

### 🐛 Bug Fixes

* explicitly specify that the PR should be created against the main branch ([#6](https://github.com/negadras/tada/issues/6)) ([2f5ad7a](https://github.com/negadras/tada/commit/2f5ad7a0308f07ab6fcbb0db024cc307d8c81a3a))
* **formula:** use tag 0.0.1 (no v-prefix) and HTTPS URL ([f93031c](https://github.com/negadras/tada/commit/f93031c4dd2c3a016b536a3da8657a0dacd5ea7f))
* specify branch name in action parameters ([ee7783f](https://github.com/negadras/tada/commit/ee7783f80bcd68f16e0272a70b97ab7ce32c4fb0))

### 👷 CI

* **release:** add release-pr.yml workflow to automate Homebrew formula bumps ([#4](https://github.com/negadras/tada/issues/4)) ([c962fd6](https://github.com/negadras/tada/commit/c962fd65b4c53501b9b7faa5093dd9523a416cd9))
* **workflow:** add GitHub Actions build workflow and update dependencies ([#2](https://github.com/negadras/tada/issues/2)) ([98efea3](https://github.com/negadras/tada/commit/98efea3de964c1a166aa2b5553fa9344bc32f6af))
