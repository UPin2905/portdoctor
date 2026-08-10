# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 0.1.x   | ✅ Yes    |

## Reporting a Vulnerability

Please do not report security vulnerabilities through public GitHub issues.

Email: security@portdoctor.dev (or open a private GitHub security advisory)

## Security Model

PortDoctor is a **read-only diagnostic tool** by default:

- `portdoctor <port>` and `portdoctor scan` never modify system state
- `portdoctor kill` requires explicit user confirmation (default: No)
- No network requests are made during normal operation
- No telemetry, no data collection, no remote logging
- User process information never leaves the local machine

## Privilege Requirements

PortDoctor does not require administrator/root privileges for normal use.
Some process details (e.g., working directory of system processes) may be unavailable
without elevated privileges. In this case, partial information is shown gracefully.
