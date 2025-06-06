# Release Automation Guide

## Overview

This repository uses automated semantic versioning based on [Conventional Commits](https://www.conventionalcommits.org/).When PRs are merged to `main`, the version is automatically bumped based on the commit messages.

## How It Works

### Version Bumping Rules

| Commit Type                   | Version Bump  | Example                                 |
|-------------------------------|---------------|-----------------------------------------|
| `fix:`                        | Patch (0.0.X) | `fix: resolve issue with todo deletion` |
| `feat:`                       | Minor (0.X.0) | `feat: add support for todo priorities` |
| `feat!:` or `BREAKING CHANGE` | Major (X.0.0) | `feat!: change CLI argument structure`  |

### Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Common Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that don't affect code meaning (formatting, etc.)
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Performance improvement
- **test**: Adding missing tests
- **chore**: Changes to build process or auxiliary tools

### Examples

#### Patch Release (Bug Fix)
```
fix: prevent todos from being deleted accidentally

Adds confirmation dialog before deletion
```

#### Minor Release (New Feature)
```
feat: add ability to set due dates on todos

Users can now set and edit due dates for each todo item
```

#### Major Release (Breaking Change)
```
feat!: redesign CLI interface

BREAKING CHANGE: The --list flag is now --show. 
The old --add flag is replaced with the 'add' subcommand.
```

## Automated Release Process

1. **Create PR** with conventional commits
2. **PR Comment** shows what version will be released
3. **Merge PR** to main
4. **Automatic Release**:
    - Version is bumped based on commits
    - CHANGELOG.md is updated
    - Homebrew formula is updated
    - Git tag is created
    - GitHub Release is published

## Setup Instructions

1. **Create VERSION file** (if not exists):
   ```bash
   echo "0.0.0" > VERSION
   ```

2. **Install commit message linting** (optional but recommended):
   ```bash
   npm install --save-dev @commitlint/cli @commitlint/config-conventional
   ```

3. **Configure git hooks** (optional):
   ```bash
   npx husky add .husky/commit-msg 'npx --no -- commitlint --edit "$1"'
   ```

## Troubleshooting

### No release created
- Check that commits follow conventional format
- Ensure PR was merged (not squashed) to preserve commit messages
- Check GitHub Actions logs for errors

### Wrong version bump
- Review commit messages for correct type prefixes
- Check for accidental breaking change indicators

### Release failed
- Ensure GitHub Actions has write permissions
- Check that all referenced files exist (VERSION, Formula/tada.rb, etc.)

## Best Practices

1. **One feature per commit**: Makes changelog cleaner
2. **Clear descriptions**: Your commit message becomes the changelog entry
3. **Use scopes**: `feat(cli): add new command` is clearer than `feat: add new command`
4. **Document breaking changes**: Always explain what changed and how to migrate

## PR Workflow Example

1. Create feature branch:
   ```bash
   git checkout -b feat/add-priorities
   ```

2. Make changes and commit:
   ```bash
   git add .
   git commit -m "feat: add priority levels to todos

   Users can now assign low, medium, or high priority to todos.
   Todos are sorted by priority in list view."
   ```

3. Push and create PR:
   ```bash
   git push origin feat/add-priorities
   ```

4. PR comment will show: "✨ New Feature - Next Version: `0.2.0`"

5. Merge PR → Automatic release!
