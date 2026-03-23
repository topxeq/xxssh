# xxssh

A cross-platform SSH client with terminal support.

[中文](README_zh.md)

## Features

- **Multi-server management**: Add, edit, delete and save multiple SSH server configurations
- **Full terminal support**: Complete terminal emulation with color support
- **Client keepalive**: Built-in heartbeat/keepalive to maintain connections
- **Cross-platform**: Runs on Linux, macOS, and Windows

## Installation

### One-liner (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/topxeq/xxssh/main/repo/install.sh | bash
```

Or with wget:

```bash
wget -qO- https://raw.githubusercontent.com/topxeq/xxssh/main/repo/install.sh | bash
```

### Manual Download

Download the appropriate binary for your platform from [Releases](https://github.com/topxeq/xxssh/releases/latest):

| Platform | Architecture | File |
|----------|--------------|------|
| Linux | amd64 | [xxssh-linux-amd64](https://github.com/topxeq/xxssh/releases/latest/download/xxssh-linux-amd64) |
| macOS | amd64 | [xxssh-darwin-amd64](https://github.com/topxeq/xxssh/releases/latest/download/xxssh-darwin-amd64) |
| macOS | arm64 | [xxssh-darwin-arm64](https://github.com/topxeq/xxssh/releases/latest/download/xxssh-darwin-arm64) |
| Windows | amd64 | [xxssh-windows-amd64.exe](https://github.com/topxeq/xxssh/releases/latest/download/xxssh-windows-amd64.exe) |

Extract and place the binary in your PATH.

## Quick Start

```bash
xxssh
```

## Build from Source

```bash
git clone https://github.com/topxeq/xxssh.git
cd xxssh/repo
go build -o xxssh ./cmd/xxssh/
```

## License

See [LICENSE](LICENSE)
