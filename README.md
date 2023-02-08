# Seal

**!!! WIP !!!**

## Develop Command

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

## Develop Environment

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
$ docker run -d -p 8000:8000 sealio/casdoor:v1.197.0-seal.2

```

### Run

```bash
$ go run cmd/server/server.go --log-debug --data-source-address="postgres://root:Root123@127.0.0.1:5432/seal?sslmode=disable" --casdoor-server="http://127.0.0.1:8000"

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

# License

Copyright (c) 2023 [Seal, Inc.](https://seal.io)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at [LICENSE](./LICENSE) file for details.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
