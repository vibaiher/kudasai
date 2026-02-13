#!/bin/bash
set -e

# Build binary to temp dir
tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

go build -o "$tmpdir/kudasai" main.go

# Run clitest from temp dir so .kudasai.json writes go there
(cd "$tmpdir" && PATH="$tmpdir:$PATH" clitest "$OLDPWD"/examples/**)
