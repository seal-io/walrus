#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

#
# util
#

function generate() {
  local target="$1"
  local task="$2"
  local path="$3"
  shift 3

  PATH="${ROOT_DIR}/.sbin:${ROOT_DIR}/.sbin/protoc/bin:${PATH}" \
    go run -mod=mod "${path}" "$@"
}

#
# lifecycle
#

function before() {
  if ! seal::protoc::protoc::validate; then
    seal::log::fatal "protoc hasn't installed"
  fi

  # NB(thxCode): install protoc-gen-gogo,
  # and trace the issue https://github.com/kubernetes/kubernetes/issues/96564.
  if [[ ! -f "${ROOT_DIR}/.sbin/protoc/bin/protoc-gen-gogo" ]]; then
    GOBIN="${ROOT_DIR}/.sbin/protoc/bin" \
      go install -mod=mod k8s.io/code-generator/cmd/go-to-protobuf/protoc-gen-gogo@latest
  fi

  if [[ ! -f "${ROOT_DIR}/.sbin/goimports" ]]; then
    GOBIN="${ROOT_DIR}/.sbin/protoc/bin" \
      go install golang.org/x/tools/cmd/goimports@latest
  fi
}

function dispatch() {
  local target="$1"
  local path="$2"
  shift 2

  local tasks=()
  IFS=" " read -r -a tasks <<<"$(seal::util::find_subdirs "${path}/gen")"

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

before

dispatch "walrus" "${ROOT_DIR}" "$@"

if [[ "${PARALLELIZE:-true}" == "true" ]]; then
  seal::util::wait_jobs || seal::log::fatal "--- GENERATE ---"
fi
seal::log::info "--- GENERATE ---"
