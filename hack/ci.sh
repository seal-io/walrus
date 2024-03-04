#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
pushd "${ROOT_DIR}" >/dev/null 2>&1

# check phase
if [[ "${CI_CHECK:-true}" == "true" ]]; then
  make deps "$@"
  make lint "$@"
  make test "$@"
fi

# publish phase
if [[ "${CI_PUBLISH:-true}" == "true" ]]; then
  make build "$@"
  make package "$@"
fi

popd >/dev/null 2>&1
