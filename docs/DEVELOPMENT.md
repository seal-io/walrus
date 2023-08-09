# Development Guide

## Setting Up Development Environment

### Local Kubernetes via [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installing-with-a-package-manager)

```bash
$ kind create cluster --name local

```

### Local Postgres via [Docker](https://docs.docker.com/desktop/install/mac-install/)

```bash
$ docker run -d -p 5432:5432 -e "POSTGRES_USER=root" -e "POSTGRES_PASSWORD=Root123" -e "POSTGRES_DB=seal" postgres:14.6

```

### Local Casdoor via [Docker](https://docs.docker.com/desktop/install/mac-install/)

```bash
$ docker run -d -p 8000:8000 sealio/casdoor:v1.344.0-seal.1

```

### (Optional) Local Redis via [Docker](https://docs.docker.com/desktop/install/mac-install/)

```bash
$ docker run -d -p 6379:6379 redis:6.2.11 redis-server --save "" --appendonly no --databases 1 --requirepass Default123
 
```

### Run

```bash
$ # default.
$ go run cmd/server/server.go --log-debug \
  --data-source-address="postgres://root:Root123@127.0.0.1:5432/seal?sslmode=disable" \
  --casdoor-server="http://127.0.0.1:8000"

$ # with cache data source.
$ go run cmd/server/server.go --log-debug \
  --data-source-address="postgres://root:Root123@127.0.0.1:5432/seal?sslmode=disable" \
  --casdoor-server="http://127.0.0.1:8000" \
  --cache-source-address="redis://default:Default123@127.0.0.1:6379/0"

```

#### Interact with [HTTPie](https://httpie.io/docs/cli/macos)

```bash
$ # get init password from console: !!! Bootstrap Admin Password: <here> !!!
$ https --verify=no POST :/account/login username=admin password=<password>

$ # access with login session.
$ https --verify=no GET :/settings Cookie:seal_session=<response from above request>

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
#
#   * [dev] `make package`, embed running resources into a Docker image on one platform.
#           - `REPO=xyz make package` package all targets named with xyz repository.
#           - `VERSION=vX.y.z+l.m make package` package all targets named with vX.y.z-l.m tag.
#           - `TAG=main make package` package all targets named with main tag.
#           - `OS=linux ARCH=arm64 make package` package all targets run on linux/arm64 arch.
#           - `PACKAGE_BUILD=false make package` prepare build resource but disable docker build.
#           - `DOCKER_USERNAME=... DOCKER_PASSWORD=... PACKAGE_PUSH=true make package` execute docker push after build.
#
#   * [ci]  `make ci-check`, execute `make deps`, `make test` and `make lint`.
#
#   * [ci]  `make ci-publish`, execute `make build` and `make package`.
#
```