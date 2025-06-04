# Documentation: https://docs.brew.sh/Formula-Cookbook
#                https://rubydoc.brew.sh/Formula
class Tada < Formula
  desc "A todo app to manage your todos from your terminal 😉"
  homepage "https://github.com/negadras/tada"
  license "MIT"
  url "ssh://git@github.com:negadras/tada.git", :using => :git, :tag => "0.0.1" # version marker, do not remove
  head "https://github.com/negadras/tada.git"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X github.com/negadras/tada/cmd.Version=#{version}")
    generate_completions_from_executable(bin/"tada", "completion")
  end

  test do
    system bin/"tada", "--version"
  end
end
