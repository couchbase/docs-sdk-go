#!/bin/bash

# exit immediately if a command fails or if there are unset vars
set -euo pipefail

EXAMPLES_DIR=$1

cd ${EXAMPLES_DIR}

echo "-- Building .go files in directory ==> ${EXAMPLES_DIR}"
for i in *.go; do
    [ -f "$i" ] || break
    
    echo "Building file: ${i}"
    go vet ${i}
done