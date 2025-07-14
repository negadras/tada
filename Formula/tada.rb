class Tada < Formula
  desc "A tada cli app to manage your todos from your terminal 😉"
  homepage "https://github.com/negadras/tada"
  license "MIT"

  # Use HTTPS and point at tag
  url "https://github.com/negadras/tada.git",
      using: :git,
      tag:      "2.1.1"

  head "https://github.com/negadras/tada.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X github.com/negadras/tada/cmd.Version=#{version}")
    generate_completions_from_executable(bin/"tada", "completion")
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/tada --version")
  end
end
