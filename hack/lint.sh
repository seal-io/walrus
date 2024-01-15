#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

function check_dirty() {
  [[ "${LINT_DIRTY:-false}" == "true" ]] || return 0

  if [[ -n "$(command -v git)" ]]; then
    if git_status=$(git status --porcelain 2>/dev/null) && [[ -n ${git_status} ]]; then
      seal::log::fatal "the git tree is dirty:\n$(git status --porcelain)"
    fi
  fi
}

function lint() {
  local path="$1"
  local path_ignored="$2"
  shift 2
  # shellcheck disable=SC2206
  local build_tags=(${*})

  [[ "${path}" == "${ROOT_DIR}" ]] || pushd "${path}" >/dev/null 2>&1

  seal::format::run "${path}" "${path_ignored}"
  if [[ ${#build_tags[@]} -gt 0 ]]; then
    GOLANGCI_LINT_CACHE="$(go env GOCACHE)/golangci-lint" seal::lint::run --build-tags="\"${build_tags[*]}\"" "${path}/..."
  else
    GOLANGCI_LINT_CACHE="$(go env GOCACHE)/golangci-lint" seal::lint::run "${path}/..."
  fi

  [[ "${path}" == "${ROOT_DIR}" ]] || popd >/dev/null 2>&1
}

function dispatch() {
  local target="$1"
  local path="$2"
  local path_ignored="$3"

  shift 3
  local specified_targets="$*"
  if [[ -n ${specified_targets} ]] && [[ ! ${specified_targets} =~ ${target} ]]; then
    return
  fi

  seal::log::debug "linting ${target}"
  if [[ "${PARALLELIZE:-false}" != "true" ]]; then
    lint "${path}" "${path_ignored}" "$(seal::target::build_tags "${target}")"
  else
    lint "${path}" "${path_ignored}" "$(seal::target::build_tags "${target}")" &
  fi
}

function after() {
  check_dirty
}

#
# main
#

seal::log::info "+++ LINT +++"

seal::commit::lint "${ROOT_DIR}"

dispatch "utils" "${ROOT_DIR}/staging/utils" "" "$@"
dispatch "walrus" "${ROOT_DIR}" "pkg/dao/model pkg/i18n" "$@"

after

if [[ "${PARALLELIZE:-false}" == "true" ]]; then
  seal::util::wait_jobs || seal::log::fatal "--- LINT ---"
fi
seal::log::info "--- LINT ---"
