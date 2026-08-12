# PortDoctor

[![CI](https://github.com/UPin2905/portdoctor/actions/workflows/ci.yml/badge.svg)](https://github.com/UPin2905/portdoctor/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/UPin2905/portdoctor)](LICENSE)
[![Release](https://img.shields.io/github/v/release/UPin2905/portdoctor?display_name=tag)](https://github.com/UPin2905/portdoctor/releases)

> Diagnose local port conflicts without the detective work.

Your dev server says port 3000 is already in use. Instead of hunting through `netstat`, `lsof`, Task Manager, Docker and process trees — just run:

```bash
portdoctor 3000
```

PortDoctor tells you **what owns the port**, where it came from, and what you can safely do about it.

Port inspection, scanning, and traffic inspection run locally with no telemetry. The optional Share action intentionally exposes one selected local port through the third-party `localhost.run` tunnel service.

---

## Example

```text
🩺 PortDoctor

Port 3000 is OCCUPIED

Process
  PID          18421
  Name         node
  Executable   /usr/local/bin/node
  Started      47 minutes ago

Command
  npm run dev

Directory
  ~/projects/shop-web

Parent
  npm (PID 18390)

Detected
  Runtime      Node.js
  Framework    Next.js

Diagnosis
  This appears to be a Next.js development server running from ~/projects/shop-web.

Suggested actions
  portdoctor kill 3000
  portdoctor find 3000
```

---

## Quick Start

### CLI

```bash
git clone https://github.com/UPin2905/portdoctor.git
cd portdoctor
go run ./cmd/portdoctor scan
```

### Windows desktop app

The desktop app is supported on Windows 10/11. Install the prerequisites below, then run:

```powershell
git clone https://github.com/UPin2905/portdoctor.git
cd portdoctor\portdoctor-ui\frontend
npm ci
cd ..
wails doctor
wails dev
```

`wails doctor` must report that WebView2 is installed before starting the app. Windows 11 normally includes it; install it if the check reports it missing.

## Prerequisites

| Component | CLI | Windows desktop app |
|-----------|:---:|:-------------------:|
| Go 1.26.5 | Required | Required |
| Node.js 22 | — | Required |
| Wails CLI 2.13.0 | — | Required |
| Microsoft WebView2 Runtime | — | Required |

Install the Wails CLI after Go, then reopen your terminal if `wails` is not found:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.13.0
```

Run `wails doctor` to verify the Windows desktop prerequisites. See the [Wails installation guide](https://wails.io/docs/v2.12.0/gettingstarted/installation/) for WebView2 and platform setup details.

## Build from Source

### CLI

```powershell
git clone https://github.com/UPin2905/portdoctor
cd portdoctor
go build -o portdoctor.exe ./cmd/portdoctor
.\portdoctor.exe scan
```

```bash
git clone https://github.com/UPin2905/portdoctor
cd portdoctor
go build -o portdoctor ./cmd/portdoctor
./portdoctor scan
```

### Pre-built binaries

Download from [Releases](https://github.com/UPin2905/portdoctor/releases).

Each release includes CLI binaries for Windows, Linux, and macOS plus SHA-256 checksums.

### Windows desktop app

Build a production executable after installing the prerequisites:

```powershell
cd portdoctor-ui
wails build
```

The executable is written to `portdoctor-ui/build/bin/portdoctor-ui.exe`.

To build only the frontend:

```bash
cd portdoctor-ui/frontend
npm ci
npm run build
```

---

## Commands

| Command | Description |
|---------|-------------|
| `portdoctor <port>` | Inspect what is using a port |
| `portdoctor inspect <port>` | Same as above (explicit form) |
| `portdoctor scan` | List all active listening ports |
| `portdoctor kill <port>` | Terminate the process using a port |
| `portdoctor find <port>` | Find the nearest available port |
| `portdoctor --version` | Show version |
| `portdoctor --help` | Show help |

---

## Examples

```bash
# Inspect port 3000
portdoctor 3000

# List all developer-relevant ports
portdoctor scan

# Kill whatever is on port 8080 (prompts for confirmation)
portdoctor kill 8080

# Kill without prompt
portdoctor kill 8080 --force

# Find a free port near 3000
portdoctor find 3000
```

---

## Platform Support

| Platform | Status |
|----------|--------|
| Windows  | ✅ Supported |
| Linux    | ✅ Supported |
| macOS    | ✅ Supported |

PortDoctor works without administrator privileges for most operations.
If process details are unavailable due to permissions, partial information is displayed gracefully.

---

## How It Works

PortDoctor resolves the following chain for each port:

```
Port → Socket → PID → Process → Parent → Runtime → Framework → Project
```

It uses native OS mechanisms:

- **Windows**: `netstat`, `tasklist`, `wmic`
- **Linux**: `/proc/net/tcp`, `/proc/<pid>/`
- **macOS**: `lsof`, `ps`

No AI. No cloud. No telemetry. Your process information never leaves your computer.

---

## Roadmap

- v0.1.0 — MVP: inspect, scan, kill, find
- v0.2.0 — `portdoctor why <port>`, JSON output
- v0.3.0 — Docker / Podman / WSL detection
- v0.4.0 — `portdoctor watch`

See [CHANGELOG.md](CHANGELOG.md) for changes in progress.

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

Good beginner areas: framework detection, runtime detection, OS compatibility, tests, documentation.

For bug reports, feature requests, and support, see [SUPPORT.md](SUPPORT.md). Please follow the [Code of Conduct](CODE_OF_CONDUCT.md).

## Security

See [SECURITY.md](SECURITY.md) for private vulnerability reporting, local data handling, and the Share feature's network boundary.

---

## License

MIT — see [LICENSE](LICENSE).
