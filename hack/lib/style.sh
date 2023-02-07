#!/usr/bin/env bash

# -----------------------------------------------------------------------------
# Lint variables helpers. These functions need the
# following variables:
#
#    GOLANGCI_LINT_VERSION  -  The Golangci-lint version, default is v1.50.1.

golangci_lint_version=${GOLANGCI_LINT_VERSION:-"v1.50.1"}

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
  $(seal::lint::golangci_lint::bin) run "$@"
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

function seal::format::run() {
  if ! seal::format::goimports::validate; then
    seal::log::fatal "cannot execute goimports as client is not found"
  fi

  seal::log::debug "goimports $*"
  $(seal::format::goimports::bin) "$@"
}
