#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

unset CDPATH

# Set no_proxy for localhost if behind a proxy, otherwise,
# the connections to localhost in scripts will time out.
export no_proxy=127.0.0.1,localhost

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
mkdir -p "${ROOT_DIR}/.sbin"
mkdir -p "${ROOT_DIR}/.dist"

for file in "${ROOT_DIR}/hack/lib/"*; do
  if [[ -f "${file}" ]] && [[ "${file}" != *"init.sh" ]]; then
    # shellcheck disable=SC1090
    source "${file}"
  fi
done

seal::log::install_errexit
seal::version::get_version_vars
