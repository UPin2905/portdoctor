# PortDoctor — Master Development Prompt

## 1. Project Overview

Build an open-source command-line utility called **PortDoctor**.

PortDoctor helps developers quickly understand **what is using a local network port, why it is running, where it came from, and what they can safely do about it**.

The core question PortDoctor should answer is:

> "What is using this port, and why?"

Existing commands such as `lsof`, `netstat`, `ss`, and `Get-NetTCPConnection` can identify ports and PIDs, but developers often need to manually combine several commands to understand the process, its parent, its working directory, the project that started it, and whether it belongs to Docker, Node.js, Python, or another development environment.

PortDoctor should combine this information into one simple, human-readable diagnostic tool.

The project must be:

- Open source
- Written in Go
- Cross-platform
- Fast
- Lightweight
- Safe by default
- Useful without configuration
- Easy to install
- Easy for contributors to understand
- Distributed as a single binary whenever possible

Supported operating systems:

- Windows
- macOS
- Linux

Do NOT require:

- AI or LLM APIs
- Cloud services
- User accounts
- Database
- Docker
- Node.js
- Python
- Telemetry

PortDoctor must work completely locally.

---

# 2. Primary Use Case

A developer runs an application and receives an error such as:

```text
Error: listen EADDRINUSE: address already in use :::3000
```

Instead of manually investigating the system, the developer runs:

```bash
portdoctor 3000
```

PortDoctor should produce something similar to:

```text
PortDoctor

Checking port 3000...

Status:      OCCUPIED
Protocol:    TCP
Address:     0.0.0.0:3000

Process
  Name:      node
  PID:       18421
  Started:   47 minutes ago
  Command:   node ./node_modules/next/dist/bin/next dev
  Directory: /Users/alice/projects/shop-web

Parent
  Name:      npm
  PID:       18390
  Command:   npm run dev

Project
  Type:      Next.js
  Directory: /Users/alice/projects/shop-web

Diagnosis
  Port 3000 is being used by a Next.js development server.

Suggested actions
  portdoctor kill 3000
  portdoctor find 3000
```

The exact output can evolve, but it must remain clean and understandable.

---

# 3. Product Philosophy

PortDoctor is NOT just another `kill-port` utility.

Its primary purpose is **diagnosis**.

The conceptual pipeline is:

```text
Port
 ↓
Socket
 ↓
PID
 ↓
Process
 ↓
Parent Process
 ↓
Command
 ↓
Working Directory
 ↓
Runtime
 ↓
Framework
 ↓
Container / Environment
 ↓
Project
```

Not every operating system will expose every field.

Missing information must never cause the entire command to fail.

Display the information that is available and gracefully omit unavailable fields.

---

# 4. MVP Scope — v0.1.0

Keep the first release deliberately small.

Implement these core commands:

```text
portdoctor <port>
portdoctor scan
portdoctor kill <port>
portdoctor find <port>
portdoctor --help
portdoctor --version
```

Do NOT add unnecessary features before the MVP works reliably.

---

# 5. Command: Inspect Port

Primary command:

```bash
portdoctor 3000
```

Alternative explicit form may also be supported:

```bash
portdoctor inspect 3000
```

The command should determine whether the port is:

```text
FREE
OCCUPIED
UNKNOWN
```

When occupied, attempt to retrieve:

```text
Port
Protocol
Listening address
PID
Process name
Executable
Command line
Working directory
Process start time
Parent PID
Parent process name
Parent command
Runtime
Framework
Container information
```

Not every value is mandatory.

---

# 6. Example Output

Example:

```text
$ portdoctor 3000

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
  This appears to be a Next.js development server
  running from ~/projects/shop-web.
```

Do not depend on emoji for meaning.

Support terminals where Unicode or colors are unavailable.

---

# 7. Command: Scan

Command:

```bash
portdoctor scan
```

Purpose:

Show useful listening ports currently active on the machine.

Example:

