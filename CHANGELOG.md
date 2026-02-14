# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased] - 2026-02-08

### Added
- **Smart Merge Detection**: Implemented accurate merge detection using `git merge-base --is-ancestor` CLI command instead of go-git history walking. This correctly handles:
  - Regular merge commits
  - Squash merges
  - Rebase merges
  - All git merge strategies
- **Configuration File Support**: Added `~/.branch-clean.yaml` configuration file support with YAML parsing
  - Configure default `stale_days` value
  - Configure default `protected` branch patterns
  - Command-line flags override config file values
- **Remote Branch Deletion**: New `--remote` flag to delete branches from remote repository (origin)
- **Non-Interactive Flags**:
  - `--force` / `-f`: Skip confirmation prompt
  - `--yes` / `-y`: Automatically answer yes to all prompts
  - Perfect for CI/CD and automation scripts
- **JSON Output Format**: `--format json` flag on `list` command for machine-readable output
  - Enables scripting and integration with other tools
  - Includes all branch metadata (name, is_merged, is_stale, last_commit, protected)
- **Version Command**: `branch-clean version` to display current version
- **Custom Error Types**: Implemented proper error types for better error handling:
  - `ProtectedBranchError` for protected branch violations
  - `ErrCurrentBranch` for current branch deletion attempts
  - `ErrDefaultBranch` for default branch deletion attempts
  - `ErrCancelled` for user cancellations
- **Comprehensive Test Suite**:
  - Added tests for merge detection edge cases
  - Added tests for configuration loading/saving
  - Added tests for filter logic with all combinations
  - Added tests for protected branch patterns
  - Added tests for error conditions
  - Total test coverage increased significantly

### Improved
- **Default Branch Detection**: Now correctly detects default branch from remote HEAD instead of using current branch
  - Tries remote origin HEAD first
  - Falls back to common branch names (main, master, develop)
  - Falls back to first available branch if no common names found
  - Prevents incorrect merge detection when on feature branches
- **Multi-Select UI**: Enhanced interactive branch selection
  - True multi-select with checkboxes `[✓]`
  - Toggle individual branches with enter key
  - Clear "Confirm selection" option
  - Larger display size (15 items)
  - Better navigation instructions
- **Error Tracking**: Proper tracking and reporting of failed deletions
  - Shows summary of successful vs failed deletions
  - Returns non-zero exit code if any deletions fail
  - Visual indicators (✓ for success, ✗ for failure)
  - Detailed error messages
- **FilterBranches Logic**: Clarified and improved filtering behavior with comprehensive documentation
  - `--merged-only`: Only merged branches
  - `--stale-only`: Only stale branches
  - Both flags: Branches that are BOTH merged AND stale
  - No flags: Branches that are merged OR stale (excludes active branches)
- **Input Validation**: Added validation for all command-line flags
  - Validates `stale-days` is positive
  - Validates output format is valid
  - Validates git repository exists before operations
- **Error Messages**: Improved error messages with actionable suggestions
  - Better context in error messages
  - Suggestions for fixing common issues
  - Clear indication of what went wrong
- **Module Path**: Fixed module path from placeholder to actual GitHub repository path

### Fixed
- **Critical: Merge Detection Bug**: Fixed incorrect merge detection that only checked if branch tip commit existed in default branch history
- **Critical: Default Branch Bug**: Fixed using current branch as default branch, which caused incorrect merge detection
- **Critical: Silent Error Handling**: Fixed silent error ignoring in `isProtected()` pattern matching
- **Critical: Failed Deletion Tracking**: Fixed issue where failed deletions didn't affect exit code
- **Missing Import**: Added missing `time` import in `ui.go`
- **Test Error Handling**: Fixed unhandled error in test setup that could cause flaky tests

### Security
- **Protected Branch Safety**: Enhanced protection against accidentally deleting important branches
  - Cannot delete current branch
  - Cannot delete default branch
  - Cannot delete pattern-protected branches
  - Proper exit code (2) for protection violations
- **Input Validation**: All user inputs are validated before operations
- **Error Wrapping**: Proper error wrapping prevents information leakage while maintaining debuggability

### Changed
- **Exit Code Handling**: Improved exit code logic with proper error type checking
  - Uses `errors.As()` and `errors.Is()` for type-safe error checking
  - No longer relies on fragile string matching
- **SelectBranches Error Handling**: No longer uses string comparison for Ctrl+C detection
  - Checks for `promptui.ErrInterrupt` first
  - Falls back to string contains check for compatibility

### Documentation
- **README**: Completely overhauled with:
  - Feature highlights
  - Installation instructions with `go install`
  - Comprehensive usage examples
  - Configuration file documentation
  - Exit code reference
  - Updated dependencies
- **Code Comments**: Added documentation comments to all exported functions
- **Test Comments**: Clear test descriptions explaining what is being tested

### Performance
- **Git Operations**: Using git CLI for merge detection may be slightly slower but is much more accurate
- **Test Suite**: Added comprehensive tests that run quickly with isolated git repositories

## Notes

This release represents a major improvement in code quality, reliability, and functionality. The critical bug fixes in merge detection and default branch detection make this release essential for all users. The new features (remote deletion, configuration file, JSON output) make the tool production-ready for both interactive use and automation.

### Breaking Changes
None - all changes are backward compatible. Existing command-line usage continues to work as before.

### Migration Guide
No migration needed. Users can optionally:
1. Create `~/.branch-clean.yaml` to set preferred defaults
2. Use new flags (`--remote`, `--force`, `--yes`, `--format json`) as needed
3. Update scripts to check exit codes for better error handling
