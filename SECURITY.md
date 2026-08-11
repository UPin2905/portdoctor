# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 0.1.x   | ✅ Yes    |

## Reporting a Vulnerability

Please do not report security vulnerabilities through public GitHub issues.

Open a [private GitHub security advisory](https://github.com/UPin2905/portdoctor/security/advisories/new). Do not include vulnerabilities in public issues.

## Security Model

PortDoctor is a **read-only diagnostic tool** by default:

- `portdoctor <port>` and `portdoctor scan` never modify system state
- `portdoctor kill` requires explicit user confirmation (default: No)
- The desktop app asks for confirmation before terminating a process or exposing a port
- No telemetry, data collection, or remote logging is built in
- Process details and traffic inspector data remain local to the desktop app

## Share Feature

The desktop app's optional Share action opens an SSH reverse tunnel through `localhost.run` and makes the selected local service reachable from the internet. Only use it for services you are authorized to expose, and never use it for sensitive or production systems. The tunnel uses SSH host-key checking and is terminated when you select Stop Share or close the application.

## Privilege Requirements

PortDoctor does not require administrator/root privileges for normal use.
Some process details (e.g., working directory of system processes) may be unavailable
without elevated privileges. In this case, partial information is shown gracefully.
