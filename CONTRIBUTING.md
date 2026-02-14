# Contributing to branch-clean

First off, thank you for considering contributing to branch-clean! It's people like you that make branch-clean such a great tool.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Commit Message Guidelines](#commit-message-guidelines)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git 2.0 or higher
- A GitHub account

### Setting Up Your Environment

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/branch-clean.git
   cd branch-clean
   ```
3. **Add the upstream repository**:
   ```bash
   git remote add upstream https://github.com/onamfc/branch-clean.git
   ```
4. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/my-new-feature
   ```

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples** (commands run, expected vs actual output)
- **Include version information** (`branch-clean version`, `go version`, OS)
- **Attach logs** (run with `--verbose` for detailed output)

Use our [bug report template](.github/ISSUE_TEMPLATE/bug_report.yml) when filing issues.

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description** of the suggested enhancement
- **Explain why this enhancement would be useful** to most users
- **Provide examples** of how the feature would be used

Use our [feature request template](.github/ISSUE_TEMPLATE/feature_request.yml).

### Your First Code Contribution

Unsure where to begin? Look for issues labeled:

- `good first issue` - Good for newcomers
- `help wanted` - Issues that need assistance
- `documentation` - Documentation improvements

### Pull Requests

1. **Discuss significant changes first**: Open an issue for substantial changes before starting work
2. **Follow the style guide**: Use `gofmt` and follow Go conventions
3. **Write tests**: Add tests for new functionality
4. **Update documentation**: Keep README and other docs in sync
5. **One feature per PR**: Keep pull requests focused on a single change

## Development Setup

### Building

```bash
# Download dependencies
go mod download

# Build the binary
go build -o branch-clean

# Build with version information
go build -ldflags "-X main.version=dev" -o branch-clean
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Testing Manually

```bash
# Test in a real git repository
cd /path/to/test/repo
/path/to/branch-clean/branch-clean --dry-run --verbose

# Test different scenarios
./branch-clean list
./branch-clean --merged-only --dry-run
./branch-clean --stale-only --stale-days 60
```

## Pull Request Process

### 1. Keep Your Fork Updated

```bash
git fetch upstream
git checkout main
git merge upstream/main
```

### 2. Create a Feature Branch

```bash
git checkout -b feature/my-feature
```

### 3. Make Your Changes

- Write clean, readable code
- Add tests for new functionality
- Update documentation as needed
- Run tests and ensure they pass

### 4. Commit Your Changes

```bash
git add .
git commit -m "feat: add new feature"
```

Follow our [commit message guidelines](#commit-message-guidelines).

### 5. Push to Your Fork

```bash
git push origin feature/my-feature
```

### 6. Open a Pull Request

- Fill out the PR template completely
- Link related issues
- Request review from maintainers
- Respond to feedback promptly

### 7. PR Review Process

- At least one maintainer approval is required
- CI tests must pass
- Code coverage should not decrease
- Documentation must be updated

## Coding Standards

### Go Style Guide

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting (run `go fmt ./...`)
- Run `go vet ./...` to catch common errors
- Use meaningful variable and function names
- Keep functions small and focused

### Code Organization

```
branch-clean/
├── main.go           # CLI entry point, commands
├── internal/         # Internal packages
│   ├── git.go       # Git operations
│   ├── ui.go        # UI and formatting
│   └── config.go    # Configuration handling
└── *_test.go        # Test files
```

### Documentation Standards

- Add godoc comments for all exported functions
- Use complete sentences in comments
- Provide usage examples for complex functionality
- Keep README.md up to date

Example:
```go
// DeleteBranch deletes a branch by name.
// Returns ErrCurrentBranch if trying to delete the currently checked out branch.
// Returns ErrDefaultBranch if trying to delete the default branch.
func (g *GitRepo) DeleteBranch(name string) error {
    // Implementation
}
```

### Error Handling

- Always handle errors explicitly
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Use custom error types for specific error conditions
- Provide actionable error messages

```go
// Good
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad
doSomething() // Ignoring error
```

## Testing Guidelines

### Test Coverage

- Aim for >80% test coverage
- Write tests for all new functionality
- Include edge cases and error conditions
- Use table-driven tests for multiple scenarios

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"case 1", "input1", "output1", false},
        {"case 2", "input2", "output2", false},
        {"error case", "bad", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Testing Best Practices

- Test behavior, not implementation
- Use descriptive test names
- Keep tests isolated and independent
- Mock external dependencies when needed
- Test both success and failure paths

## Commit Message Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, no logic change)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks, dependency updates

### Examples

```
feat(ui): add JSON output format support

Implements --format=json flag for list command.
Includes JSON marshaling with proper struct tags.

Closes #123
```

```
fix(git): correct merge detection for squash merges

Previous implementation failed for squash merges because it only
checked if branch tip commit existed in default branch history.
Now uses git merge-base --is-ancestor for accurate detection.

Fixes #456
```

```
docs: update installation instructions

- Add go install method as recommended approach
- Include instructions for building from source
- Add troubleshooting section for common issues
```

### Subject Line Rules

- Use imperative mood ("add feature" not "added feature")
- Don't capitalize first letter
- No period at the end
- Limit to 50 characters
- Be descriptive but concise

### Body

- Wrap at 72 characters
- Explain what and why, not how
- Reference issues and PRs

## Release Process

Releases are automated using GitHub Actions and GoReleaser:

1. Maintainers merge PRs to `main`
2. Update CHANGELOG.md
3. Create and push a version tag:
   ```bash
   git tag -a v1.2.3 -m "Release v1.2.3"
   git push origin v1.2.3
   ```
4. GitHub Actions automatically:
   - Runs tests
   - Builds binaries for multiple platforms
   - Creates GitHub release with artifacts
   - Updates Homebrew formula

## Questions?

- **General questions**: Open a [GitHub Discussion](https://github.com/onamfc/branch-clean/discussions)
- **Bug reports**: Use the [bug report template](.github/ISSUE_TEMPLATE/bug_report.yml)
- **Feature requests**: Use the [feature request template](.github/ISSUE_TEMPLATE/feature_request.yml)

## Recognition

Contributors will be recognized in:
- GitHub contributors page
- Release notes
- Project README (for significant contributions)

Thank you for contributing to branch-clean! 🎉
