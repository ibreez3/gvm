#!/bin/bash

# Homebrew Tap Setup Script
# This script helps create and publish the homebrew-gvm tap repository

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GITHUB_USER="ibreez3"
TAP_REPO="homebrew-gvm"
TAP_DIR="$(pwd)/homebrew-tap"

print_msg() {
    local color=$1
    local msg=$2
    echo -e "${color}${msg}${NC}"
}

print_msg "$GREEN" "========================================"
print_msg "$GREEN" "Homebrew Tap Setup for GVM"
print_msg "$GREEN" "========================================"
echo ""

# Check if homebrew-tap directory exists
if [ ! -d "$TAP_DIR" ]; then
    print_msg "$YELLOW" "Error: homebrew-tap directory not found in current directory"
    print_msg "$YELLOW" "Please run this script from the gvm repository root"
    exit 1
fi

cd "$TAP_DIR"

print_msg "$BLUE" "Step 1: Initializing Git repository..."
if [ ! -d ".git" ]; then
    git init
    git add .
    git commit -m "Initial commit: Add GVM formula"

    print_msg "$YELLOW" "Git repository initialized"
else
    print_msg "$YELLOW" "Git repository already exists"
fi

echo ""
print_msg "$BLUE" "Step 2: Creating GitHub repository..."
print_msg "$YELLOW" "Please create a new repository on GitHub:"
print_msg "$YELLOW" "  - Repository name: $TAP_REPO"
print_msg "$YELLOW" "  - Description: Homebrew tap for GVM (Go Version Manager)"
print_msg "$YELLOW" "  - Visibility: Public"
print_msg "$YELLOW" "  - DO NOT initialize with README, .gitignore, or license"
echo ""
read -p "Press Enter after creating the repository on GitHub..."

echo ""
print_msg "$BLUE" "Step 3: Setting up remote..."
if ! git remote get-url origin &>/dev/null; then
    git remote add origin "git@github.com:${GITHUB_USER}/${TAP_REPO}.git"
    print_msg "$YELLOW" "Remote 'origin' added"
else
    print_msg "$YELLOW" "Remote 'origin' already exists"
fi

echo ""
print_msg "$BLUE" "Step 4: Pushing to GitHub..."
read -p "Enter branch name (main/master) [main]: " branch
branch=${branch:-main}

git branch -M "$branch"
print_msg "$YELLOW" "Pushing to origin/$branch..."
git push -u origin "$branch"

echo ""
print_msg "$GREEN" "========================================"
print_msg "$GREEN" "Setup Complete!"
print_msg "$GREEN" "========================================"
echo ""
print_msg "$GREEN" "Your Homebrew tap is now live at:"
print_msg "$BLUE" "  https://github.com/${GITHUB_USER}/${TAP_REPO}"
echo ""
print_msg "$YELLOW" "Users can now install gvm via:"
print_msg "$BLUE" "  brew tap ${GITHUB_USER}/gvm"
print_msg "$BLUE" "  brew install gvm"
echo ""
print_msg "$YELLOW" "Updating the formula on release:"
print_msg "$YELLOW" "When you release a new version of gvm, update:"
print_msg "$BLUE" "  $TAP_DIR/Formula/gvm.rb"
print_msg "$YELLOW" "Update the 'url' and 'sha256' fields, then commit and push."
echo ""
