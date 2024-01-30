#!/usr/bin/env bash

readonly DEFAULT_BUILD_TAGS=(
  "netgo"
  "urfave_cli_no_docs"
  "jsoniter"
  "sqlite_json"
  "ginx"
)

function seal::target::build_tags() {
  local target="${1:-}"

  local tags
  if [[ -n "${BUILD_TAGS:-}" ]]; then
    IFS="," read -r -a tags <<<"${BUILD_TAGS}"
  else
    case "${target}" in
    utils)
      tags=()
      ;;
    *)
      tags=("${DEFAULT_BUILD_TAGS[@]}")
      ;;
    esac
  fi

  if [[ ${#tags[@]} -ne 0 ]]; then
    echo -n "${tags[@]}"
  fi
}

readonly DEFAULT_BUILD_PLATFORMS=(
  linux/amd64
  linux/arm64
)

readonly DEFAULT_BUILD_CLI_PLATFORMS=(
  linux/amd64
  linux/arm64
  darwin/amd64
  darwin/arm64
  windows/amd64
  windows/arm64
)

function seal::target::build_platforms() {
  local target="${1:-}"
  local task="${2:-}"

  local platforms
  if [[ -z "${OS:-}" ]] && [[ -z "${ARCH:-}" ]]; then
    if [[ -n "${BUILD_PLATFORMS:-}" ]]; then
      IFS="," read -r -a platforms <<<"${BUILD_PLATFORMS}"
    else
      case "${target}" in
      utils)
        platforms=()
        ;;
      *)
        case "${task}" in
        cli)
          platforms=("${DEFAULT_BUILD_CLI_PLATFORMS[@]}")
          ;;
        *)
          platforms=("${DEFAULT_BUILD_PLATFORMS[@]}")
          ;;
        esac
        ;;
      esac
    fi
  else
    local os="${OS:-$(seal::util::get_raw_os)}"
    local arch="${ARCH:-$(seal::util::get_raw_arch)}"
    platforms=("${os}/${arch}")
  fi

  if [[ ${#platforms[@]} -ne 0 ]]; then
    echo -n "${platforms[@]}"
  fi
}

function seal::target::package_platform() {
  if [[ -z "${OS:-}" ]] && [[ -z "${ARCH:-}" ]]; then
    echo -n "linux/$(seal::util::get_raw_arch)"
  else
    echo -n "${OS}/${ARCH}"
  fi
}
