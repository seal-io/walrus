#!/usr/bin/env bash

# -----------------------------------------------------------------------------
# Image variables helpers. These functions need the
# following variables:
#
#          DOCKER_VERSION   -  The Docker version for running, default is 20.10.
#         DOCKER_USERNAME   -  The username of image registry.
#         DOCKER_PASSWORD   -  The password of image registry.

docker_version=${DOCKER_VERSION:-"20.10"}
docker_username=${DOCKER_USERNAME:-}
docker_password=${DOCKER_PASSWORD:-}

function seal::image::docker::install() {
  curl --retry 3 --retry-all-errors --retry-delay 3 \
    -sSfL "https://get.docker.com" | sh -s VERSION="${docker_version}"
}

function seal::image::docker::validate() {
  # shellcheck disable=SC2046
  if [[ -n "$(command -v $(seal::image::docker::bin))" ]]; then
    return 0
  fi

  seal::log::info "installing docker"
  if seal::image::docker::install; then
    seal::log::info "docker: $($(seal::image::docker::bin) version --format '{{.Server.Version}}' 2>&1)"
    return 0
  fi
  seal::log::error "no docker available"
  return 1
}

function seal::image::docker::bin() {
  echo -n "docker"
}

function seal::image::name() {
  if [[ -n "${IMAGE:-}" ]]; then
    echo -n "${IMAGE}"
  else
    echo -n "$(basename "${ROOT_DIR}")" 2>/dev/null
  fi
}

function seal::image::tag() {
  echo -n "${TAG:-${GIT_VERSION}}" | sed -E 's/[^a-zA-Z0-9\.]+/-/g' 2>/dev/null
}

function seal::image::login() {
  if seal::image::docker::validate; then
    if [[ -n ${docker_username} ]] && [[ -n ${docker_password} ]]; then
      seal::log::debug "docker login ${*:-} -u ${docker_username} -p ***"
      if ! docker login "${*:-}" -u "${docker_username}" -p "${docker_password}" >/dev/null 2>&1; then
        seal::log::fatal "failed: docker login ${*:-} -u ${docker_username} -p ***"
      fi
    fi
    return 0
  fi

  seal::log::fatal "cannot execute image login as client is not found"
}

function seal::image::build::within_container() {
  if seal::image::docker::validate; then
    if ! $(seal::image::docker::bin) buildx inspect --builder="seal"; then
      seal::log::debug "setting up qemu"
      $(seal::image::docker::bin) run \
        --rm \
        --privileged \
        tonistiigi/binfmt:qemu-v7.0.0 --install all
      seal::log::debug "setting up buildx"
      $(seal::image::docker::bin) buildx create \
        --name="seal" \
        --driver="docker-container" \
        --buildkitd-flags="--allow-insecure-entitlement security.insecure --allow-insecure-entitlement network.host" \
        --use \
        --bootstrap
    fi

    return 0
  fi

  seal::log::fatal "cannot execute image build as client is not found"
}

function seal::image::build() {
  if seal::image::docker::validate; then
    seal::log::debug "docker buildx build $*"
    $(seal::image::docker::bin) buildx build "$@"

    return 0
  fi

  seal::log::fatal "cannot execute image build as client is not found"
}
