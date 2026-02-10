#!/bin/bash

# Formula Update Script for GVM Homebrew Tap
# This script updates the formula when a new gvm version is released

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

GITHUB_USER="ibreez3"
REPO="gvm"
TAP_DIR="$(pwd)/homebrew-tap"
FORMULA_FILE="$TAP_DIR/Formula/gvm.rb"

print_msg() {
    local color=$1
    local msg=$2
    echo -e "${color}${msg}${NC}"
}

print_msg "$GREEN" "========================================"
print_msg "$GREEN" "GVM Homebrew Formula Updater"
print_msg "$GREEN" "========================================"
echo ""

# Get version from argument or prompt
if [ -n "$1" ]; then
    VERSION="$1"
else
    read -p "Enter GVM version (e.g., 0.2.0): " VERSION
fi

# Remove 'v' prefix if present
VERSION=${VERSION#v}

print_msg "$YELLOW" "Updating formula for version v${VERSION}..."

# Check if tap directory exists
if [ ! -d "$TAP_DIR" ]; then
    print_msg "$RED" "Error: homebrew-tap directory not found"
    print_msg "$YELLOW" "Please run this script from the gvm repository root"
    exit 1
fi

# Download the release archive to get the checksum
ARCHIVE_URL="https://github.com/${GITHUB_USER}/${REPO}/archive/refs/tags/v${VERSION}.tar.gz"
TMP_FILE=$(mktemp)

print_msg "$YELLOW" "Downloading release archive to calculate SHA256..."
if ! curl -fsSL "$ARCHIVE_URL" -o "$TMP_FILE"; then
    print_msg "$RED" "Failed to download release archive"
    print_msg "$YELLOW" "Please check that version v${VERSION} exists at:"
    print_msg "$YELLOW" "  https://github.com/${GITHUB_USER}/${REPO}/releases"
    rm -f "$TMP_FILE"
    exit 1
fi

# Calculate SHA256
SHA256=$(shasum -a 256 "$TMP_FILE" | awk '{print $1}')
rm -f "$TMP_FILE"

print_msg "$GREEN" "SHA256: $SHA256"

# Update the formula file
print_msg "$YELLOW" "Updating formula file..."

cat > "$FORMULA_FILE" << EOF
# Homebrew formula for gvm
# Repository: https://github.com/${GITHUB_USER}/homebrew-gvm
# To install: brew tap ${GITHUB_USER}/gvm && brew install gvm

class Gvm < Formula
  desc "Go Version Manager - Manage multiple Go versions easily"
  homepage "https://github.com/${GITHUB_USER}/${REPO}"
  url "https://github.com/${GITHUB_USER}/${REPO}/archive/refs/tags/v${VERSION}.tar.gz"
  sha256 "$SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    system bin/"gvm", "--version"
    system bin/"gvm", "--help"
  end
end
EOF

print_msg "$GREEN" "Formula updated successfully!"

# Commit and push
cd "$TAP_DIR"

print_msg "$YELLOW" "Committing changes..."
git add Formula/gvm.rb
git commit -m "Update gvm to v${VERSION}"

print_msg "$YELLOW" "Pushing to GitHub..."
git push

echo ""
print_msg "$GREEN" "========================================"
print_msg "$GREEN" "Update Complete!"
print_msg "$GREEN" "========================================"
echo ""
print_msg "$GREEN" "Formula updated and pushed to:"
print_msg "$GREEN" "  https://github.com/${GITHUB_USER}/homebrew-gvm"
echo ""
print_msg "$YELLOW" "Users can upgrade with:"
print_msg "$YELLOW" "  brew upgrade gvm"
echo ""
