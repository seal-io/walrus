# Development Guide

## Environment Setup

### Prerequisites

- [Install Go](https://golang.org/doc/install): Required, for `go run`, and [version sensitive](../go.mod).
- Provide local Kubernetes cluster.
    - [x] [OrbStack](https://docs.orbstack.dev/install): Approved, nice to work, but unsupported in Windows platform at
      present.
    - [ ] [Kubernetes of Docker Desktop](https://docs.docker.com/desktop/kubernetes/): Unsupported, bypass APIService is
      not working, it doesn't support tuning.
    - [ ] [KinD](https://kind.sigs.k8s.io/docs/user/quick-start/#installing-with-a-package-manager): Unsupported, bypass
      APIService is not working properly and requires further investigation.

### Go Run

```bash
$ # Setup local Kubernetes cluster.

$ # Run the Walrus without auths.
$ go run cmd/server/main.go --v=4 --kube-leader-election=false --disable-auths=true

```

#### Interact with [Kubectl](https://kubernetes.io/docs/tasks/tools/)

```bash
$ # Operate with Walrus API.
$ kubectl get apiservices | grep walrus

```

#### Interact with [Swagger UI](https://github.com/swagger-api/swagger-ui)

```bash
$ # Open Swagger UI.
$ open https://localhost/swagger/

```

## Commands

The Makefile includes some useful commands for development. You can build locally using the `make` command. You can
run `make help` for details. The output is shown as below.

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
