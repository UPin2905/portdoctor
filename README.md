# PortDoctor

> Diagnose local port conflicts without the detective work.

Your dev server says port 3000 is already in use. Instead of hunting through `netstat`, `lsof`, Task Manager, Docker and process trees — just run:

```bash
portdoctor 3000
```

PortDoctor tells you **what owns the port**, where it came from, and what you can safely do about it.

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

## Installation

### From source

```bash
git clone https://github.com/UPin2905/portdoctor
cd portdoctor
go build -o portdoctor ./cmd/portdoctor
```

### Pre-built binaries

Download from [Releases](https://github.com/UPin2905/portdoctor/releases).

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

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

Good beginner areas: framework detection, runtime detection, OS compatibility, tests, documentation.

---

## License

MIT — see [LICENSE](LICENSE).

