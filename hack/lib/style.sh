#!/usr/bin/env bash

# -----------------------------------------------------------------------------
# Lint variables helpers. These functions need the
# following variables:
#
#    GOLANGCI_LINT_VERSION  -  The Golangci-lint version, default is v1.50.1.
#        COMMITSAR_VERSION  -  The Commitsar version, default is v0.20.1.

golangci_lint_version=${GOLANGCI_LINT_VERSION:-"v1.50.1"}
commitsar_version=${COMMITSAR_VERSION:-"v0.20.1"}

function seal::lint::golangci_lint::install() {
  curl --retry 3 --retry-all-errors --retry-delay 3 -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${ROOT_DIR}/.sbin" "${golangci_lint_version}"
}

function seal::lint::golangci_lint::validate() {
  # shellcheck disable=SC2046
  if [[ -n "$(command -v $(seal::lint::golangci_lint::bin))" ]]; then
    if [[ $($(seal::lint::golangci_lint::bin) --version 2>&1 | cut -d " " -f 4 2>&1 | head -n 1) == "${golangci_lint_version#v}" ]]; then
      return 0
    fi
  fi

  seal::log::info "installing golangci-lint ${golangci_lint_version}"
  if seal::lint::golangci_lint::install; then
    seal::log::info "golangci_lint $($(seal::lint::golangci_lint::bin) --version 2>&1 | cut -d " " -f 4 2>&1 | head -n 1)"
    return 0
  fi
  seal::log::error "no golangci-lint available"
  return 1
}

function seal::lint::golangci_lint::bin() {
  local bin="golangci-lint"
  if [[ -f "${ROOT_DIR}/.sbin/golangci-lint" ]]; then
    bin="${ROOT_DIR}/.sbin/golangci-lint"
  fi
  echo -n "${bin}"
}

function seal::lint::run() {
  if ! seal::lint::golangci_lint::validate; then
    seal::log::warn "using go fmt/vet instead golangci-lint"
    shift 1
    local fmt_args=()
    local vet_args=()
    for arg in "$@"; do
      if [[ "${arg}" == "--build-tags="* ]]; then
        arg="${arg//--build-/-}"
        vet_args+=("${arg}")
        continue
      fi
      fmt_args+=("${arg}")
      vet_args+=("${arg}")
    done
    seal::log::debug "go fmt ${fmt_args[*]}"
    go fmt "${fmt_args[@]}"
    seal::log::debug "go vet ${vet_args[*]}"
    go vet "${vet_args[@]}"
    return 0
  fi

  seal::log::debug "golangci-lint $*"
  $(seal::lint::golangci_lint::bin) run --fix "$@"
}

function seal::format::goimports::install() {
  GOBIN="${ROOT_DIR}/.sbin" go install golang.org/x/tools/cmd/goimports@latest
}

function seal::format::goimports::validate() {
  # shellcheck disable=SC2046
  if [[ -n "$(command -v $(seal::format::goimports::bin))" ]]; then
    return 0
  fi

  seal::log::info "installing goimports"
  if seal::format::goimports::install; then
    return 0
  fi
  seal::log::error "no goimports available"
  return 1
}

function seal::format::goimports::bin() {
  local bin="goimports"
  if [[ -f "${ROOT_DIR}/.sbin/goimports" ]]; then
    bin="${ROOT_DIR}/.sbin/goimports"
  fi
  echo -n "${bin}"
}

function seal::format::gofumpt::install() {
  GOBIN="${ROOT_DIR}/.sbin" go install mvdan.cc/gofumpt@latest
}

function seal::format::gofumpt::validate() {
  # shellcheck disable=SC2046
  if [[ -n "$(command -v $(seal::format::gofumpt::bin))" ]]; then
    return 0
  fi

  seal::log::info "installing gofumpt"
  if seal::format::gofumpt::install; then
    return 0
  fi
  seal::log::error "no gofumpt available"
  return 1
}

function seal::format::gofumpt::bin() {
  local bin="gofumpt"
  if [[ -f "${ROOT_DIR}/.sbin/gofumpt" ]]; then
    bin="${ROOT_DIR}/.sbin/gofumpt"
  fi
  echo -n "${bin}"
}

# install golines
function seal::format::golines::install() {
  GOBIN="${ROOT_DIR}/.sbin" go install github.com/segmentio/golines@latest
}

function seal::format::golines::validate() {
  # shellcheck disable=SC2046
  if [[ -n "$(command -v $(seal::format::golines::bin))" ]]; then
    return 0
  fi

  seal::log::info "installing golines"
  if seal::format::golines::install; then
    return 0
  fi
  seal::log::error "no golines available"
  return 1
}

