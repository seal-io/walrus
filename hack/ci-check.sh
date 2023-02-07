#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
pushd "${ROOT_DIR}" >/dev/null 2>&1

make deps "$@"
make lint "$@"
make test "$@"

popd >/dev/null 2>&1
