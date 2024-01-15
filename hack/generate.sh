#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

function generate() {
  local target="$1"
  local task="$2"
  local path="$3"

  go run -mod=mod "${path}"

  # FIXME(thxCode): remove this after bumped entc version.
  if [[ "${task}" == "entc" ]]; then
    local gofmt_opts=(
      "-w"
      "-r"
      "interface{} -> any"
      "${ROOT_DIR}/pkg/dao/model"
    )
    gofmt "${gofmt_opts[@]}"
  fi
}

function dispatch() {
  local target="$1"
  local path="$2"

  shift 2
  local specified_targets="$*"
  if [[ -n ${specified_targets} ]] && [[ ! ${specified_targets} =~ ${target} ]]; then
    return
  fi

  local tasks=()
  # shellcheck disable=SC2086
  IFS=" " read -r -a tasks <<<"$(seal::util::find_subdirs ${path}/gen)"

  for task in "${tasks[@]}"; do
    seal::log::debug "generating ${target} ${task}"
    if [[ "${PARALLELIZE:-true}" == "false" ]]; then
      generate "${target}" "${task}" "${path}/gen/${task}"
    else
      generate "${target}" "${task}" "${path}/gen/${task}" &
    fi
  done
}

#
# main
#

seal::log::info "+++ GENERATE +++"

dispatch "walrus" "${ROOT_DIR}" "$@"

if [[ "${PARALLELIZE:-true}" == "true" ]]; then
  seal::util::wait_jobs || seal::log::fatal "--- GENERATE ---"
fi
seal::log::info "--- GENERATE ---"
