#!/bin/bash
EXAMPLES_DIR=$1

echo "-- Building .go files in directory ==> ${EXAMPLES_DIR}"

cd ${EXAMPLES_DIR}
for i in *.go; do
    [ -f "$i" ] || break
    go vet $i
done

