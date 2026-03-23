#!/bin/bash
set -e

REPO="topxeq/xxssh"
INSTALL_DIR="${INSTALL_DIR:-$HOME/bin}"
BINARY_NAME="xxssh"

detect_os() {
    local os_info="$(uname -s 2>/dev/null || echo '')"
    case "$os_info" in
        Linux*)     echo "linux";;
        Darwin*)    echo "darwin";;
        CYGWIN*)    echo "windows";;
        MSYS*)      echo "windows";;
        MINGW*)     echo "windows";;
        MINGW64*)   echo "windows";;
        *)
            # Fallback: check environment
            if echo "$os_info" | grep -qi "mingw\|msys\|cygwin"; then
                echo "windows"
            else
                echo "unsupported"
            fi
            ;;
    esac
}

detect_arch() {
    local arch_info="$(uname -m 2>/dev/null || echo '')"
    case "$arch_info" in
        x86_64)     echo "amd64";;
        amd64)      echo "amd64";;
        aarch64|arm64) echo "arm64";;
        *)          echo "amd64";;
    esac
}

get_download_url() {
    local os=$1
    local arch=$2
    local version=$3
    local ext=""
    if [ "$os" = "windows" ]; then
        ext=".exe"
    fi
    echo "https://github.com/${REPO}/releases/download/${version}/${BINARY_NAME}-${os}-${arch}${ext}"
}

version="${1:-v1.0.0}"
os=$(detect_os)
arch=$(detect_arch)

echo "Detected OS: $os, Arch: $arch"

if [ "$os" = "unsupported" ]; then
    echo "Error: Unsupported operating system" >&2
    exit 1
fi

echo "Installing ${BINARY_NAME} ${version} for ${os}/${arch}..."

# Detect if we need .exe extension
ext=""
if [ "$os" = "windows" ]; then
    ext=".exe"
fi

tmp_dir=$(mktemp -d)
cd "$tmp_dir"

# Download the binary
url=$(get_download_url "$os" "$arch" "$version")
echo "Downloading from: $url"

if command -v curl &> /dev/null; then
    curl -fsSL "$url" -o "${BINARY_NAME}${ext}" || { echo "Download failed"; exit 1; }
elif command -v wget &> /dev/null; then
    wget -q "$url" -O "${BINARY_NAME}${ext}" || { echo "Download failed"; exit 1; }
else
    echo "Error: curl or wget is required" >&2
    exit 1
fi

# Make it executable (skip on Windows)
if [ "$os" != "windows" ]; then
    chmod +x "${BINARY_NAME}${ext}"
fi

# Create install directory if needed
mkdir -p "$INSTALL_DIR"

# Move binary to install location (rename to just xxssh, no platform suffix)
if [ "$os" = "windows" ]; then
    mv "${BINARY_NAME}${ext}" "${INSTALL_DIR}/${BINARY_NAME}.exe"
    echo "Successfully installed to ${INSTALL_DIR}/${BINARY_NAME}.exe"
else
    mv "${BINARY_NAME}${ext}" "${INSTALL_DIR}/${BINARY_NAME}"
    echo "Successfully installed to ${INSTALL_DIR}/${BINARY_NAME}"
fi

# Cleanup
cd /
rm -rf "$tmp_dir"

echo "Add ${INSTALL_DIR} to your PATH if needed"
