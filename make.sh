#!/usr/bin/env bash

set -o errexit
set -o pipefail

# Use xtrace for debugging.
#set -o xtrace

# Capture arguments then disable unset variables.
FUNC=$1
set -o nounset

REPO_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BUILD_DIR="${REPO_ROOT}/builds"

VERSION="v0.1.0"
COMMIT=$(git rev-list -1 HEAD)
DATE=$(date -uR)
GOVERSION=$(go version | awk '{print $3 " " $4}')

function build() {
    cd "$REPO_ROOT"
    BUILD_GOOS=${1:-$(go env GOOS)}
    BUILD_GOARCH=${2:-$(go env GOARCH)}
    BUILD_OUTPUT="${BUILD_DIR}/qlcplus-http-api-${VERSION}-${BUILD_GOOS}-${BUILD_GOARCH}/qlcplus-http-api"
    LDFLAGS="-X \"main.Version=${VERSION}\""
    LDFLAGS="${LDFLAGS} -X \"main.Platform=${BUILD_GOOS}/${BUILD_GOARCH}\""
    LDFLAGS="${LDFLAGS} -X \"main.Commit=${COMMIT}\""
    LDFLAGS="${LDFLAGS} -X \"main.BuildDate=${DATE}\""
    LDFLAGS="${LDFLAGS} -X \"main.GoVersion=${GOVERSION}\""
    echo "Building qlcplus-http-api ${VERSION} ${BUILD_GOOS} ${BUILD_GOARCH}."
    GOOS=$BUILD_GOOS GOARCH=$BUILD_GOARCH CGO_ENABLED=0 GO111MODULE=on go build -ldflags "${LDFLAGS}" -o "${BUILD_OUTPUT}" .
}

function build-all() {
    build linux amd64
    build linux arm
    build linux arm64
    build darwin amd64
    build windows amd64
}

function package-all() {
    if [ ! -d "${BUILD_DIR}" ]; then
        echo "No builds to package."
        exit
    fi
    for BUILD in "${BUILD_DIR}"/*; do
        if [ ! -d "${BUILD}" ]; then
            continue
        fi
        ZIP_NAME=${BUILD_DIR}/$(basename "${BUILD}").zip
        if [ -f "${ZIP_NAME}" ]; then
            rm "${ZIP_NAME}"
        fi
        echo "Packaging build ${ZIP_NAME}."
        cd "${BUILD}"
        zip -q "${ZIP_NAME}" ./*
    done
}

function install() {
    build
    cp "${BUILD_OUTPUT}" "${GOPATH}/bin"
}

function clean() {
    rm -r "${BUILD_DIR}"
}

# Check the first argument passed is a function.
if [ "$(type -t $FUNC)" != "function" ]; then
    # If not warn and print the list of functions.
    echo "No target: $FUNC"
    echo "Try: $0 { $(compgen -A function | tr '\n' ' ')}"
    exit 1
fi
# Run the function named by the first argument.
$FUNC
