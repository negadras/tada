## [2.3.0](https://github.com/negadras/tada/compare/2.2.0...2.3.0) (2025-07-15)

### ✨ Features

* **tui:** implement add functionality for todos and quotes ([#40](https://github.com/negadras/tada/issues/40)) ([9771f32](https://github.com/negadras/tada/commit/9771f32d15c2e624bc95afdc6a40668f6b477a58))

## [2.2.0](https://github.com/negadras/tada/compare/2.1.1...2.2.0) (2025-07-14)

### ✨ Features

* **quote:** migrate quotes to database with full CRUD operations ([#39](https://github.com/negadras/tada/issues/39)) ([95625a6](https://github.com/negadras/tada/commit/95625a69e3747065e466a4030512fa8462f2c019))

## [2.1.1](https://github.com/negadras/tada/compare/2.1.0...2.1.1) (2025-07-14)

### 🐛 Bug Fixes

* **deps:** update module github.com/charmbracelet/bubbletea to v1.3.6 ([#38](https://github.com/negadras/tada/issues/38)) ([7726190](https://github.com/negadras/tada/commit/7726190087311a84e5ef420e903f80b49f7200ea))

## [2.1.0](https://github.com/negadras/tada/compare/2.0.0...2.1.0) (2025-07-04)

### ✨ Features

* **ui:** add interactive table view for list command using Bubble Tea ([#32](https://github.com/negadras/tada/issues/32)) ([e9ede61](https://github.com/negadras/tada/commit/e9ede61e68351e3d2c1e31d6f63a9289aad58636))

### ♻️ Refactoring

* extract delete command into separate package ([#36](https://github.com/negadras/tada/issues/36)) ([28f5da0](https://github.com/negadras/tada/commit/28f5da0eb4ff51f4485f6d5c6f8cd289d8448748))
* extract delete command into separate package ([#36](https://github.com/negadras/tada/issues/36)) ([#37](https://github.com/negadras/tada/issues/37)) ([7bd7aa3](https://github.com/negadras/tada/commit/7bd7aa3facda96047a97e3ad7cf10af5a4df968f))

### ✅ Tests

* add test coverage for internal todo and helpers package ([#35](https://github.com/negadras/tada/issues/35)) ([36a822d](https://github.com/negadras/tada/commit/36a822d0695affa839d72b0639329ee39045e57c))
* add unit tests for commands and improve documentation ([#34](https://github.com/negadras/tada/issues/34)) ([8a3c478](https://github.com/negadras/tada/commit/8a3c478d862165b8d35c08121101b2fe6ae0dc44))

## [2.0.0](https://github.com/negadras/tada/compare/1.4.0...2.0.0) (2025-06-09)

### ⚠ BREAKING CHANGES

* None - all commands work the same way

### ✨ Features

* simplify todo app architecture and reduce over-engineering ([#31](https://github.com/negadras/tada/issues/31)) ([225fefc](https://github.com/negadras/tada/commit/225fefc06b0d83bf3d251c46d3175e1516ef7b27))

## [1.4.0](https://github.com/negadras/tada/compare/1.3.0...1.4.0) (2025-06-09)

### ✨ Features

* **cmd:** add delete command to remove todos ([#29](https://github.com/negadras/tada/issues/29)) ([eba7ce0](https://github.com/negadras/tada/commit/eba7ce0a5a4077edaa6a065cf6c70fff8efb2e04))

## [1.3.0](https://github.com/negadras/tada/compare/1.2.0...1.3.0) (2025-06-09)

### ✨ Features

* **cmd:** add done command to mark todos as completed ([#28](https://github.com/negadras/tada/issues/28)) ([1268669](https://github.com/negadras/tada/commit/1268669449f6416d3c71c6eb3e3c515edeb8adda))

### ♻️ Refactoring

* restructure project with improved architecture ([#27](https://github.com/negadras/tada/issues/27)) ([cbf1930](https://github.com/negadras/tada/commit/cbf19302e3507ee3c40262ba6d503303188e377c))

## [1.2.0](https://github.com/negadras/tada/compare/1.1.1...1.2.0) (2025-06-09)

### ✨ Features

* implement list todo tasks with optional filter ([#26](https://github.com/negadras/tada/issues/26)) ([35d2a3b](https://github.com/negadras/tada/commit/35d2a3b1aac9ec05b865991a61a5b6a96008fe0d))

## [1.1.1](https://github.com/negadras/tada/compare/1.1.0...1.1.1) (2025-06-09)

### 🐛 Bug Fixes

* **ci:** repair auto-release workflow configuration ([#25](https://github.com/negadras/tada/issues/25)) ([67d3a2e](https://github.com/negadras/tada/commit/67d3a2e84fe45d75ad79c653de2afd69c37b7433))

## [1.1.0](https://github.com/negadras/tada/compare/1.0.4...1.1.0) (2025-06-09)

### ✨ Features

* implement SQLite db and add command ([#24](https://github.com/negadras/tada/issues/24)) ([df3cbf7](https://github.com/negadras/tada/commit/df3cbf7b9840a9145595182526ae182d325e58ad))

## [1.0.4](https://github.com/negadras/tada/compare/1.0.3...1.0.4) (2025-06-07)

### 🐛 Bug Fixes

* **commitlint:** use .mjs extension for config file ([#21](https://github.com/negadras/tada/issues/21)) ([642c641](https://github.com/negadras/tada/commit/642c641324a3db15b1cc2b422eb9492117eade6d))

### ♻️ Refactoring

* some refacotr and adding a task list to be done ([#15](https://github.com/negadras/tada/issues/15)) ([73c600b](https://github.com/negadras/tada/commit/73c600b40e17d098248ba0f3692ddda594078c11))

## [1.0.3](https://github.com/negadras/tada/compare/1.0.2...1.0.3) (2025-06-04)

### 🐛 Bug Fixes

* update-files script was auto-release action content ([#12](https://github.com/negadras/tada/issues/12)) ([6b7a5f8](https://github.com/negadras/tada/commit/6b7a5f88fc634fde051194a3a0b99d1a7971c2a3))

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
