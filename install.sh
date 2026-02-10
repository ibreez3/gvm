#!/bin/bash

# GVM (Go Version Manager) Installer
# This script downloads and installs the latest version of gvm

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# GitHub repository
REPO="ibreez3/gvm"

# Print colored message
print_msg() {
    local color=$1
    local msg=$2
    echo -e "${color}${msg}${NC}"
}

# Detect OS and architecture
detect_os_arch() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)

    case $os in
        darwin)
            OS="Darwin"
            ;;
        linux)
            OS="Linux"
            ;;
        msys*|mingw*|cygwin*)
            print_msg "$RED" "Windows is not supported by this installer. Please download from https://github.com/${REPO}/releases"
            exit 1
            ;;
        *)
            print_msg "$RED" "Unsupported OS: $os"
            exit 1
            ;;
    esac

    case $arch in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            print_msg "$RED" "Unsupported architecture: $arch"
            exit 1
            ;;
    esac

    print_msg "$GREEN" "Detected OS: $OS, Architecture: $ARCH"
}

# Get latest version from GitHub
get_latest_version() {
    print_msg "$YELLOW" "Fetching latest version from GitHub..."
    VERSION=$(curl -s https://api.github.com/repos/${REPO}/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/v//')

    if [ -z "$VERSION" ]; then
        print_msg "$RED" "Failed to fetch latest version"
        exit 1
    fi

    print_msg "$GREEN" "Latest version: v${VERSION}"
}

# Download and install gvm
install_gvm() {
    local filename="gvm_${OS}_${ARCH}.tar.gz"
    local download_url="https://github.com/${REPO}/releases/download/v${VERSION}/${filename}"
    local tmp_dir=$(mktemp -d)

    print_msg "$YELLOW" "Downloading gvm v${VERSION}..."
    if ! curl -fsSL "$download_url" -o "${tmp_dir}/${filename}"; then
        print_msg "$RED" "Failed to download gvm"
        rm -rf "$tmp_dir"
        exit 1
    fi

    print_msg "$YELLOW" "Extracting archive..."
    cd "$tmp_dir"
    tar -xzf "$filename"

    # Determine install directory
    INSTALL_DIR="/usr/local/bin"
    if [ ! -w "$INSTALL_DIR" ] && [ -n "$HOME" ]; then
        # Try to use ~/.local/bin or create it
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    fi

    # Check if we can write to install directory
    if [ ! -w "$INSTALL_DIR" ]; then
        print_msg "$YELLOW" "Requesting sudo permissions to install to $INSTALL_DIR..."
        if ! sudo mkdir -p "$INSTALL_DIR" 2>/dev/null; then
            print_msg "$RED" "Failed to create install directory"
            rm -rf "$tmp_dir"
            exit 1
        fi
        SUDO="sudo"
    else
        SUDO=""
    fi

    print_msg "$YELLOW" "Installing gvm to $INSTALL_DIR..."
    $SUDO cp gvm "$INSTALL_DIR/"
    $SUDO chmod +x "$INSTALL_DIR/gvm"

    # Cleanup
    cd -
    rm -rf "$tmp_dir"

    print_msg "$GREEN" "gvm v${VERSION} installed successfully to $INSTALL_DIR/gvm"
}

# Verify installation
verify_installation() {
    if command -v gvm &> /dev/null; then
        print_msg "$GREEN" "Installation verified! gvm is now available."
        gvm version || true
    else
        print_msg "$YELLOW" "gvm is installed but not in PATH."
        print_msg "$YELLOW" "Please add the following to your shell profile (~/.bashrc or ~/.zshrc):"
        if [ "$INSTALL_DIR" = "$HOME/.local/bin" ]; then
            print_msg "$YELLOW" "  export PATH=\"\$HOME/.local/bin:\$PATH\""
        fi
    fi
}

# Print post-installation instructions
print_instructions() {
    echo ""
    print_msg "$GREEN" "========================================="
    print_msg "$GREEN" "Next steps:"
    print_msg "$GREEN" "========================================="
    echo ""
    echo "  1. Initialize gvm:"
    echo "     $ gvm init"
    echo ""
    echo "  2. Reload your shell or run:"
    echo "     $ source ~/.zshrc  # or source ~/.bashrc"
    echo ""
    echo "  3. Install a Go version:"
    echo "     $ gvm install 1.22.5"
    echo ""
    echo "  4. For more information, run:"
    echo "     $ gvm --help"
    echo ""
}

# Main installation flow
main() {
    print_msg "$GREEN" "========================================="
    print_msg "$GREEN" "GVM (Go Version Manager) Installer"
    print_msg "$GREEN" "========================================="
    echo ""

    detect_os_arch
    get_latest_version
    install_gvm
    verify_installation
    print_instructions
}

main
