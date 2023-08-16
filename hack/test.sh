#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

TEST_DIR="${ROOT_DIR}/.dist/test"
mkdir -p "${TEST_DIR}"

function test() {
  local target="$1"
  local path="$2"

  [[ "${path}" == "${ROOT_DIR}" ]] || pushd "${path}" >/dev/null 2>&1

  local tags=()
  # shellcheck disable=SC2086
  IFS=" " read -r -a tags <<<"$(seal::target::build_tags ${target})"

  CGO_ENABLED=1 go test \
    -v \
    -failfast \
    -race \
    -cover \
    -timeout=10m \
    -tags="${tags[*]}" \
    -coverprofile="${TEST_DIR}/${target}-coverage.out" \
    "${path}/..."

  [[ "${path}" == "${ROOT_DIR}" ]] || popd >/dev/null 2>&1
}

function dispatch() {
  local target="$1"
  local path="$2"

  shift 2
  local specified_targets="$*"
  if [[ -n ${specified_targets} ]] && [[ ! ${specified_targets} =~ ${target} ]]; then
    return
  fi

  seal::log::debug "testing ${target}"
  if [[ "${PARALLELIZE:-true}" == "false" ]]; then
    test "${target}" "${path}"
  else
    test "${target}" "${path}" &
  fi
}

#
# main
#

seal::log::info "+++ TEST +++"

dispatch "utils" "${ROOT_DIR}/staging/utils" "$@"
dispatch "walrus" "${ROOT_DIR}" "$@"

if [[ "${PARALLELIZE:-true}" == "true" ]]; then
  seal::util::wait_jobs || seal::log::fatal "--- TEST ---"
fi
seal::log::info "--- TEST ---"