```text
$ portdoctor scan

PORT    PROCESS       PID      PROJECT             TYPE
3000    next-server   18421    shop-web            Next.js
5173    node          19211    admin-dashboard     Vite
5432    postgres      2281     -                   PostgreSQL
6379    redis         2298     -                   Redis
8000    uvicorn       20122    ai-api              Python
```

Do not dump hundreds of irrelevant system sockets by default.

Prioritize developer-relevant listening ports.

Provide an optional future flag such as:

```bash
portdoctor scan --all
```

but `--all` does not have to be part of v0.1.0 unless implementation is trivial.

---

# 8. Command: Kill

Command:

```bash
portdoctor kill 3000
```

PortDoctor must identify the process before terminating anything.

Display:

```text
Port 3000 is used by:

node
PID 18421
~/projects/shop-web

Terminate this process? [y/N]
```

Default must be **No**.

Never terminate a process automatically from the basic inspection command.

Support:

```bash
portdoctor kill 3000 --force
```

only if safe to implement.

Prefer graceful termination before forceful termination.

Conceptually:

```text
Graceful termination
        ↓
wait
        ↓
check process
        ↓
force termination if explicitly requested
```

Never silently kill unrelated processes.

---

# 9. Command: Find Free Port

Command:

```bash
portdoctor find 3000
```

If port 3000 is occupied, find the next reasonable available port.

Example:

```text
Port 3000 is occupied.

Nearest available port:
3001
```

Potential future syntax:

```bash
portdoctor find 3000 --count 5
```

Example:

```text
Available ports:

3001
3002
3004
3005
3006
```

For v0.1.0, returning one available port is sufficient.

---

# 10. Process Investigation

Create an abstraction for process inspection.

Conceptually:

```go
type ProcessInfo struct {
    PID          int
    Name         string
    Executable   string
    CommandLine  string
    WorkingDir   string
    ParentPID    int
    StartTime    time.Time
}
```

The actual implementation may differ.

The important requirement is to isolate operating-system-specific code.

Suggested architecture:

```text
internal/process/
    process.go
    process_windows.go
    process_linux.go
    process_darwin.go
```

Use Go build tags or platform-specific files where appropriate.

Avoid spreading OS-specific conditionals throughout the entire codebase.

---

# 11. Port Inspection

Create a reusable representation such as:

```go
type PortInfo struct {
    Port       int
    Protocol   string
    Address    string
    PID        int
    State      string
}
```

Port inspection logic should be separate from CLI rendering.

Example architecture:

```text
internal/port/
    port.go
    port_windows.go
    port_linux.go
    port_darwin.go
```

---

# 12. Runtime Detection

Attempt to identify common development runtimes.

Initial runtime detection:

```text
Node.js
Python
Go
Java
Ruby
PHP
.NET
Rust
Docker
Unknown
```

Use deterministic evidence such as:

- executable name
- command line
- parent process
- project files

Do NOT use AI.

Example:

```text
node → Node.js
python / python3 → Python
java → Java
dotnet → .NET
go → Go
cargo → Rust
ruby → Ruby
php → PHP
```

Runtime detection must be implemented as independent logic that can be extended later.

---

# 13. Framework Detection

For the MVP, implement lightweight framework detection.

Node.js:

```text
Next.js
Vite
Nuxt
NestJS
Express
React development server
Unknown Node.js application
```

Python:

```text
FastAPI
Django
Flask
Uvicorn
Gunicorn
Unknown Python application
```

Detection can use:

- process command line
- executable
- working directory
- package.json
- pyproject.toml
- requirements.txt
- known command patterns

Examples:

```text
next dev
→ Next.js

vite
→ Vite

uvicorn
→ Uvicorn / likely FastAPI

manage.py runserver
→ Django

flask run
→ Flask
```

Framework detection should never claim certainty without reasonable evidence.

If uncertain, report the runtime only.

---

# 14. Project Detection

When a working directory is available, attempt to determine the project root.

Look upward for common project markers:

```text
.git/
go.mod
package.json
pyproject.toml
requirements.txt
Cargo.toml
pom.xml
build.gradle
composer.json
Gemfile
*.sln
*.csproj
```

