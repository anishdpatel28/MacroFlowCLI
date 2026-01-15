# This is a template for the Homebrew formula
# To publish: Create a tap repository at github.com/anishdpatel28/homebrew-macroflow
# Then add this formula as Formula/macroflow.rb

class Macroflow < Formula
  desc "Directory-scoped command aliases for developers"
  homepage "https://github.com/anishdpatel28/macroflow"
  url "https://github.com/anishdpatel28/macroflow/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "YOUR_SHA256_HERE"
  license "MIT"

  depends_on "go" => :build

  def install
    cd "MacroFlowCLI" do
      system "go", "build", *std_go_args(output: bin/"macro"), "main.go"
    end
  end

  test do
    system "#{bin}/macro", "--help"
  end
end

# Steps to publish on Homebrew:
#
# 1. Create a GitHub repository: homebrew-macroflow
#    URL: https://github.com/anishdpatel28/homebrew-macroflow
#
# 2. Add this formula to: Formula/macroflow.rb
#
# 3. Tag and release your MacroFlow project:
#    git tag v1.0.0
#    git push origin v1.0.0
#
# 4. Update the URL and sha256 in this formula:
#    - Download the release tarball
#    - Run: shasum -a 256 v1.0.0.tar.gz
#    - Update the sha256 value above
#
# 5. Users can then install with:
#    brew tap anishdpatel28/macroflow
#    brew install macroflow
