#!/usr/bin/env bash

function seal::cgo::musl_libc::install() {
  local build_os
  build_os="$(seal::util::get_raw_os)"
  if [[ "${build_os}" == "darwin" ]]; then
    seal::util::warn "please run 'xcode-select --install' to install xcode compiler,"
    seal::util::warn "or run 'softwareupdate --all --install --force' to update xcode software,"
    seal::util::warn "finally, run 'brew install filosottile/musl-cross/musl-cross --with-aarch64' to install musl cross build chain"
    return 1
  fi

  if [[ "${build_os}" != "linux" ]]; then
    seal::util::error "cannot execute cgo build via musl gcc without linux platform"
    return 1
  fi
  local platform
  platform="$(go env GOOS)/$(go env GOARCH)"
  local src
  case "${platform}" in
  linux/amd64)
    src="/musl-cross/x86_64-linux-musl-cross.tgz"
    ;;
  linux/arm64)
    src="/musl-cross/aarch64-linux-musl-cross.tgz"
    ;;
  windows/amd64)
    src="/musl-cross/x86_64-w64-mingw32-cross.tgz"
    ;;
  *)
    return 1
    ;;
  esac
  local directory="${ROOT_DIR}/.sbin/musl-${platform////-}"
  mkdir -p "${directory}"
  docker run \
    --rm \
    --pull=always \
    --interactive \
    --volume /tmp:/host \
    sealio/meta:hamal cp -f "${src}" "/host/$(basename "${src}")"
  tar -xvf "/tmp/$(basename "${src}")" \
    --directory "${directory}" \
    --strip-components 1 \
    --no-same-owner \
    --exclude ./*/usr
}

function seal::cgo::musl_libc::validate() {
  # shellcheck disable=SC2046
  if seal::cgo::musl_libc::bin; then
    return 0
  fi

  seal::log::info "installing musl libc"
  if seal::cgo::musl_libc::install; then
    return 0
  fi
  seal::log::error "no musl libc available"
  return 1
}

function seal::cgo::musl_libc::bin() {
  local platform
  platform="$(go env GOOS)/$(go env GOARCH)"

  PATH="${PATH}:${ROOT_DIR}/.sbin/musl-${platform////-}/bin"
  case "${platform}" in
  linux/amd64)
    if [[ -n "$(command -v x86_64-linux-musl-gcc)" ]] &&
      [[ -n "$(command -v x86_64-linux-musl-g++)" ]] &&
      [[ -n "$(command -v x86_64-linux-musl-ar)" ]]; then
      return 0
    fi
    ;;
  linux/arm64)
    if [[ -n "$(command -v aarch64-linux-musl-gcc)" ]] &&
      [[ -n "$(command -v aarch64-linux-musl-g++)" ]] &&
      [[ -n "$(command -v aarch64-linux-musl-ar)" ]]; then
      return 0
    fi
    ;;
  windows/amd64)
    if [[ -n "$(command -v x86_64-w64-mingw32-gcc)" ]] &&
      [[ -n "$(command -v x86_64-w64-mingw32-g++)" ]] &&
      [[ -n "$(command -v x86_64-w64-mingw32-ar)" ]]; then
      return 0
    fi
    ;;
  esac

  return 1
}

function seal::cgo::musl_libc::build() {
  local platform
  platform="$(go env GOOS)/$(go env GOARCH)"

  if seal::cgo::musl_libc::validate; then
    case "${platform}" in
    linux/amd64)
      CC="x86_64-linux-musl-gcc" CXX="x86_64-linux-musl-g++" AR="x86_64-linux-musl-ar" CGO_ENABLED=1 \
        go build "$@"
      return 0
      ;;
    linux/arm64)
      CC="aarch64-linux-musl-gcc" CXX="aarch64-linux-musl-g++" AR="aarch64-linux-musl-ar" CGO_ENABLED=1 \
        go build "$@"
      return 0
      ;;
    windows/amd64)
      CC="x86_64-w64-mingw32-gcc" CXX="x86_64-w64-mingw32-g++" AR="x86_64-w64-mingw32-ar" CGO_ENABLED=1 \
        go build "$@"
      return 0
      ;;
    esac
  fi

  seal::log::fatal "cannot execute cgo build ${platform} via musl libc"
}
