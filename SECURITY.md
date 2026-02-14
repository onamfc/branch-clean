# Security Policy

## Supported Versions

The following versions of branch-clean currently receive security updates:

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |
| < latest | :x:               |

We recommend always running the latest version.

## Reporting a Vulnerability

We take the security of branch-clean seriously. If you discover a security vulnerability, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

### How to Report

1. **Email**: Send a detailed report to the project maintainer via a private channel. You can reach us by opening a [private security advisory](https://github.com/onamfc/branch-clean/security/advisories/new) on GitHub.
2. **Include the following information**:
   - A description of the vulnerability
   - Steps to reproduce the issue
   - The potential impact of the vulnerability
   - Any suggested fixes, if you have them

### What to Expect

- **Acknowledgement**: We will acknowledge receipt of your report within **48 hours**.
- **Assessment**: We will investigate and assess the severity of the vulnerability within **7 days**.
- **Resolution**: We aim to release a fix for confirmed vulnerabilities within **30 days**, depending on complexity.
- **Disclosure**: We will coordinate with you on a public disclosure timeline once a fix is available.

You will be kept informed of our progress throughout the process.

## Scope

The following areas are in scope for security reports:

- **Command injection** through branch names, configuration values, or CLI arguments
- **Path traversal** via configuration file handling or git operations
- **Arbitrary code execution** through crafted git repositories or configurations
- **Sensitive data exposure** such as unintended logging of credentials or tokens
- **Dependency vulnerabilities** in direct or transitive Go module dependencies

### Out of Scope

- Vulnerabilities in git itself (report these to the [git project](https://git-scm.com/community))
- Issues that require physical access to the machine running branch-clean
- Social engineering attacks
- Denial of service through excessively large repositories (this is a local CLI tool)

## Security Best Practices for Users

- **Keep branch-clean updated** to the latest version
- **Review configuration files** (`.branch-clean.yml`) before running in new repositories, especially cloned from untrusted sources
- **Use `--dry-run`** when running against unfamiliar repositories to preview actions before making changes
- **Avoid running with elevated privileges** — branch-clean does not require root or administrator access

## Dependency Management

We monitor our Go module dependencies for known vulnerabilities using:

- GitHub Dependabot alerts
- `govulncheck` as part of our development process

When a vulnerability is discovered in a dependency, we will assess its impact on branch-clean and release a patched version as needed.

## Recognition

We appreciate the efforts of security researchers who help keep branch-clean safe. With your permission, we will acknowledge your contribution in the release notes for the version containing the fix.
