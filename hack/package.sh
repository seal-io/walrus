#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

PACKAGE_DIR="${ROOT_DIR}/.dist/package"
mkdir -p "${PACKAGE_DIR}"
PACKAGE_TMP_DIR="${PACKAGE_DIR}/tmp"
mkdir -p "${PACKAGE_TMP_DIR}"

#
# util
#

function get_ui() {
  local default_tag="latest"
  local path="${1}"
  local tag="${2}"

  mkdir -p "${PACKAGE_TMP_DIR}/ui"

  mkdir -p "${path}"
  if ! curl --retry 3 --retry-all-errors --retry-delay 3 -sSfL "https://walrus-ui-1303613262.cos.accelerate.myqcloud.com/releases/${tag}.tar.gz" 2>/dev/null |
    tar -xzf - --directory "${PACKAGE_TMP_DIR}/ui" 2>/dev/null; then

    if [[ "${tag:-}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$ ]]; then
      seal::log::fatal "failed to download '${tag}' ui archive"
    fi

    seal::log::warn "failed to download '${tag}' ui archive, fallback to '${default_tag}' ui archive"
    if ! curl --retry 3 --retry-all-errors --retry-delay 3 -sSfL "https://walrus-ui-1303613262.cos.accelerate.myqcloud.com/releases/${default_tag}.tar.gz" |
      tar -xzf - --directory "${PACKAGE_TMP_DIR}/ui" 2>/dev/null; then
      seal::log::fatal "failed to download '${default_tag}' ui archive"
    fi
  fi
  cp -a "${PACKAGE_TMP_DIR}/ui/dist/." "${path}"

  rm -rf "${PACKAGE_TMP_DIR}/ui"
}

function get_files() {
  local path="${1}"

  mkdir -p "${path}"
  if ! git clone --depth 1 https://github.com/seal-io/walrus-file-hub "${path}"; then
    seal::log::fatal "failed to download walrus-file-hub"
  fi
}

function pack_image() {
  local target="$1"
  local task="$2"
  local path="$3"
  local context="$4"
  shift 4

  mkdir -p "${context}/image"

  if [[ "${target}" == "walrus" ]]; then
    # get ui
    get_ui "${context}/image/var/lib/walrus/ui" "${TAG:-${GIT_VERSION}}"
    # get files
    get_files "${context}/image/var/lib/walrus/files"
  fi

  local name="${target}-${task}"
  if [[ "${target}" == "${task}" ]]; then
    name="${target}"
  fi
  local tag
  tag="${REPO:-sealio}/${name}:$(seal::image::tag)"
  local platform
  platform="$(seal::target::package_platform)"

  [[ "${PACKAGE_BUILD:-true}" == "true" ]] || return 0

  # shellcheck disable=SC2086
  local cache="type=registry,ref=sealio/build-cache:${target}-${task}"
  local output="type=image,push=${PACKAGE_PUSH:-false}"

  seal::image::build \
    --tag="${tag}" \
    --platform="${platform}" \
    --cache-from="${cache}" \
    --output="${output}" \
    --progress="plain" \
    --file="${context}/image/Dockerfile" \
    "${context}"
}

function pack_release() {
  local target="$1"
  local task="$2"
  local path="$3"
  local context="$4"
  shift 4

  mkdir -p "${context}/release"

  # copy shell script.
  cp -f "${ROOT_DIR}/hack/mirror/"walrus-* "${context}/release"
  # mutate walrus-images.txt.
  seal::util::sed "s/docker.io\/sealio\/walrus:.*$/docker.io\/sealio\/walrus:$(seal::image::tag)/g" "${context}/release/walrus-images.txt"
  # copy cli.
  cp -f "${context}/build/"walrus-cli-* "${context}/release"
  # create checksum.
  find "${context}/release" -type f -exec shasum -a 256 {} \; | grep -v -E "sha256sums" | sed -e "s#${context}/release/##g" >"${context}/release/sha256sums.txt"
}

function package() {
  local target="$1"
  local task="$2"
  local path="$3"
  shift 3

  local actions=()
  IFS=" " read -r -a actions <<<"$(seal::util::find_subdirs "${path}")"

  for action in "${actions[@]}"; do
    # prepare context.
    local context="${PACKAGE_DIR}/${target}/${task}"
    rm -rf "${context}"
    mkdir -p "${context}"

    # copy build result to "${context}/build/".
    cp -rf "${ROOT_DIR}/.dist/build/${target}" "${context}/build"

    case "${action}" in
    image)
      # copy image assets to "${context}/image/".
      cp -a "${path}/image" "${context}/image"
      pack_image "${target}" "${task}" "${path}" "${context}" "$@"
      ;;
    release)
      pack_release "${target}" "${task}" "${path}" "${context}" "$@"
      ;;
    esac
  done
}

#
# lifecycle
#

function before() {
  [[ "${PACKAGE_PUSH:-false}" == "true" ]] || return 0

  seal::image::login
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
  IFS=" " read -r -a tasks <<<"$(seal::util::find_subdirs "${path}"/pack)"

  for task in "${tasks[@]}"; do
    seal::log::debug "packaging ${target} ${task}"
    if [[ "${PARALLELIZE:-true}" == "false" ]]; then
      package "${target}" "${task}" "${path}/pack/${task}"
    else
      package "${target}" "${task}" "${path}/pack/${task}" &
    fi
  done
}

#
# main
#

seal::log::info "+++ PACKAGE +++" "tag: $(seal::image::tag)"

before

dispatch "walrus" "${ROOT_DIR}" "$@"

if [[ "${PARALLELIZE:-true}" == "true" ]]; then
  seal::util::wait_jobs || seal::log::fatal "--- PACKAGE ---"
fi
seal::log::info "--- PACKAGE ---"