Stop searching after a reasonable number of parent directories.

Never recursively scan the entire filesystem.

Return something conceptually similar to:

```go
type ProjectInfo struct {
    Root      string
    Name      string
    Runtime   string
    Framework string
}
```

Project name may be determined from:

1. package/project metadata
2. directory name

Avoid expensive operations.

---

# 15. Parent Process Chain

One of PortDoctor's differentiating features should be the ability to explain where a process came from.

Example:

```text
Port 3000
└── node PID 18421
    └── npm PID 18390
        └── powershell PID 9021
```

Future command:

```bash
portdoctor why 3000
```

Possible output:

```text
Port 3000
└── node PID 18421
    └── started by npm
        └── npm run dev
            └── ~/projects/shop-web
                └── Next.js project
```

For v0.1.0, retrieving at least the immediate parent process is enough.

Design the internal API so full process-tree support can be added later.

---

# 16. Docker Detection

Docker support is useful but must NOT delay the initial working implementation.

If feasible, detect whether a port belongs to:

```text
docker-proxy
Docker Desktop
containerd
podman
```

Potential future output:

```text
Container
  Runtime: Docker
  Name:    local-postgres
  ID:      a4b3c912
  Image:   postgres:17
```

Do not require Docker to be installed.

If Docker is unavailable, PortDoctor must continue normally.

Docker integration should eventually live in something similar to:

```text
internal/container/
```

---

# 17. Windows Requirements

Windows is a first-class platform.

PortDoctor should work from:

```text
PowerShell
Command Prompt
Windows Terminal
```

Potential Windows data sources include:

```text
Get-NetTCPConnection
Get-Process
Win32_Process
netstat
tasklist
PowerShell APIs
Windows system APIs
```

Prefer reliable Go/system APIs where practical.

Avoid requiring administrator privileges for normal inspection.

If elevated privileges are required for a particular field, return partial information rather than failing completely.

Example:

```text
Working directory unavailable (permission denied)
```

Do not display a frightening error stack for expected permission limitations.

---

# 18. Linux Requirements

Potential information sources:

```text
/proc
ss
lsof
netstat
```

Prefer `/proc` or native mechanisms when practical.

Do not assume every distribution includes `lsof`.

PortDoctor should work on common distributions such as:

```text
Ubuntu
Debian
Fedora
Arch Linux
```

---

# 19. macOS Requirements

Potential sources:

```text
lsof
ps
sysctl
native APIs
```

Keep platform-specific implementation isolated.

---

# 20. CLI Design

The CLI must feel simple.

Good:

```bash
portdoctor 3000
portdoctor scan
portdoctor kill 3000
portdoctor find 3000
```

Avoid unnecessarily verbose commands such as:

```bash
portdoctor network inspect --tcp --port-number 3000
```

The common case should require the fewest keystrokes possible.

---

# 21. CLI Library

Use a mature Go CLI library if it meaningfully improves maintainability.

Preferred candidate:

```text
Cobra
```

However, if the MVP remains extremely small, the standard library may also be acceptable.

Do not add dependencies without a clear reason.

---

# 22. Output Architecture

Keep data collection separate from rendering.

Bad:

```text
OS query
→ immediately print
→ more OS query
→ immediately print
```

Preferred:

```text
collect
   ↓
normalize
   ↓
diagnose
   ↓
render
```

This allows future output formats such as:

```bash
portdoctor 3000 --json
```

Possible future JSON:

```json
{
  "port": 3000,
  "status": "occupied",
  "protocol": "tcp",
  "process": {
    "pid": 18421,
    "name": "node"
  },
  "project": {
    "name": "shop-web",
    "framework": "nextjs"
  }
}
```

JSON support is desirable but not mandatory for the first working commit.

---

# 23. Exit Codes

Use predictable exit codes.

Suggested:

```text
0 = command completed successfully
1 = general error
2 = invalid arguments
3 = permission problem
4 = process/port information unavailable
```

Do not use random exit codes.

Document them later.

---

# 24. Error Handling

