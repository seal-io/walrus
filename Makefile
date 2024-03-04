SHELL := /bin/bash

# Borrowed from https://stackoverflow.com/questions/18136918/how-to-get-current-relative-directory-of-your-makefile
curr_dir := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

# Borrowed from https://stackoverflow.com/questions/2214575/passing-arguments-to-make-run
rest_args := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
$(eval $(rest_args):;@:)

targets := $(shell ls $(curr_dir)/hack | grep '.sh' | sed 's/\.sh//g')
$(targets):
	@$(curr_dir)/hack/$@.sh $(rest_args)

help:
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
	@echo

.DEFAULT_GOAL := build
.PHONY: $(targets) cmd client docs gen hack manager server staging
