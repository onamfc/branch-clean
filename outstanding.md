1. CODE_OF_CONDUCT.md

  Standard for open source projects to set community expectations.

  2. SECURITY.md

  Security policy for reporting vulnerabilities.

  3. LICENSE - Missing Copyright Holder

  Your LICENSE file at branch-clean/LICENSE:3 has "Copyright (c) 2025" but
  is missing the copyright holder name.

  4. .golangci.yml

  Your CI workflow references golangci-lint, but there's no configuration
  file. This could cause inconsistent linting.

  Pre-Release Checklist

  5. Verify Tests Pass

  Run go test ./... to ensure all tests pass before release.

  6. Verify Build Works

  Run go build -v to ensure the project compiles successfully.

  7. Create Initial Release

  Push a git tag (e.g., v1.0.0) to trigger the release workflow and generate
   pre-built binaries.

  8. Update README Badges

  Add GitHub workflow status badges:
  - CI status badge
  - Release version badge
  - Go Report Card badge
  - Code coverage badge (if using Codecov)

  9. Repository Settings

  On GitHub, configure:
  - Repository description and topics/tags
  - Homepage URL (if applicable)
  - Repository social preview image
  - Enable/disable features (Wikis, Projects, Discussions)
  - Branch protection rules for main branch

  Optional but Recommended

  10. MAINTAINERS.md

  Document who maintains the project and how decisions are made.

  11. Release Process Documentation

  Document how to cut a new release (tagging conventions, changelog
  updates).

  12. Example Configurations

  Add example .branch-clean.yaml files in a examples/ directory.

  13. Go Report Card

  Submit your repo to https://goreportcard.com for code quality analysis.

