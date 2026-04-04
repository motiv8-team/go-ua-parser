# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 1.x     | Yes       |
| < 1.0   | No        |

## Reporting a Vulnerability

If you discover a security vulnerability in go-ua-parser, please report it responsibly.

**Do not open a public GitHub issue for security vulnerabilities.**

Instead, please report via one of these methods:

1. **GitHub Security Advisories**: Use the [Security tab](https://github.com/motiv8-team/go-ua-parser/security/advisories/new) to create a private advisory
2. **Email**: Send details to the repository maintainers

### What to Include

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Response Timeline

- **Acknowledgment**: Within 48 hours
- **Assessment**: Within 1 week
- **Fix**: Dependent on severity, typically within 2 weeks for critical issues

### Scope

This library parses untrusted User-Agent strings. Relevant security concerns include:

- **Denial of Service**: Maliciously crafted UA strings that cause excessive CPU or memory usage
- **Panic/crash**: Input that causes an unrecovered panic
- **Memory safety**: Buffer overflows or out-of-bounds access

The library is fuzz-tested to mitigate these risks. If you find a UA string that triggers any of the above, that is a valid security report.

### Out of Scope

- UA spoofing (users can send any UA string — this is by design)
- Incorrect parsing results (these are bugs, not security issues — use regular GitHub Issues)
