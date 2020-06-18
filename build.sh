#!/bin/sh -e

indent='s/^/    /'
no_fetch=""

while test $# -gt 0; do
  case "$1" in
    --no-fetch)   no_fetch="true";  shift 1 ;;
    *)                              break   ;;
  esac
done

export CGO_ENABLED=0
export GO111MODULE=on

git_commit=$(git rev-parse HEAD)
git_dirty=$(git status --porcelain)
git_dirty=${git_dirty:+"CHANGES"}

echo "==> Build"
echo "    Git Commit: $git_commit"
echo "    Git Status: ${git_dirty:-"Clean"}"

ldflags="-X brickrecon/version.GitCommit=$git_commit -X brickrecon/version.Prerelease=$git_dirty"

if [ -n "$no_fetch" ]; then
  echo "--> Downloading Modules"
  go mod download \
    | sed -e "$indent"
fi

echo "--> Building..."

go build -a -installsuffix cgo -ldflags "$ldflags" -o brickrecon  \
  | sed -e "$indent"

echo "==> Build Succeeded"
