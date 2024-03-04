#!/usr/bin/env bash

# -----------------------------------------------------------------------------
# Lint variables helpers. These functions need the
# following variables:
#
# GOIMPORT_REVISER_VERSION  -  The Goimports-reviser version, default is v3.6.4.
#    GOLANGCI_LINT_VERSION  -  The Golangci-lint version, default is v1.56.2.
#        COMMITSAR_VERSION  -  The Commitsar version, default is v0.20.2.

goimports_reviser_version=${GOIMPORT_REVISER_VERSION:-"v3.6.4"}
golangci_lint_version=${GOLANGCI_LINT_VERSION:-"v1.56.2"}
commitsar_version=${COMMITSAR_VERSION:-"v0.20.2"}

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
    seal::log::info "golangci-lint $($(seal::lint::golangci_lint::bin) --version 2>&1 | cut -d " " -f 4 2>&1 | head -n 1)"
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

function seal::lint::goimports_reviser::install() {
  local os
  os="$(seal::util::get_raw_os)"
  local arch
  arch="$(seal::util::get_raw_arch)"
  curl --retry 3 --retry-all-errors --retry-delay 3 \
    -o /tmp/commitsar.tar.gz \
    -sSfL "https://github.com/incu6us/goimports-reviser/releases/download/${goimports_reviser_version}/goimports-reviser_${goimports_reviser_version#v}_${os}_${arch}.tar.gz"
  tar -zxvf /tmp/commitsar.tar.gz \
    --directory "${ROOT_DIR}/.sbin" \
    --no-same-owner \
    --exclude ./LICENSE \
    --exclude ./README.md
  chmod a+x "${ROOT_DIR}/.sbin/goimports-reviser"
}

function seal::lint::goimports_reviser::validate() {
  # shellcheck disable=SC2046
  if [[ -n "$(command -v $(seal::lint::goimports_reviser::bin))" ]]; then
    if [[ $($(seal::lint::goimports_reviser::bin) -version | grep tag | cut -d " " -f 2 2>&1 | head -n 1) == "${goimports_reviser_version}" ]]; then
      return 0
    fi
  fi

  seal::log::info "installing goimports-reviser"
  if seal::lint::goimports_reviser::install; then
    seal::log::info "goimports-reviser $($(seal::lint::goimports_reviser::bin) -version | grep tag | cut -d " " -f 2 2>&1 | head -n 1)"
    return 0
  fi
  seal::log::error "no goimports-reviser available"
  return 1
}

function seal::lint::goimports_reviser::bin() {
  local bin="goimports-reviser"
  if [[ -f "${ROOT_DIR}/.sbin/goimports-reviser" ]]; then
    bin="${ROOT_DIR}/.sbin/goimports-reviser"
  fi
  echo -n "${bin}"
}

function seal::lint::run() {
  if ! seal::lint::goimports_reviser::validate; then
    seal::log::fatal "cannot execute goimports-reviser as client is not found"
  fi

  local goimport_target="${*:$#}"
  goimport_target="${goimport_target//\/.../}"
  local goimports_opts=(
    "-rm-unused"
    "-use-cache"
    "-imports-order=std,general,company,project,blanked,dotted"
    "-output=file"
  )
  local goimports_excludes=()
  for arg in "$@"; do
    if [[ "${arg}" == "--skip-dirs="* ]]; then
      arg="${arg//--skip-dirs=/}"
      goimports_excludes+=("${arg}")
    fi
  done
  if [[ -n "${goimports_excludes[*]}" ]]; then
    goimports_opts+=("-excludes=$(seal::util::join_array "," "${goimports_excludes[@]}")")
  fi
  if [[ "${goimport_target}" == "${ROOT_DIR}" ]]; then
    seal::log::debug "go list -f \"{{.Dir}}\" ./... | xargs -I {} find {} -maxdepth 1 -type f -name '*.go' | xargs -I {} goimports-reviser ${goimports_opts[*]} {}"
    go list -f "{{.Dir}}" ./... | xargs -I {} find {} -maxdepth 1 -type f -name '*.go' | xargs -I {} "$(seal::lint::goimports_reviser::bin)" "${goimports_opts[@]}" {}
  else
    seal::log::debug "pushd \"${goimport_target}\" >/dev/null 2>&1; go list -f \"{{.Dir}}\" ./... | xargs -I {} find {} -maxdepth 1 -type f -name '*.go' | xargs -I {} goimports-reviser ${goimports_opts[*]} {}; popd"
    # shellcheck disable=SC2164
    pushd "${goimport_target}" >/dev/null 2>&1
    go list -f "{{.Dir}}" ./... | xargs -I {} find {} -maxdepth 1 -type f -name '*.go' | xargs -I {} "$(seal::lint::goimports_reviser::bin)" "${goimports_opts[@]}" {}
    # shellcheck disable=SC2164
    popd >/dev/null 2>&1
  fi

  if ! seal::lint::golangci_lint::validate; then
    seal::log::fatal "cannot execute golangci-lint as client is not found"
  fi

  local golangci_lint_opts=(
    "--fix"
  )
  seal::log::debug "golangci-lint run ${golangci_lint_opts[*]} $*"
  $(seal::lint::golangci_lint::bin) run "${golangci_lint_opts[@]}" "$@"
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
    --exclude ./README.md
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
