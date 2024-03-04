# Development Guide

## Setup Development Environment

### Local Kubernetes via [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installing-with-a-package-manager)

```bash
$ kind create cluster --name local

```

### Go Run

```bash
$ # default.
$ go run cmd/server/main.go --log-verbosity=4

$ # with specified kubernetes cluster.
$ go run cmd/server/main.go --log-verbosity=4 --kubeconfig=/path/to/kubeconfig
```

#### Interact with [HTTPie](https://httpie.io/docs/cli/macos)

```bash
$ # get init password from console: !!! Bootstrap Admin Password: <here> !!!
$ https --verify=no POST :/account/login username=admin password=<password>

$ # access with login session.
$ https --verify=no GET :/settings Cookie:walrus_session=<response from above request>

```

#### Interact with [Swagger UI](https://github.com/swagger-api/swagger-ui)

```bash
$ open http://127.0.0.1/swagger

```

## Development Commands

The Makefile includes some useful commands for development. You can build locally using the `make` command. You can run `make help` for details. The output is shown as below.

```bash
$ make help
#
# Usage:
#
#   * [dev] `make deps`, get dependencies.
#           - `make deps update`, update dependencies.
#
#   * [dev] `make generate`, generate something.
#
#   * [dev] `make lint`, check style.
#           - `BUILD_TAGS="jsoniter" make lint` check with specified tags.
#           - `LINT_DIRTY=true make lint` verify whether the code tree is dirty.
#
#   * [dev] `make test`, execute unit testing.
#           - `BUILD_TAGS="jsoniter" make test` test with specified tags.
#
#   * [dev] `make build`, execute cross building.
#           - `VERSION=vX.y.z+l.m make build` build all targets with vX.y.z+l.m version.
#           - `OS=linux ARCH=arm64 make build` build all targets run on linux/arm64 arch.
#           - `BUILD_TAGS="jsoniter" make build` build with specified tags.
#           - `BUILD_PLATFORMS="linux/amd64,linux/arm64" make build` do multiple platforms go build.
#           - `WALRUS_TELEMETRY_API_KEY="phc_xxx" make build` build with telemetry api key.
#
#   * [dev] `make package`, embed running resources into a Docker image on one platform.
#           - `REPO=xyz make package` package all targets named with xyz repository.
#           - `VERSION=vX.y.z+l.m make package` package all targets named with vX.y.z-l.m tag.
#           - `TAG=main make package` package all targets named with main tag.
#           - `OS=linux ARCH=arm64 make package` package all targets run on linux/arm64 arch.
#           - `PACKAGE_BUILD=false make package` prepare build resource but disable docker build.
#           - `DOCKER_USERNAME=... DOCKER_PASSWORD=... PACKAGE_PUSH=true make package` execute docker push after build.
#
#   * [ci]  `make ci`, execute `make deps`, `make lint`, `make test`, `make build` and `make package`.
#           - `CI_CHECK=false make ci` only execute `make build` and `make package`.
#           - `CI_PUBLISH=false make ci` only execute `make deps`, `make lint` and `make test`.
#
```