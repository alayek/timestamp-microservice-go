#!/usr/env/bin sh

set -exo pipefail

BUILD_DATE=$(date -u)

PACKAGE_NAME="main"

# set it empty initially
LDFLAGS=
LDFLAGS="${LDFLAGS} -X '${PACKAGE_NAME}.BuildDate=${BUILD_DATE}'"

# build command
go build -ldflags "${LDFLAGS}" -o ./timestamp