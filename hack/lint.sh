#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

#
# util
#

function lint() {
  local target="$1"
  local path="$2"
  local path_ignored="$3"
  shift 3

  [[ "${path}" == "${ROOT_DIR}" ]] || pushd "${path}" >/dev/null 2>&1

  local tags=()
  IFS=" " read -r -a tags <<<"$(seal::target::build_tags "${target}")"

  local ignored_path=()
  IFS=" " read -r -a ignored_path <<<"${path_ignored}"

  local opts=()
  if [[ ${#tags[@]} -gt 0 ]]; then
    opts+=("--build-tags=\"${tags[*]}\"")
  fi
  for ig in "${ignored_path[@]}"; do
    opts+=("--skip-dirs=${ig}")
  done
  opts+=("${path}/...")
  GOLANGCI_LINT_CACHE="$(go env GOCACHE)/golangci-lint" seal::lint::run "${opts[@]}"

  [[ "${path}" == "${ROOT_DIR}" ]] || popd >/dev/null 2>&1
}

#
# lifecycle
#

function dispatch() {
  local target="$1"
  local path="$2"
  local path_ignored="$3"
  shift 3

  seal::log::debug "linting ${target}"
  if [[ "${PARALLELIZE:-false}" != "true" ]]; then
    lint "${target}" "${path}" "${path_ignored}"
  else
    lint "${target}" "${path}" "${path_ignored}" &
  fi
}

function after() {
  [[ "${LINT_DIRTY:-false}" == "true" ]] || return 0

  if [[ -n "$(command -v git)" ]]; then
    if git_status=$(git status --porcelain 2>/dev/null) && [[ -n ${git_status} ]]; then
      seal::log::fatal "the git tree is dirty:\n$(git status --porcelain)"
    fi
  fi
}

#
# main
#

seal::log::info "+++ LINT +++"

seal::commit::lint "${ROOT_DIR}"

dispatch "utils" "${ROOT_DIR}/staging/github.com/seal-io/utils" "" "$@"
dispatch "code-generator" "${ROOT_DIR}/staging/github.com/seal-io/code-generator" "" "$@"
dispatch "walrus" "${ROOT_DIR}" "staging pkg/clients" "$@"

after

if [[ "${PARALLELIZE:-false}" == "true" ]]; then
  seal::util::wait_jobs || seal::log::fatal "--- LINT ---"
fi
seal::log::info "--- LINT ---"
