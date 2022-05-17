#!/usr/bin/env bash

set -eu

# Runs docker container.
# Main purpose is to use on Windows but can be used by all, especially to avoid
# installing required tooling (make, go, etc).
#
# Usage:
#
#   # will target host's OS and architecture
#   ./build.sh make build
#   ./build.sh make all
#
#   # override target OS/architecture
#   Z_GOOS=linux Z_GOARCH=amd ./build.sh make build
#   Z_GOOS=linux Z_GOARCH=amd ./build.sh make all
#
Z_GOOS="${Z_GOOS:-"$(go env GOOS)"}"
Z_GOARCH="${Z_GOARCH:-"$(go env GOARCH)"}"

DOCKER_IMAGE="golang:1.17"

cd "$(dirname "${BASH_SOURCE}")"

echo "Running in docker container ${DOCKER_IMAGE}..."

(
    # Make sure we only attach TTY if we have it, CI builds won't have it.
    declare TTY_FLAG=""
    if [ -t 1 ]
    then
        TTY_FLAG="-t"
    fi

    # Annoying issue with ownership of files in mapped volumes.
    # Need to run with same UID and GID in container as we do
    # on the machine, otherwise all output will be owned by root.
    # Doesn't happen on OS X but does on Linux. So we will do
    # UID and GID for Linux only (this won't work on OS X anyway).
    declare USER_FLAG=""
    if test "Linux" == "$(uname || true)"
    then
        USER_FLAG="-u $(id -u):$(id -g)"
    fi


    # Another annoying Linux issue: failed to initialize build cache at /.cache/go-build
    # This is because we don't have HOME defined for Linux.
    # See:
    #   - https://forum.golangbridge.org/t/go-build-failed/19551
    #   - https://github.com/openshift/release/issues/9748
    Z_GOCACHE="$(go env GOCACHE)"

    set -x
    docker run \
      -i \
      ${TTY_FLAG} \
      ${USER_FLAG} \
      --rm \
      -v "${PWD}:${PWD}":cached \
      -v "${Z_GOCACHE}":"${Z_GOCACHE}":cached \
      -w ${PWD} \
      -e GOCACHE="${Z_GOCACHE}" \
      -e Z_GOOS="${Z_GOOS}" \
      -e Z_GOARCH="${Z_GOARCH}" \
      -e Z_GOBIN="${PWD}/bin_tools/linux_amd64" \
      ${DOCKER_IMAGE} \
      "$@"
  )

