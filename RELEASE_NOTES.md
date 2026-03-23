# xxssh v1.0.0

A cross-platform SSH client with terminal support.

## Features

- Multi-server management (add, edit, delete SSH servers)
- Full terminal support with colors
- Client keepalive (heartbeat)
- Cross-platform: Linux, macOS, Windows

## Installation

### One-liner (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/topxeq/xxssh/main/install.sh | bash
```

Or with wget:

```bash
wget -qO- https://raw.githubusercontent.com/topxeq/xxssh/main/install.sh | bash
```

### Manual Download

Download the appropriate binary for your platform:

| Platform     | Architecture | File                          |
|-------------|--------------|-------------------------------|
| Linux       | amd64        | xxssh-linux-amd64             |
| macOS       | amd64        | xxssh-darwin-amd64           |
| macOS       | arm64        | xxssh-darwin-arm64           |
| Windows     | amd64        | xxssh-windows-amd64.exe      |

Extract and place the binary in your PATH.

## Usage

```bash
xxssh
```

## Downloads

- [xxssh-linux-amd64](https://github.com/topxeq/xxssh/releases/download/v1.0.0/xxssh-linux-amd64)
- [xxssh-darwin-amd64](https://github.com/topxeq/xxssh/releases/download/v1.0.0/xxssh-darwin-amd64)
- [xxssh-darwin-arm64](https://github.com/topxeq/xxssh/releases/download/v1.0.0/xxssh-darwin-arm64)
- [xxssh-windows-amd64.exe](https://github.com/topxeq/xxssh/releases/download/v1.0.0/xxssh-windows-amd64.exe)
