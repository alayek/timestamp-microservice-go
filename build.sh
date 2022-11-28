#!/usr/env/bin sh

set -euxo pipefail

if [ -z "${GIT_VERSION}"]; then
  GIT_VERSION=$(git rev-parse --short HEAD)
fi
PACKAGE_NAME="main"

# set it empty initially
LDFLAGS=
LDFLAGS="${LDFLAGS} -X '${PACKAGE_NAME}.CommitID=${GIT_VERSION}'"

# build command
go build -ldflags "${LDFLAGS}" -o ./timestamp