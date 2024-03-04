#!/usr/bin/env bash

# -----------------------------------------------------------------------------
# Protoc variables helpers. These functions need the
# following variables:
#
#              PROTOC_VERSION  -  The protoc version, default is v23.4.
#     PROTOC_GEN_GOGO_VERSION  -  The protoc version, default is master.
#       PROTOC_GEN_GO_VERSION  -  The protoc-gen-go version, default is v1.5.4.
#  PROTOC_GEN_GO_GRPC_VERSION  -  The protoc-gen-go-grpc version, default is v1.62.1.
#
# Refs:
# - https://grpc.io/docs/protoc-installation/
# - https://grpc.io/docs/languages/go/quickstart/

protoc_version=${PROTOC_VERSION:-"v23.4"}
protoc_gen_go_version=${PROTOC_GEN_GO_VERSION:-"v1.5.4"}
protoc_gen_go_grpc_version=${PROTOC_GEN_GO_GRPC_VERSION:-"v1.62.1"}

function seal::protoc::protoc::install() {
  local os
  os=$(echo -n "$(uname -s)" | tr '[:upper:]' '[:lower:]')
  if [[ "${os}" == "darwin" ]]; then
    os="osx"
  fi
  local arch
  arch=$(uname -m)
  if [[ "${arch}" == "arm64" ]]; then
    arch="aarch_64"
  fi
  curl --retry 3 --retry-all-errors --retry-delay 3 \
    -o /tmp/protoc.zip \
    -sSfL "https://github.com/protocolbuffers/protobuf/releases/download/${protoc_version}/protoc-${protoc_version#v}-${os}-${arch}.zip"
  unzip -qu /tmp/protoc.zip -d "${ROOT_DIR}/.sbin/protoc"
}

function seal::protoc::protoc::validate() {
  # shellcheck disable=SC2046
  if [[ -n "$(command -v $(seal::protoc::bin))" ]]; then
    if [[ "$($(seal::protoc::bin) --version 2>&1)" == "libprotoc ${protoc_version#v}" ]]; then
      return 0
    fi
  fi

  seal::log::info "installing protoc ${protoc_version}"
  if seal::protoc::protoc::install; then
    seal::log::info "protoc $($(seal::protoc::bin) --version 2>&1 | cut -d " " -f 2 2>&1 | head -n 1)"
    return 0
  fi
  seal::log::error "no protoc available"
  return 1
}

function seal::protoc::protoc_gen_go::install() {
  local bin="${ROOT_DIR}/.sbin/protoc/bin"
  mkdir -p "${bin}"
  GOBIN="${bin}" go install "google.golang.org/protobuf/cmd/protoc-gen-go@${protoc_gen_go_version}"
}

function seal::protoc::protoc_gen_go::validate() {
  PATH="${PATH}:${ROOT_DIR}/.sbin/protoc/bin"
  if [[ -f "${ROOT_DIR}/.sbin/protoc/bin/protoc-gen-go" ]]; then
    return 0
  fi

  seal::log::info "installing protoc-gen-go ${protoc_gen_go_version}"
  if seal::protoc::protoc_gen_go::install; then
    seal::log::info "protoc-gen-go installed"
    return 0
  fi
  seal::log::error "no protoc-gen-go available"
  return 1
}

function seal::protoc::protoc_gen_go_grpc::install() {
  local bin="${ROOT_DIR}/.sbin/protoc/bin"
  mkdir -p "${bin}"
  GOBIN="${bin}" go install "google.golang.org/grpc/cmd/protoc-gen-go-grpc@${protoc_gen_go_grpc_version}"
}

function seal::protoc::protoc_gen_go_grpc::validate() {
  PATH="${PATH}:${ROOT_DIR}/.sbin/protoc/bin"
  if [[ -f "${ROOT_DIR}/.sbin/protoc/bin/protoc-gen-go-grpc" ]]; then
    return 0
  fi

  seal::log::info "installing protoc-gen-go-grpc ${protoc_gen_go_grpc_version}"
  if seal::protoc::protoc_gen_go_grpc::install; then
    seal::log::info "protoc-gen-go-grpc installed"
    return 0
  fi
  seal::log::error "no protoc-gen-go-grpc available"
  return 1
}

function seal::protoc::bin() {
  local bin="protoc"
  if [[ -f "${ROOT_DIR}/.sbin/protoc/bin/protoc" ]]; then
    bin="${ROOT_DIR}/.sbin/protoc/bin/protoc"
  fi
  echo -n "${bin}"
}

function seal::protoc::generate() {
  if ! seal::protoc::protoc::validate; then
    seal::log::error "cannot execute protoc as it hasn't installed"
    return 1
  fi

  if ! seal::protoc::protoc_gen::validate; then
    seal::log::error "cannot execute protoc-gen as it hasn't installed"
    return 1
  fi

  if ! seal::protoc::protoc_gen_grpc::validate; then
    seal::log::error "cannot execute protoc-gen-grpc as it hasn't installed"
    return 1
  fi

  local filepath="${1:-}"
  if [[ ! -f ${filepath} ]]; then
    seal::log::warn "${filepath} isn't existed"
    return 1
  fi
  local filedir
  filedir=$(dirname "${filepath}")
  local filename
  filename=$(basename "${filepath}" ".proto")

  # generate
  $(seal::protoc::bin) \
    --proto_path="${filedir}" \
    --proto_path="${ROOT_DIR}/.sbin/protoc/include" \
    --go_out="${filedir}" \
    --go_opt=paths=source_relative \
    --go-grpc_out="${filedir}" \
    --go-grpc_opt=paths=source_relative \
    "${filepath}"

  # format
  local tmpfile
  tmpfile=$(mktemp)
  local generated_files=(
    "${filedir}/${filename}.pb.go"
    "${filedir}/${filename}_grpc.pb.go"
  )
  for generated_file in "${generated_files[@]}"; do
    if [[ -f "${generated_file}" ]]; then
      sed "2d" "${generated_file}" >"${tmpfile}" &&
        mv "${tmpfile}" "${generated_file}"
      cat "${ROOT_DIR}/hack/boilerplate/go.txt" "${generated_file}" >"${tmpfile}" &&
        mv "${tmpfile}" "${generated_file}"
      go fmt "${generated_file}" >/dev/null 2>&1
    fi
  done
}
