#!/bin/bash

EXAMPLES_DIR=$1

GO_VERSION=$2
SDK_VERSION=$3

# exit immediately if a command fails or if there are unset vars
set -euo pipefail

cd ${EXAMPLES_DIR}

# Set go module version to what has been provided.
# NOTE: Does nothing if the version is the same as the default.
go mod edit -go=${GO_VERSION}

# Sets gocb version to what has been provided.
# NOTE: Does nothing if the version is the same as the default.
go get github.com/couchbase/gocb/v2@${SDK_VERSION}

# Run 'go mod tidy' to update go.sum file and cleanup any old dependencies.
go mod tidy
