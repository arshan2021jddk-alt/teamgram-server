#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-300}"
GO_TEST_FLAGS="${GO_TEST_FLAGS:--run TestDoesNotExist -count=1}"

packages=(
  "./app/service/status/internal/core"
  "./app/service/biz/user/internal/core"
  "./app/bff/users/internal/core"
  "./app/bff/usernames/internal/core"
  "./app/bff/messages/internal/core"
  "./app/bff/drafts/internal/core"
  "./app/bff/chats/internal/core"
  "./app/bff/chatinvites/internal/core"
  "./app/bff/dialogs/internal/core"
  "./app/bff/notification/internal/core"
  "./pkg/code/me"
)

failures=0

echo "== Build/compile check started =="
echo "timeout per package: ${TIMEOUT_SECONDS}s"
echo "go test flags: ${GO_TEST_FLAGS}"

for p in "${packages[@]}"; do
  echo
  echo "=== Checking ${p} ==="
  out_file="/tmp/build-check-$(echo "$p" | tr '/.' '__').log"

  if timeout "${TIMEOUT_SECONDS}s" go test "$p" ${GO_TEST_FLAGS} >"$out_file" 2>&1; then
    echo "PASS: ${p}"
  else
    ec=$?
    if [[ $ec -eq 124 ]]; then
      echo "TIMEOUT: ${p} (>${TIMEOUT_SECONDS}s)"
    else
      echo "FAIL: ${p} (exit=${ec})"
    fi
    echo "--- last output ---"
    tail -n 120 "$out_file" || true
    echo "-------------------"
    failures=$((failures + 1))
  fi
done

echo
if [[ $failures -gt 0 ]]; then
  echo "Build/compile check finished with ${failures} failing package(s)."
  exit 1
fi

echo "Build/compile check passed for all packages."