function seal::format::golines::bin() {
  local bin="golines"
  if [[ -f "${ROOT_DIR}/.sbin/golines" ]]; then
    bin="${ROOT_DIR}/.sbin/golines"
  fi
  echo -n "${bin}"
}

# install wsl(Whitespace Linter)
function seal::format::wsl::install() {
  GOBIN="${ROOT_DIR}/.sbin" go install github.com/bombsimon/wsl/v4/cmd...@master
}

function seal::format::wsl::validate() {
  # shellcheck disable=SC2046
  if [[ -n "$(command -v $(seal::format::wsl::bin))" ]]; then
    return 0
  fi

  seal::log::info "installing wsl"
  if seal::format::wsl::install; then
    return 0
  fi
  seal::log::error "no wsl available"
  return 1
}

function seal::format::wsl::bin() {
  local bin="wsl"
  if [[ -f "${ROOT_DIR}/.sbin/wsl" ]]; then
    bin="${ROOT_DIR}/.sbin/wsl"
  fi
  echo -n "${bin}"
}

function seal::format::run() {
  if ! seal::format::goimports::validate; then
    seal::log::fatal "cannot execute goimports as client is not found"
  fi

  seal::log::debug "goimports $*"
  $(seal::format::goimports::bin) "$@"

  # format go files
  set +e
  local path=$4
  
  if ! seal::format::gofumpt::validate; then
    seal::log::fatal "cannot execute gofumpt as client is not found"
  fi

  seal::log::debug "gofumpt -extra -l -w  $path"
  $(seal::format::gofumpt::bin) -extra -l -w "$path"

  if ! seal::format::golines::validate; then
    seal::log::fatal "cannot execute golines as client is not found"
  fi

  local golines_opts="-w -m 120 --no-reformat-tags --base-formatter=$(seal::format::gofumpt::bin)"
  seal::log::debug "golines $golines_opts $path"
  $(seal::format::golines::bin) $golines_opts "$path"

  if ! seal::format::wsl::validate; then
    seal::log::fatal "cannot execute wsl as client is not found"
  fi

  local generated_pkg="pkg/dao/model"
  local wsl_opts="--allow-assign-and-anything --allow-trailing-comment --force-short-decl-cuddling=false --fix"
  seal::log::debug "wsl $wsl_opts $path/..."
  go list $path/... | grep -v $generated_pkg | xargs $(seal::format::wsl::bin) $wsl_opts

  set -e
}

function seal::commit::commitsar::install() {
  local os
  os="$(seal::util::get_raw_os)"
  local arch
  arch="$(seal::util::get_raw_arch)"
  curl --retry 3 --retry-all-errors --retry-delay 3 \
    -o /tmp/commitsar.tar.gz \
    -sSfL "https://github.com/aevea/commitsar/releases/download/${commitsar_version}/commitsar_${commitsar_version#v}_${os}_${arch}.tar.gz"
  tar -zxvf /tmp/commitsar.tar.gz \
    --directory "${ROOT_DIR}/.sbin" \
    --no-same-owner \
    --exclude ./LICENSE \
    --exclude ./*.md
  chmod a+x "${ROOT_DIR}/.sbin/commitsar"
}

function seal::commit::commitsar::validate() {
  # shellcheck disable=SC2046
  if [[ -n "$(command -v $(seal::commit::commitsar::bin))" ]]; then
    if [[ $($(seal::commit::commitsar::bin) version 2>&1 | cut -d " " -f 7 2>&1 | head -n 1 | xargs echo -n) == "${commitsar_version#v}" ]]; then
      return 0
    fi
  fi

  seal::log::info "installing commitsar ${commitsar_version}"
  if seal::commit::commitsar::install; then
    seal::log::info "commitsar $($(seal::commit::commitsar::bin) version 2>&1 | cut -d " " -f 7 2>&1 | head -n 1 | xargs echo -n)"
    return 0
  fi
  seal::log::error "no commitsar available"
  return 1
}

function seal::commit::commitsar::bin() {
  local bin="commitsar"
  if [[ -f "${ROOT_DIR}/.sbin/commitsar" ]]; then
    bin="${ROOT_DIR}/.sbin/commitsar"
  fi
  echo -n "${bin}"
}

function seal::commit::lint() {
  if ! seal::commit::commitsar::validate; then
    seal::log::fatal "cannot execute commitsar as client is not found"
  fi

  seal::log::debug "commitsar $*"
  $(seal::commit::commitsar::bin) "$@"
}
