#!/bin/bash
set -e

# Build binary to temp dir
bindir=$(mktemp -d)
trap 'rm -rf "$bindir"' EXIT

go build -o "$bindir/kudasai" main.go

# Run each test file in its own temp dir to avoid contamination
for testfile in examples/*.txt; do
  tmpdir=$(mktemp -d)
  (cd "$tmpdir" && PATH="$bindir:$PATH" clitest "$OLDPWD/$testfile")
  rm -rf "$tmpdir"
done
