#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

TEST_DIR="${ROOT_DIR}/.dist/test"
mkdir -p "${TEST_DIR}"

#
# util
#

function test() {
  local target="$1"
  local path="$2"
  local path_ignored="$3"
  shift 3

  [[ "${path}" == "${ROOT_DIR}" ]] || pushd "${path}" >/dev/null 2>&1

  local tags=()
  IFS=" " read -r -a tags <<<"$(seal::target::build_tags "${target}")"

  local ignored_path=()
  IFS=" " read -r -a ignored_path <<<"${path_ignored}"

  if [[ ${#ignored_path[@]} -gt 0 ]]; then
    CGO_ENABLED=1 go list ./... | grep -v -E "$(seal::util::join_array "|" "${ignored_path[@]}")" | xargs -I {} \
      go test \
      -v \
      -failfast \
      -race \
      -cover \
      -timeout=30m \
      -tags="${tags[*]}" \
      -coverprofile="${TEST_DIR}/${target}-coverage.out" \
      {}"/..."
  else
    CGO_ENABLED=1 go test \
      -v \
      -failfast \
      -race \
      -cover \
      -timeout=30m \
      -tags="${tags[*]}" \
      -coverprofile="${TEST_DIR}/${target}-coverage.out" \
      ./...
  fi

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

  seal::log::debug "testing ${target}"
  if [[ "${PARALLELIZE:-true}" == "false" ]]; then
    test "${target}" "${path}" "${path_ignored}"
  else
    test "${target}" "${path}" "${path_ignored}" &
  fi
}

#
# main
#

seal::log::info "+++ TEST +++"

dispatch "utils" "${ROOT_DIR}/staging/github.com/seal-io/utils" "" "$@"
dispatch "code-generator" "${ROOT_DIR}/staging/github.com/seal-io/code-generator" "" "$@"
dispatch "walrus" "${ROOT_DIR}" "staging pkg/clients" "$@"

if [[ "${PARALLELIZE:-true}" == "true" ]]; then
  seal::util::wait_jobs || seal::log::fatal "--- TEST ---"
fi
seal::log::info "--- TEST ---"
