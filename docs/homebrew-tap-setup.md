# Homebrew Tap Setup Guide

This guide explains how to set up and maintain the Homebrew tap for GVM.

## Overview

The Homebrew tap allows users to install GVM via Homebrew:

```bash
brew tap ibreez3/gvm
brew install gvm
```

## Repository Structure

The tap files are located in the `homebrew-tap/` directory:

```
homebrew-tap/
├── Formula/
│   └── gvm.rb          # The Homebrew formula
├── LICENSE             # MIT License
└── README.md           # Tap documentation
```

## Initial Setup

### Step 1: Create the Tap Repository

Run the setup script:

```bash
chmod +x scripts/setup-brew-tap.sh
./scripts/setup-brew-tap.sh
```

This script will:
1. Initialize a Git repository in `homebrew-tap/`
2. Guide you through creating the repository on GitHub
3. Set up the remote and push the initial commit

### Step 2: Manual Setup (Alternative)

If you prefer to set up manually:

```bash
cd homebrew-tap
git init
git add .
git commit -m "Initial commit: Add GVM formula"
git remote add origin git@github.com:ibreez3/homebrew-gvm.git
git branch -M main
git push -u origin main
```

The tap repository should be created at: **https://github.com/ibreez3/homebrew-gvm**

## Updating on New Release

When you release a new version of GVM, you need to update the formula:

### Automated Update

Run the update script:

```bash
# Run from gvm repository root
chmod +x scripts/update-brew-formula.sh
./scripts/update-brew-formula.sh v0.3.0
```

This script will:
1. Download the release archive
2. Calculate the SHA256 checksum
3. Update the formula with the new version and checksum
4. Commit and push the changes

### Manual Update

If you prefer to update manually:

1. Download the release archive and calculate SHA256:
   ```bash
   curl -fsSL https://github.com/ibreez3/gvm/archive/refs/tags/v0.3.0.tar.gz | shasum -a 256
   ```

2. Update `homebrew-tap/Formula/gvm.rb`:
   ```ruby
   url "https://github.com/ibreez3/gvm/archive/refs/tags/v0.3.0.tar.gz"
   sha256 "<calculated_sha256>"
   ```

3. Commit and push:
   ```bash
   cd homebrew-tap
   git add Formula/gvm.rb
   git commit -m "Update gvm to v0.3.0"
   git push
   ```

## Formula Reference

The formula (`homebrew-tap/Formula/gvm.rb`) contains:

```ruby
class Gvm < Formula
  desc "Go Version Manager - Manage multiple Go versions easily"
  homepage "https://github.com/ibreez3/gvm"
  url "https://github.com/ibreez3/gvm/archive/refs/tags/v{VERSION}.tar.gz"
  sha256 "{CHECKSUM}"
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
```

## User Experience

### Installation

Users can install GVM with two commands:

```bash
brew tap ibreez3/gvm
brew install gvm
```

### Upgrade

```bash
brew upgrade gvm
```

### Uninstall

```bash
brew uninstall gvm
brew untap ibreez3/gvm
```

## Maintenance

### CI/CD Integration (Optional)

You can automate formula updates by adding to your CI workflow:

```yaml
# .github/workflows/release.yml
- name: Update Homebrew formula
  run: |
    chmod +x scripts/update-brew-formula.sh
    ./scripts/update-brew-formula.sh ${{ github.ref_name }}
  env:
    GITHUB_TOKEN: ${{ secrets.GH_PAT }}
```

Note: This requires a GitHub PAT with push access to the homebrew-gvm repository.

## Troubleshooting

### Formula not found

If users get "formula not found" error:

1. Check if the tap exists: `brew tap-info ibreez3/gvm`
2. Re-tap: `brew untap ibreez3/gvm && brew tap ibreez3/gvm`
3. Check the repository: https://github.com/ibreez3/homebrew-gvm

### SHA256 mismatch

If installation fails with SHA256 error:

1. Verify the release exists: https://github.com/ibreez3/gvm/releases
2. Re-run the update script to recalculate checksum

### Installation fails

If installation fails:

1. Check if Go is installed: `brew install go`
2. Check build logs: `brew install gvm --verbose`
3. Try installing from source for debugging