Errors must be human-readable.

Bad:

```text
panic: runtime error: invalid memory address
```

Good:

```text
PortDoctor could not inspect PID 18421.

Reason:
Permission denied.

Try running the command with elevated privileges if you need
additional process information.
```

Do not expose stack traces during normal use.

Stack traces may be available through future debug functionality.

---

# 25. Safety Requirements

PortDoctor interacts with system processes, so safety matters.

Rules:

1. Inspection must never modify system state.
2. Never kill anything from `portdoctor <port>`.
3. `kill` requires explicit user action.
4. Default confirmation must be No.
5. Clearly display PID and process name before termination.
6. Handle PID disappearance/race conditions.
7. Verify the process still owns the relevant port when possible before killing it.
8. Do not use shell command construction with unsanitized user input.
9. Validate port numbers.
10. Valid port range is:

```text
1–65535
```

Reject:

```text
0
-1
65536
abc
3.14
```

with clear messages.

---

# 26. Performance

Basic inspection should feel instant.

Target:

```text
portdoctor 3000
```

should normally finish in well under one second.

Avoid:

- scanning the entire filesystem
- expensive recursive directory searches
- network requests
- unnecessary subprocess creation
- loading large configuration files

PortDoctor must never contact the internet during normal operation.

---

# 27. Privacy

PortDoctor is local-first.

Do NOT implement:

```text
telemetry
analytics
tracking
automatic crash uploads
remote logging
user identification
```

No system information should leave the machine.

This can later become a project selling point:

> Your process information never leaves your computer.

---

# 28. Proposed Project Structure

Start with a clean Go structure.

Example:

```text
portdoctor/
│
├── cmd/
│   └── portdoctor/
│       └── main.go
│
├── internal/
│   ├── cli/
│   │   ├── root.go
│   │   ├── inspect.go
│   │   ├── scan.go
│   │   ├── kill.go
│   │   └── find.go
│   │
│   ├── port/
│   │   ├── port.go
│   │   ├── port_windows.go
│   │   ├── port_linux.go
│   │   └── port_darwin.go
│   │
│   ├── process/
│   │   ├── process.go
│   │   ├── process_windows.go
│   │   ├── process_linux.go
│   │   └── process_darwin.go
│   │
│   ├── project/
│   │   └── detect.go
│   │
│   ├── runtime/
│   │   └── detect.go
│   │
│   ├── framework/
│   │   └── detect.go
│   │
│   └── output/
│       └── terminal.go
│
├── tests/
│
├── .github/
│   └── workflows/
│
├── .gitignore
├── CONTRIBUTING.md
├── LICENSE
├── README.md
├── SECURITY.md
├── go.mod
├── go.sum
└── Makefile
```

This is a guideline, not an absolute requirement.

Do not create unnecessary abstraction layers merely to match this tree.

---

# 29. Code Quality

Use idiomatic Go.

Requirements:

```bash
go fmt ./...
go vet ./...
go test ./...
```

should succeed.

Prefer:

- small packages
- clear interfaces
- descriptive names
- short functions
- explicit error handling
- minimal dependencies

Avoid:

- giant utility packages
- unnecessary global variables
- premature abstraction
- reflection unless justified
- clever code that reduces readability

Add comments where the reason behind code is not obvious.

Do not comment every trivial line.

---

# 30. Testing

Testing is required.

Unit test:

```text
port validation
runtime detection
framework detection
project root detection
output formatting
free-port search
```

Platform-specific process/port discovery should be tested where practical.

Use dependency injection or interfaces where needed to avoid making every test depend on real system processes.

Integration tests may start temporary TCP listeners.

Example Go test strategy:

```go
listener, err := net.Listen("tcp", "127.0.0.1:0")
```

Then inspect the dynamically allocated port.

Never rely on port 3000 being free during automated tests.

---

# 31. CI

Create GitHub Actions.

At minimum run:

```text
go test
go vet
go build
```

against:

```text
Ubuntu
Windows
macOS
```

when practical.

Example matrix:

