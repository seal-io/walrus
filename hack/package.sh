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

function download_ui() {
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

function download_walrus_file_hub() {
    local path="${1}"

    mkdir -p "${path}"
    if ! git clone --depth 1 https://github.com/seal-io/walrus-file-hub "${path}"; then
      seal::log::fatal "failed to download walrus-file-hub"
    fi
}

function setup_image_package() {
  if [[ "${PACKAGE_PUSH:-false}" == "true" ]]; then
    seal::image::login
  fi
}

function setup_image_package_context() {
  local target="$1"
  local task="$2"
  local path="$3"
  local tag="$(seal::image::tag)"

  local context="${PACKAGE_DIR}/${target}/${task}"
  # create targeted dist
  rm -rf "${context}"
  mkdir -p "${context}"
  # copy targeted source to dist
  cp -rf "${path}/image" "${context}/image"
  # copy built result to dist
  cp -rf "${ROOT_DIR}/.dist/build/${target}" "${context}/build"

  # NB(thxCode): mutate the package context below.
  case "${target}" in
  walrus)
    case "${task}" in
    server)
      download_ui "${context}/image/var/lib/walrus/ui" "${tag}"
      download_walrus_file_hub "${context}/image/var/lib/walrus/walrus-file-hub"
      ;;
    esac
    ;;
  esac

  echo -n "${context}"
}

function package() {
  local target="$1"
  local task="$2"
  local path="$3"

  # shellcheck disable=SC2155
  local tag="${REPO:-sealio}/${target}:$(seal::image::tag)"
  # shellcheck disable=SC2155
  local platform="$(seal::target::package_platform)"

  # shellcheck disable=SC2155
  local context="$(setup_image_package_context "${target}" "${task}" "${path}")"

  if [[ "${PACKAGE_BUILD:-true}" == "true" ]]; then
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

  else
    local build_dir="${PACKAGE_DIR}/${target}/${task}/build"
    local release_dir="${PACKAGE_DIR}/${target}/${task}/release"
    mkdir -p "${release_dir}"

    # copy walrus-images.txt to release dir and update image tag
    sed "s/docker.io\/sealio\/walrus:.*$/docker.io\/sealio\/walrus:$(seal::image::tag)/g" "${ROOT_DIR}/hack/mirror/walrus-images.txt" > "${release_dir}/walrus-images.txt"

    # rename cli to walrus-cli and move to cli dir
    for file in "${build_dir}/cli"*; do
      mv "$file" "$(echo "$file" | sed 's/cli/walrus-cli/')";
    done
    mkdir -p "${build_dir}/cli"
    mv "${build_dir}/walrus-cli"* "${build_dir}/cli"

    # copy assets to release dir
    cp "${build_dir}/cli/walrus-cli"* "${release_dir}"
    cp "${ROOT_DIR}/hack/mirror/walrus-load-images.sh" "${release_dir}"
    cp "${ROOT_DIR}/hack/mirror/walrus-save-images.sh" "${release_dir}"

    # create sha256sums.txt
    shasum -a 256 "${release_dir}"/* | sed -e "s#${release_dir}/##g" > "${release_dir}/sha256sums.txt"
  fi
}

function before() {
  setup_image_package
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
  IFS=" " read -r -a tasks <<<"$(seal::util::find_subdirs ${path}/pack)"

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
