#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

#
# util
#

function mod() {
  local path="$1"
  shift 1

  [[ "${path}" == "${ROOT_DIR}" ]] || pushd "${path}" >/dev/null 2>&1

  if [[ -n "$*" ]] && [[ "$*" =~ update$ ]]; then
    go get -u ./...
  fi

  go mod tidy
  go mod download

  [[ "${path}" == "${ROOT_DIR}" ]] || popd >/dev/null 2>&1
}

#
# lifecycle
#

function dispatch() {
  local target="$1"
  local path="$2"
  shift 2

  seal::log::debug "modding ${target}"
  if [[ "${PARALLELIZE:-true}" == "false" ]]; then
    mod "${path}" "$@"
  else
    mod "${path}" "$@" &
  fi
}

#
# main
#

seal::log::info "+++ MOD +++"

dispatch "utils" "${ROOT_DIR}/staging/github.com/seal-io/utils" "$@"
dispatch "code-generator" "${ROOT_DIR}/staging/github.com/seal-io/code-generator" "$@"
dispatch "walrus" "${ROOT_DIR}" "$@"

if [[ "${PARALLELIZE:-true}" == "true" ]]; then
  seal::util::wait_jobs || seal::log::fatal "--- MOD ---"
fi
seal::log::info "--- MOD ---"
