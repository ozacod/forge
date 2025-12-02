#!/bin/sh
# cpx installer - C++ Project Generator
# Usage: sh -c "$(curl -fsSL https://raw.githubusercontent.com/ozacod/cpx/master/install.sh)"

set -e

REPO="ozacod/cpx"
BINARY_NAME="cpx"

# Determine install directory (prefer user-local, fallback to system with sudo)
get_install_dir() {
    # Check for user-specified directory
    if [ -n "$CPX_INSTALL_DIR" ]; then
        echo "$CPX_INSTALL_DIR"
        return
    fi
    
    # Prefer ~/.local/bin (no sudo needed)
    LOCAL_BIN="$HOME/.local/bin"
    if [ -d "$LOCAL_BIN" ] && [ -w "$LOCAL_BIN" ]; then
        echo "$LOCAL_BIN"
        return
    fi
    
    # Check if /usr/local/bin is writable
    if [ -w "/usr/local/bin" ]; then
        echo "/usr/local/bin"
        return
    fi
    
    # Default to ~/.local/bin (will be created)
    echo "$LOCAL_BIN"
}

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_banner() {
    printf "\n"
    printf "%b   ██████╗██████╗ ██╗  ██╗%b\n" "$CYAN" "$NC"
    printf "%b  ██╔════╝██╔══██╗╚██╗██╔╝%b\n" "$CYAN" "$NC"
    printf "%b  ██║     ██████╔╝ ╚███╔╝ %b\n" "$CYAN" "$NC"
    printf "%b  ██║     ██╔═══╝  ██╔██╗ %b\n" "$CYAN" "$NC"
    printf "%b  ╚██████╗██║     ██╔╝ ██╗%b\n" "$CYAN" "$NC"
    printf "%b   ╚═════╝╚═╝     ╚═╝  ╚═╝%b\n" "$CYAN" "$NC"
    printf "\n"
    printf "  %bC++ Project Generator - Cpx Your Code!%b\n" "$YELLOW" "$NC"
    printf "\n"
}

detect_os() {
    OS="$(uname -s)"
    case "$OS" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)       echo "unknown" ;;
    esac
}

detect_arch() {
    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)            echo "unknown" ;;
    esac
}

check_dependencies() {
    if ! command -v curl > /dev/null 2>&1; then
        if ! command -v wget > /dev/null 2>&1; then
            printf "%bError: curl or wget is required%b\n" "$RED" "$NC"
            exit 1
        fi
        DOWNLOADER="wget"
    else
        DOWNLOADER="curl"
    fi
}

get_latest_version() {
    if [ "$DOWNLOADER" = "curl" ]; then
        VERSION=$(curl -sI "https://github.com/$REPO/releases/latest" | grep -i "location:" | sed 's/.*tag\///' | tr -d '\r\n')
    else
        VERSION=$(wget -qO- --server-response "https://github.com/$REPO/releases/latest" 2>&1 | grep -i "location:" | sed 's/.*tag\///' | tr -d '\r\n')
    fi
    
    if [ -z "$VERSION" ]; then
        VERSION="v1.0.0"
    fi
    echo "$VERSION"
}

download_binary() {
    OS=$1
    ARCH=$2
    VERSION=$3
    
    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        printf "%bError: Unsupported platform: %s/%s%b\n" "$RED" "$OS" "$ARCH" "$NC"
        exit 1
    fi
    
    BINARY_NAME_PLATFORM="$BINARY_NAME-$OS-$ARCH"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME_PLATFORM="${BINARY_NAME_PLATFORM}.exe"
    fi
    
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY_NAME_PLATFORM"
    
    printf "%bDownloading %s from %s...%b\n" "$CYAN" "$BINARY_NAME_PLATFORM" "$DOWNLOAD_URL" "$NC"
    
    INSTALL_DIR=$(get_install_dir)
    TARGET_PATH="$INSTALL_DIR/$BINARY_NAME"
    
    # Create install directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"
    
    if [ "$DOWNLOADER" = "curl" ]; then
        if ! curl -fSL "$DOWNLOAD_URL" -o "$TARGET_PATH"; then
            printf "%bError: Failed to download binary%b\n" "$RED" "$NC"
            exit 1
        fi
    else
        if ! wget -q "$DOWNLOAD_URL" -O "$TARGET_PATH"; then
            printf "%bError: Failed to download binary%b\n" "$RED" "$NC"
            exit 1
        fi
    fi
    
    chmod +x "$TARGET_PATH"
    
    printf "%bSuccessfully installed %s to %s%b\n" "$GREEN" "$BINARY_NAME" "$TARGET_PATH" "$NC"
    
    # Check if binary is in PATH
    if ! command -v "$BINARY_NAME" > /dev/null 2>&1; then
        printf "%bWarning: %s is not in your PATH.%b\n" "$YELLOW" "$BINARY_NAME" "$NC"
        printf "Add this to your shell profile (.bashrc, .zshrc, etc.):\n"
        printf "  export PATH=\"\$PATH:%s\"\n" "$INSTALL_DIR"
    else
        printf "%b%s is ready to use!%b\n" "$GREEN" "$BINARY_NAME" "$NC"
        printf "Run '%b%s version%b' to verify installation.\n" "$CYAN" "$BINARY_NAME" "$NC"
    fi
}

main() {
    print_banner
    
    check_dependencies
    
    OS=$(detect_os)
    ARCH=$(detect_arch)
    VERSION=$(get_latest_version)
    
    printf "%bDetected: %s/%s%b\n" "$CYAN" "$OS" "$ARCH" "$NC"
    printf "%bLatest version: %s%b\n" "$CYAN" "$VERSION" "$NC"
    printf "\n"
    
    download_binary "$OS" "$ARCH" "$VERSION"
}

main
