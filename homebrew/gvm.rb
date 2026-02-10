# Homebrew formula for gvm
# To install: brew install ibreez3/gvm/gvm

class Gvm < Formula
  desc "Go Version Manager - Manage multiple Go versions easily"
  homepage "https://github.com/ibreez3/gvm"
  url "https://github.com/ibreez3/gvm/archive/refs/tags/v0.2.0.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
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
