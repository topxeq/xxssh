# xxssh

跨平台 SSH 客户端，带终端支持。

[English](README.md)

## 功能特点

- **多服务器管理**：添加、编辑、删除和保存多个 SSH 服务器配置
- **完整终端支持**：完整的终端仿真，支持色彩显示
- **客户端保活**：内置心跳/保活机制，保持连接
- **跨平台**：支持 Linux、macOS 和 Windows

## 安装

### 一键安装 (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/topxeq/xxssh/master/install.sh | bash
```

或使用 wget：

```bash
wget -qO- https://raw.githubusercontent.com/topxeq/xxssh/master/install.sh | bash
```

### Windows

**PowerShell（推荐）：**

```powershell
irm https://github.com/topxeq/xxssh/releases/latest/download/xxssh-windows-amd64.exe -o xxssh.exe
```

**CMD：**

```cmd
curl -fsSL https://github.com/topxeq/xxssh/releases/latest/download/xxssh-windows-amd64.exe -o xxssh.exe
```

**Git Bash / MSYS2：**

```bash
curl -fsSL https://raw.githubusercontent.com/topxeq/xxssh/master/install.sh | bash
```

下载后直接运行 `xxssh.exe`，或将其加入 PATH 环境变量。

### 手动下载

从 [Releases](https://github.com/topxeq/xxssh/releases/latest) 下载对应平台的二进制文件：

| 平台 | 架构 | 文件 |
|------|------|------|
| Linux | amd64 | [xxssh-linux-amd64](https://github.com/topxeq/xxssh/releases/latest/download/xxssh-linux-amd64) |
| macOS | amd64 | [xxssh-darwin-amd64](https://github.com/topxeq/xxssh/releases/latest/download/xxssh-darwin-amd64) |
| macOS | arm64 | [xxssh-darwin-arm64](https://github.com/topxeq/xxssh/releases/latest/download/xxssh-darwin-arm64) |
| Windows | amd64 | [xxssh-windows-amd64.exe](https://github.com/topxeq/xxssh/releases/latest/download/xxssh-windows-amd64.exe) |

下载后解压，将二进制文件放入 PATH 中。

## 快速开始

```bash
xxssh
```

## 从源码构建

```bash
git clone https://github.com/topxeq/xxssh.git
cd xxssh/repo
go build -o xxssh ./cmd/xxssh/
```

## 开源协议

见 [LICENSE](LICENSE)
