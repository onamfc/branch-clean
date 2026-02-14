# branch-clean

> A powerful CLI tool to safely clean up merged and stale git branches with interactive selection and comprehensive safety features.

[![CI](https://github.com/onamfc/branch-clean/actions/workflows/ci.yml/badge.svg)](https://github.com/onamfc/branch-clean/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/onamfc/branch-clean?style=flat)](https://github.com/onamfc/branch-clean/releases/latest)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/onamfc/branch-clean)](https://goreportcard.com/report/github.com/onamfc/branch-clean)
[![codecov](https://codecov.io/gh/onamfc/branch-clean/branch/main/graph/badge.svg)](https://codecov.io/gh/onamfc/branch-clean)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage](#usage)
  - [Basic Commands](#basic-commands)
  - [Interactive Mode](#interactive-mode)
  - [List Mode](#list-mode)
  - [Command-Line Flags](#command-line-flags)
- [Configuration](#configuration)
- [Common Use Cases](#common-use-cases)
- [CI/CD Integration](#cicd-integration)
- [Exit Codes](#exit-codes)
- [Filtering Logic](#filtering-logic)
- [Safety Features](#safety-features)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

`branch-clean` helps you maintain a clean git repository by identifying and removing branches that have been merged or are no longer actively developed. Unlike simple git commands, branch-clean provides:

- **Smart detection** of merged branches (handles squash and rebase merges)
- **Interactive selection** with visual feedback
- **Safety checks** to prevent accidental deletion of important branches
- **Automation support** for CI/CD pipelines
- **Configuration management** for team consistency

Perfect for developers working on large projects with many feature branches, or teams maintaining multiple repositories.

---

## Features

### 🎯 Smart Branch Detection

- **Accurate Merge Detection**: Uses `git merge-base --is-ancestor` to correctly identify merged branches
  - ✅ Regular merge commits
  - ✅ Squash merges
  - ✅ Rebase merges
  - ✅ All git merge strategies
- **Automatic Default Branch Detection**: Intelligently detects the default branch from remote HEAD
- **Stale Branch Detection**: Identifies branches with no activity based on configurable age threshold

### 🖱️ Interactive Experience

- **Multi-Select Interface**: Use checkboxes to select multiple branches at once
- **Visual Feedback**: Color-coded status indicators (merged/stale/active)
- **Confirmation Prompts**: Review selections before deletion
- **Dry-Run Mode**: Preview changes without making any modifications

### 🛡️ Safety First

- **Protected Branches**: Glob pattern matching prevents deletion of critical branches
- **Current Branch Protection**: Cannot delete the branch you're currently on
- **Default Branch Protection**: Prevents accidental deletion of main/master
- **Detailed Error Reporting**: Clear feedback when operations fail
- **Deletion Summary**: Shows count of successful vs failed deletions

### 🤖 Automation Ready

- **Non-Interactive Flags**: `--force` and `--yes` for CI/CD pipelines
- **JSON Output**: Machine-readable format for scripting
- **Configuration Files**: Team-wide defaults via YAML config
- **Proper Exit Codes**: Script-friendly error handling
- **Remote Deletion**: Clean up both local and remote branches

### 📊 Flexible Output

- **Table Format**: Human-readable colorized tables (default)
- **JSON Format**: Structured data for parsing and automation
- **Verbose Mode**: Detailed operation logging for debugging

---

## Quick Start

```bash
# Install
go install github.com/onamfc/branch-clean@latest

# Navigate to your git repository
cd /path/to/your/repo

# See what branches can be cleaned (dry-run)
branch-clean --dry-run

# Interactive cleanup of merged branches
branch-clean --merged-only

# Non-interactive cleanup of stale branches (30+ days old)
branch-clean --stale-only --force

# List all branches with status information
branch-clean list
```

---

## Installation

### Option 1: Install using `go install` (Recommended)

```bash
go install github.com/onamfc/branch-clean@latest
```

This installs the latest version directly from GitHub. Ensure `$GOPATH/bin` or `$HOME/go/bin` is in your `PATH`.

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/onamfc/branch-clean.git
cd branch-clean

# Download dependencies
go mod download

# Build the binary
go build -o branch-clean

# Install to your PATH (optional)
sudo mv branch-clean /usr/local/bin/

# Or build with version information
go build -ldflags "-X main.version=1.0.0" -o branch-clean
```

### Option 3: Download Pre-built Binary

*(Coming soon: Download pre-built binaries from GitHub Releases)*

### Verify Installation

```bash
branch-clean version
```

---

## Usage

### Basic Commands

```bash
# Interactive mode (default)
branch-clean

# List branches with status
branch-clean list

# Show version information
branch-clean version

# Get help
branch-clean --help
branch-clean list --help
```

### Interactive Mode

The default interactive mode provides a visual interface for branch selection:

```bash
branch-clean
```

**How it works:**
1. Shows all merged/stale branches with status indicators
2. Use `↑`/`↓` arrow keys to navigate
3. Press `Enter` to toggle selection (checkbox appears)
4. Navigate to "Confirm selection" and press `Enter`
5. Review the summary and confirm deletion

**Example output:**
```
Select branches to delete (↑/↓ to navigate, enter to toggle/confirm)
→ [✓] feature/old-implementation [merged]
  [ ] bugfix/login-issue [stale]
  [✓] feature/deprecated-api [merged]
  ✓ Confirm selection

You are about to delete 2 branch(es):
  - feature/old-implementation
  - feature/deprecated-api
Continue [y/N]: y

✓ Deleted local branch feature/old-implementation
✓ Deleted local branch feature/deprecated-api

Deleted 2 of 2 branches
```

### List Mode

View all branches with detailed status information:

```bash
# Table format (default)
branch-clean list

# JSON format for scripting
branch-clean list --format json

# Filter to specific branch types
branch-clean list --merged-only
branch-clean list --stale-only
```

**Example table output:**
```
Branch                          Status      Age          Last Commit
--------------------------------------------------------------------------------
feature/user-authentication     merged      15 days ago  2026-01-24
bugfix/memory-leak             stale       45 days ago  2025-12-24
feature/api-v2                 active      2 days ago   2026-02-06
```

**Example JSON output:**
```json
[
  {
    "name": "feature/user-authentication",
    "is_merged": true,
    "is_stale": false,
    "last_commit": "2026-01-24T10:30:00Z",
    "protected": false
  },
  {
    "name": "bugfix/memory-leak",
    "is_merged": false,
    "is_stale": true,
    "last_commit": "2025-12-24T14:20:00Z",
    "protected": false
  }
]
```

### Command-Line Flags

#### Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dry-run` | `-d` | `false` | Preview changes without deleting branches |
| `--stale-days` | `-s` | `30` | Days since last commit to consider branch stale |
| `--protect` | `-p` | `main, master, develop, release/*` | Protected branch patterns (glob) |
| `--merged-only` | `-m` | `false` | Only show/delete merged branches |
| `--stale-only` | | `false` | Only show/delete stale branches |
| `--verbose` | `-v` | `false` | Enable verbose output |
| `--force` | `-f` | `false` | Skip confirmation prompt |
| `--yes` | `-y` | `false` | Auto-answer yes to all prompts |
| `--remote` | | `false` | Also delete branches from remote (origin) |

#### List Command Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `table` | Output format: `table` or `json` |

---

## Configuration

### Configuration File

Create `~/.branch-clean.yaml` to set default values:

```yaml
# Number of days before a branch is considered stale
stale_days: 60

# Protected branch patterns (glob syntax)
protected:
  - main
  - master
  - develop
  - staging
  - production
  - release/*
  - hotfix/*
  - feature/important-*
```

### Configuration Priority

Settings are applied in this order (highest priority first):

1. **Command-line flags** (e.g., `--stale-days 90`)
2. **Configuration file** (`~/.branch-clean.yaml`)
3. **Built-in defaults** (stale_days: 30, protected: main, master, develop, release/*)

### Example: Team Configuration

Share a config template with your team:

```bash
# Create team configuration
cat > ~/.branch-clean.yaml <<EOF
stale_days: 45
protected:
  - main
  - develop
  - staging
  - release/*
  - hotfix/*
EOF

# Team members can override specific settings
branch-clean --stale-days 30  # Override just stale_days
```

---

## Common Use Cases

### 1. Weekly Cleanup of Merged Branches

```bash
# Interactive - review before deleting
branch-clean --merged-only

# Non-interactive - for scripts
branch-clean --merged-only --force
```

### 2. Remove Branches Older Than 90 Days

```bash
# See what would be deleted
branch-clean --stale-only --stale-days 90 --dry-run

# Delete after review
branch-clean --stale-only --stale-days 90
```

### 3. Clean Up Both Local and Remote Branches

```bash
# Preview changes
branch-clean --merged-only --remote --dry-run

# Execute cleanup
branch-clean --merged-only --remote --force
```

### 4. Strict Cleanup (Merged AND Stale)

```bash
# Only delete branches that are both merged AND stale
branch-clean --merged-only --stale-only --stale-days 60
```

### 5. Export Branch Information for Analysis

```bash
# Export to JSON for analysis
branch-clean list --format json > branches.json

# Parse with jq
branch-clean list --format json | jq '.[] | select(.is_stale == true) | .name'

# Count merged branches
branch-clean list --format json | jq '[.[] | select(.is_merged == true)] | length'
```

### 6. Protect Custom Branch Patterns

```bash
# Protect all branches starting with "prod-"
branch-clean --protect "main" --protect "prod-*" --protect "release/*"

# Or add to config file
echo "protected:
  - main
  - prod-*
  - release/*" > ~/.branch-clean.yaml
```

### 7. Verbose Debugging

```bash
# See detailed operation logs
branch-clean --merged-only --verbose
```

---

## CI/CD Integration

### GitHub Actions

```yaml
name: Cleanup Stale Branches

on:
  schedule:
    - cron: '0 0 * * 0'  # Every Sunday at midnight
  workflow_dispatch:  # Allow manual trigger

jobs:
  cleanup:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
        with:
          fetch-depth: 0  # Fetch all history for all branches

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install branch-clean
        run: go install github.com/onamfc/branch-clean@latest

      - name: Cleanup merged branches
        run: |
          branch-clean --merged-only --force --remote
        env:
          GIT_AUTHOR_NAME: 'GitHub Actions'
          GIT_AUTHOR_EMAIL: 'actions@github.com'
```

### GitLab CI

```yaml
cleanup_branches:
  stage: maintenance
  image: golang:1.21
  only:
    - schedules
  script:
    - go install github.com/onamfc/branch-clean@latest
    - branch-clean --merged-only --force
  when: manual
```

### Jenkins Pipeline

```groovy
pipeline {
    agent any

    triggers {
        cron('0 0 * * 0')  // Weekly on Sunday
    }

    stages {
        stage('Cleanup Branches') {
            steps {
                script {
                    sh '''
                        go install github.com/onamfc/branch-clean@latest
                        branch-clean --merged-only --stale-only --stale-days 60 --force
                    '''
                }
            }
        }
    }

    post {
        failure {
            mail to: 'team@example.com',
                 subject: "Branch cleanup failed: ${env.JOB_NAME}",
                 body: "Check ${env.BUILD_URL} for details"
        }
    }
}
```

### Cron Job (Unix/Linux/macOS)

```bash
# Add to crontab: crontab -e
# Run every Sunday at 2 AM
0 2 * * 0 cd /path/to/repo && branch-clean --merged-only --force >> /var/log/branch-clean.log 2>&1
```

### Pre-commit Hook

```bash
# .git/hooks/pre-push
#!/bin/bash

# Remind about stale branches before push
STALE_COUNT=$(branch-clean list --stale-only --format json | jq 'length')

if [ "$STALE_COUNT" -gt 0 ]; then
    echo "⚠️  You have $STALE_COUNT stale branches. Consider running 'branch-clean' to clean up."
fi
```

---

## Exit Codes

branch-clean uses standard exit codes for integration with scripts and automation:

| Exit Code | Meaning | Description |
|-----------|---------|-------------|
| `0` | Success | All operations completed successfully |
| `1` | General Error | Git errors, validation failures, file I/O errors, etc. |
| `2` | Protected Branch | Attempted to delete protected, current, or default branch |

### Using Exit Codes in Scripts

```bash
#!/bin/bash

# Basic error handling
branch-clean --merged-only --force
if [ $? -ne 0 ]; then
    echo "ERROR: Branch cleanup failed"
    exit 1
fi

# Detailed error handling
branch-clean --merged-only --force
EXIT_CODE=$?

case $EXIT_CODE in
    0)
        echo "✓ Cleanup successful"
        ;;
    1)
        echo "✗ Cleanup failed with errors"
        exit 1
        ;;
    2)
        echo "⚠ Attempted to delete protected branch"
        exit 1
        ;;
    *)
        echo "✗ Unknown error"
        exit 1
        ;;
esac

# Continue with next steps...
```

---

## Filtering Logic

Understanding how `--merged-only` and `--stale-only` flags interact:

| Flags | Behavior |
|-------|----------|
| *None* | Shows branches that are **merged OR stale** (excludes active branches) |
| `--merged-only` | Shows **only merged** branches |
| `--stale-only` | Shows **only stale** branches |
| `--merged-only --stale-only` | Shows branches that are **both merged AND stale** |

### Examples

```bash
# Scenario 1: No filters - shows merged OR stale
branch-clean
# Shows: merged branches + stale branches
# Excludes: active unmerged branches

# Scenario 2: Merged only
branch-clean --merged-only
# Shows: only merged branches
# Excludes: stale but unmerged, active branches

# Scenario 3: Stale only
branch-clean --stale-only --stale-days 60
# Shows: only branches with no commits in 60+ days
# Excludes: merged but recent, active branches

# Scenario 4: Both filters (strictest)
branch-clean --merged-only --stale-only --stale-days 90
# Shows: only branches that are BOTH merged AND 90+ days old
# Excludes: everything else (safest option)
```

---

## Safety Features

### 1. Protected Branches

Branches are protected if they match any pattern:

```bash
# Default protected patterns
main, master, develop, release/*

# Add custom patterns
branch-clean --protect "staging" --protect "hotfix/*"

# In config file
protected:
  - main
  - master
  - develop
  - staging
  - production
  - release/*
  - hotfix/*
```

**Pattern Matching:**
- `main` - Exact match
- `release/*` - Wildcard (e.g., matches `release/v1.0`, `release/v2.0`)
- `*-production` - Suffix wildcard (e.g., matches `app-production`)

### 2. Current Branch Protection

Cannot delete the branch you're currently on:

```bash
$ git branch
* feature/my-work
  feature/old-stuff

$ branch-clean
# feature/my-work will not appear in the list
```

### 3. Default Branch Protection

The default branch (usually `main` or `master`) is automatically protected:

```bash
$ branch-clean --merged-only --force
# Will not delete main/master even if it somehow appears merged
```

### 4. Confirmation Prompts

Unless using `--force` or `--yes`, you'll always see:

```
You are about to delete 3 branch(es):
  - feature/old-feature
  - bugfix/ancient-bug
  - experiment/failed-test
Continue [y/N]:
```

### 5. Dry-Run Mode

Always test with `--dry-run` first:

```bash
$ branch-clean --merged-only --dry-run

[DRY RUN] Would delete:
  - feature/completed-feature
  - bugfix/fixed-issue
```

---

## Troubleshooting

### Issue: "not a git repository"

**Cause:** You're not in a git repository directory.

**Solution:**
```bash
# Navigate to your git repository first
cd /path/to/your/git/repo

# Verify it's a git repo
git status

# Then run branch-clean
branch-clean list
```

### Issue: "failed to determine default branch"

**Cause:** Repository has no branches or no remote configured.

**Solution:**
```bash
# Check if remote exists
git remote -v

# Add remote if missing
git remote add origin https://github.com/user/repo.git

# Fetch from remote
git fetch origin

# Or specify default branch manually (temporary workaround)
git symbolic-ref refs/remotes/origin/HEAD refs/remotes/origin/main
```

### Issue: "command not found: branch-clean"

**Cause:** `branch-clean` is not in your PATH.

**Solution:**
```bash
# Check if Go bin directory is in PATH
echo $PATH | grep go

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$PATH:$HOME/go/bin"

# Or install to system directory
sudo cp $(which branch-clean) /usr/local/bin/
```

### Issue: Merge detection is inaccurate

**Cause:** Local repository is out of sync with remote.

**Solution:**
```bash
# Fetch latest from remote
git fetch --all --prune

# Then run branch-clean
branch-clean list
```

### Issue: Cannot delete remote branches

**Cause:** No push permission or remote doesn't exist.

**Solution:**
```bash
# Verify push permission
git push --dry-run

# Check if branch exists on remote
git ls-remote --heads origin

# Delete local branch only (without --remote flag)
branch-clean --merged-only
```

### Issue: "stale-days must be positive"

**Cause:** Invalid `--stale-days` value.

**Solution:**
```bash
# Use positive integer
branch-clean --stale-days 30  # ✓ Correct
branch-clean --stale-days -7  # ✗ Wrong
branch-clean --stale-days 0   # ✗ Wrong
```

---

## FAQ

### Q: Will this delete branches I'm still working on?

**A:** No. Active branches (not merged and not stale) are automatically excluded. Use `--dry-run` to preview before deleting.

### Q: Can I undo deletions?

**A:** Local branch deletions can be recovered using git reflog:

```bash
# Find the commit hash of the deleted branch
git reflog

# Recreate the branch
git branch recovered-branch <commit-hash>
```

**Note:** Remote deletions cannot be easily undone if others have already fetched the changes.

### Q: Does this work with GitHub/GitLab/Bitbucket?

**A:** Yes! branch-clean works with any git repository, regardless of hosting provider. Use `--remote` to delete from the remote server.

### Q: How does it detect merged branches?

**A:** Uses `git merge-base --is-ancestor` which accurately detects all types of merges:
- Regular merge commits
- Squash merges
- Rebase merges
- Fast-forward merges

### Q: Can I use this in CI/CD?

**A:** Yes! Use `--force` and `--yes` flags for non-interactive operation:

```bash
branch-clean --merged-only --force --remote
```

### Q: What happens if deletion fails?

**A:** branch-clean tracks failures and reports them:

```
✓ Deleted feature/branch1
✗ Failed to delete feature/branch2: permission denied
✓ Deleted feature/branch3

Deleted 2 of 3 branches
```

Exit code will be 1 to indicate partial failure.

### Q: How do I protect specific branches?

**A:** Use the `--protect` flag or configuration file:

```bash
# Command line
branch-clean --protect "release/*" --protect "hotfix/*"

# Config file
echo "protected:
  - release/*
  - hotfix/*" > ~/.branch-clean.yaml
```

### Q: Can I see what would be deleted without actually deleting?

**A:** Yes, use `--dry-run`:

```bash
branch-clean --merged-only --dry-run
```

### Q: What's the difference between `--force` and `--yes`?

**A:**
- `--force`: Skips the confirmation prompt (but still shows selection UI)
- `--yes`: Auto-answers "yes" to all prompts
- Both are useful for automation, but `--force` is more common

---

## Development

### Prerequisites

- Go 1.21 or higher
- Git

### Building from Source

```bash
# Clone the repository
git clone https://github.com/onamfc/branch-clean.git
cd branch-clean

# Install dependencies
go mod download

# Build
go build -o branch-clean

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific tests
go test -v ./internal -run TestFilterBranches
```

### Running Tests

```bash
# All tests
go test ./...

# With verbose output
go test -v ./...

# With coverage report
go test -cover ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Specific package
go test ./internal

# Specific test
go test ./internal -run TestMergeDetection
```

### Project Structure

```
branch-clean/
├── main.go                 # CLI entry point, commands, flags
├── internal/
│   ├── git.go             # Git operations, merge detection
│   ├── git_test.go        # Git tests
│   ├── ui.go              # UI rendering, branch selection
│   ├── ui_test.go         # UI tests
│   ├── config.go          # Configuration file handling
│   └── config_test.go     # Config tests
├── go.mod                 # Dependencies
├── go.sum                 # Dependency checksums
├── README.md              # This file
├── CHANGELOG.md           # Version history
├── IMPROVEMENTS.md        # Technical improvements doc
└── LICENSE                # MIT License
```

### Adding New Features

1. **Write tests first** (TDD approach)
2. **Update documentation** (README, CHANGELOG)
3. **Follow Go conventions** (gofmt, golint)
4. **Add examples** in README

### Debugging

```bash
# Enable verbose mode
branch-clean --verbose --merged-only

# Use dry-run for safe testing
branch-clean --dry-run --verbose

# Check git operations directly
git merge-base --is-ancestor feature/branch main
echo $?  # 0 = merged, 1 = not merged
```

---

## Contributing

We welcome contributions! Here's how you can help:

### Reporting Issues

1. Check if issue already exists
2. Provide minimal reproduction steps
3. Include version information (`branch-clean version`)
4. Share relevant error messages

### Submitting Pull Requests

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Make** your changes
4. **Add tests** for new functionality
5. **Run tests** (`go test ./...`)
6. **Commit** your changes (`git commit -m 'Add amazing feature'`)
7. **Push** to your fork (`git push origin feature/amazing-feature`)
8. **Open** a Pull Request

### Code Style

- Follow standard Go formatting (`gofmt`)
- Add comments for exported functions
- Write table-driven tests
- Keep functions focused and small
- Use meaningful variable names

### Commit Messages

```
feat: Add support for custom config file locations
fix: Correct merge detection for squash merges
docs: Update installation instructions
test: Add tests for filter logic
chore: Update dependencies
```

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [promptui](https://github.com/manifoldco/promptui) for interactive prompts
- Powered by [go-git](https://github.com/go-git/go-git) for Git operations

---

## Support

- **Issues:** [GitHub Issues](https://github.com/onamfc/branch-clean/issues)
- **Discussions:** [GitHub Discussions](https://github.com/onamfc/branch-clean/discussions)
- **Documentation:** [This README](https://github.com/onamfc/branch-clean#readme)

---

**Made with ❤️ by developers, for developers**