```text
ubuntu-latest
windows-latest
macos-latest
```

Every pull request should be automatically tested.

---

# 32. Releases

Prepare the project to produce binaries such as:

```text
portdoctor-windows-amd64.exe
portdoctor-windows-arm64.exe

portdoctor-linux-amd64
portdoctor-linux-arm64

portdoctor-darwin-amd64
portdoctor-darwin-arm64
```

Eventually use GitHub Releases.

A release automation tool such as GoReleaser may be added if it meaningfully simplifies distribution.

Do not introduce complex release infrastructure before the application itself works.

---

# 33. README

Create a polished README.

Suggested structure:

```text
# PortDoctor

One command to understand what's using your port.

[demo]

## Why PortDoctor?

## Installation

## Quick Start

## Commands

## Examples

## Platform Support

## How It Works

## Roadmap

## Contributing

## Security

## License
```

The opening should immediately explain the problem.

Example:

```text
Your dev server says port 3000 is already in use.

Instead of hunting through netstat, lsof, Task Manager,
Docker and process trees:

    portdoctor 3000

PortDoctor tells you what owns the port, where it came from,
and what you can safely do about it.
```

Do not exaggerate features that have not been implemented.

---

# 34. Open-Source Requirements

Use an appropriate permissive license.

Preferred:

```text
MIT License
```

Create:

```text
LICENSE
CONTRIBUTING.md
SECURITY.md
CODE_OF_CONDUCT.md
```

`CODE_OF_CONDUCT.md` can be added before public launch if not necessary for the first commit.

Make contribution paths obvious.

Good beginner contribution areas:

```text
framework detection
runtime detection
OS compatibility
tests
documentation
terminal rendering
package manager installation
```

---

# 35. Contributor-Friendly Architecture

Detection systems should be easy to extend.

For example, adding Bun should not require rewriting the application.

Conceptually:

```go
Detector
    Detect(ProcessInfo, ProjectInfo) Detection
```

Potential contributor additions:

```text
Bun
Deno
Laravel
Rails
Spring Boot
ASP.NET
Phoenix
SvelteKit
Astro
Nuxt
NestJS
Podman
WSL
Docker
```

Do not overengineer a plugin system in v0.1.0.

Simply organize code so new detectors can be added cleanly.

---

# 36. Future Roadmap

Do NOT implement all of these now.

They are future possibilities.

## v0.2

```text
portdoctor why <port>
better process tree
JSON output
better framework detection
```

## v0.3

```text
Docker detection
Podman detection
WSL detection
```

## v0.4

```text
portdoctor watch
```

Example:

```bash
portdoctor watch
```

Show ports appearing and disappearing during development.

## v0.5

Detect port conflicts before starting common development commands.

Potential future feature:

```bash
portdoctor run npm run dev
```

PortDoctor checks the expected port first and explains conflicts.

Do not implement this in the MVP.

---

# 37. Important Non-Goals

PortDoctor is NOT:

- a firewall
- a packet sniffer
- a vulnerability scanner
- a network intrusion tool
- a remote port scanner
- a full process manager
- a Docker replacement
- a Kubernetes management tool
- an observability platform

Focus only on **local developer port diagnostics**.

Do not add remote host scanning to the MVP.

---

# 38. MVP Acceptance Criteria

Version `v0.1.0` can be considered ready when the following works reliably.

### Inspection

```bash
portdoctor 3000
```

can:

- validate the port
- determine whether it is free or occupied
- identify the owning PID when permitted
- identify process name
- display available process metadata
- attempt working-directory detection
- attempt parent-process detection
- produce readable output

### Scan

```bash
portdoctor scan
```

can:

- list active listening ports
- display associated processes when available
- avoid crashing because one process cannot be inspected

### Kill

```bash
portdoctor kill 3000
```

can:

- identify the owning process
- display it to the user
- request confirmation
- terminate it safely when confirmed
- gracefully handle already-terminated processes

### Find

```bash
portdoctor find 3000
```

can:

- determine whether 3000 is available
- otherwise return a nearby free port

### Engineering

The repository must:

```text
build successfully
pass tests
pass go vet
be gofmt formatted
build on Windows
build on Linux
build on macOS
contain a useful README
contain an open-source license
```

---

# 39. Development Strategy

Do NOT attempt to implement the entire specification in one giant change.

Build incrementally.

Recommended order:

### Phase 1 — Skeleton

Create:

```text
Go module
CLI entry point
version command
help
port validation
```

Verify:

```bash
go build
go test ./...
```

### Phase 2 — Basic Port Detection

Implement:

```bash
portdoctor 3000
```

Initially determine:

```text
FREE
OCCUPIED
```

Do not worry about framework detection yet.

### Phase 3 — PID and Process

Resolve:

```text
port → PID → process
```

Start with the current development OS, but preserve cross-platform architecture.

### Phase 4 — Process Context

Add:

```text
command line
working directory
parent process
start time
```

where supported.

### Phase 5 — Runtime / Project Detection

Add:

```text
Node.js
Python
Go
Java
etc.
```

Then project root detection.

### Phase 6 — Framework Detection

Add a small set:

```text
Next.js
Vite
Django
Flask
FastAPI/Uvicorn
```

### Phase 7 — Scan

Implement:

```bash
portdoctor scan
```

### Phase 8 — Kill

Implement safe termination.

### Phase 9 — Find

Implement nearest free port.

### Phase 10 — Polish

Add:

```text
tests
CI
README
documentation
release configuration
```

---

# 40. Instructions to the Coding Agent

When implementing this project:

1. Read this entire specification before making architectural decisions.

2. Do not blindly implement every future feature.

3. Focus on the smallest reliable MVP.

4. Before adding a dependency, explain why it is needed.

5. Prefer the Go standard library when practical.

6. Keep OS-specific implementation isolated.

7. Never sacrifice Windows support merely because Unix implementation is easier.

8. Never make a destructive action the default.

9. Add tests for important logic.

10. Do not fake system information.

If information cannot be determined, return:

```text
Unknown
```

or omit the field.

Never invent values.

11. Do not add telemetry.

12. Do not require network access.

13. Do not require external accounts or API keys.

14. Keep commands intuitive.

15. Optimize for developer experience.

16. Keep commits logically separated when possible.

17. Before declaring a feature complete, run:

```bash
go fmt ./...
go vet ./...
go test ./...
go build ./cmd/portdoctor
```

18. Fix failures rather than hiding them.

19. Document platform limitations.

20. Do not claim cross-platform support until the relevant builds/tests succeed.

---

# 41. First Task

Start by implementing only the initial foundation.

Do NOT implement the entire project yet.

Perform these tasks:

1. Initialize the Go project.
2. Create a clean repository structure.
3. Implement the PortDoctor CLI entry point.
4. Implement:

```bash
portdoctor --help
portdoctor --version
portdoctor <port>
```

5. Validate port numbers from `1` to `65535`.

6. Create the internal port inspection abstraction.

7. Implement the simplest reliable method to determine:

```text
Port 3000: FREE
```

or:

```text
Port 3000: OCCUPIED
```

8. Add unit tests.

9. Add a minimal README.

10. Run:

```bash
go fmt ./...
go vet ./...
go test ./...
go build ./cmd/portdoctor
```

11. Report:

- files created
- architectural decisions
- dependencies added
- tests performed
- known limitations
- recommended next task

STOP after this foundation works.

Do not proceed to process/PID detection until the basic architecture and tests are working correctly.

---

# 42. Project Identity

Project name:

**PortDoctor**

CLI command:

```bash
portdoctor
```

Suggested tagline:

> Diagnose local port conflicts without the detective work.

Alternative:

> Find out what's using your port — and why.

Repository description:

> A fast, cross-platform CLI that tells developers what is using a local port, where the process came from, and what they can safely do about it.

Primary language:

**Go**

Target users:

**Software developers**

Primary platforms:

**Windows, macOS, Linux**

License:

**MIT**

Core principle:

> Diagnose first. Modify only when explicitly requested.