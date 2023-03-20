#!/bin/bash

# ===============================================
# NOTE: Any changes made to this file will not be 
# automatically reflected in `make cb-start` as
# the Makefile does not use the mounted version
# of this file in the Docker image. You will need
# to rebuild the image via `make cb-build` before
# running `make cb-start`.
# ===============================================

# exit immediately if a command fails or if there are unset vars
set -euo pipefail

CB_USER="${CB_USER:-Administrator}"
CB_PSWD="${CB_PSWD:-password}"
CB_HOST=localhost

CB_BUCKET_RAMSIZE="${CB_BUCKET_RAMSIZE:-128}"

echo "couchbase-cli bucket-create travel-sample..."
/opt/couchbase/bin/couchbase-cli bucket-create \
    -c ${CB_HOST} -u ${CB_USER} -p ${CB_PSWD} \
    --bucket travel-sample \
    --bucket-type couchbase \
    --bucket-ramsize ${CB_BUCKET_RAMSIZE} \
    --bucket-replica 0 \
    --bucket-priority low \
    --bucket-eviction-policy fullEviction \
    --enable-flush 1 \
    --enable-index-replica 0 \
    --wait

sleep 5

echo "cbimport travel-sample..."
/opt/couchbase/bin/cbimport json --format sample --verbose \
    -c ${CB_HOST} -u ${CB_USER} -p ${CB_PSWD} \
    -b travel-sample \
    -d file:///opt/couchbase/samples/travel-sample.zip

echo "create airports dataset"
curl --fail -s -u ${CB_USER}:${CB_PSWD} -H "Content-Type: application/json" -d '{
    "statement": "CREATE DATASET airports ON `travel-sample` WHERE `type`=\"airport\";",
    "pretty":true,
    "client_context_id":"test"
}' http://${CB_HOST}:8095/analytics/service

echo "create scoped airport dataset"
echo "Skipping..."
# These are already setup in the official Couchbase Enterprise docker image.
#curl --fail -v -u ${CB_USER}:${CB_PSWD} -H "Content-Type: application/json" -d '{
#    "statement": "ALTER COLLECTION `travel-sample`.`inventory`.`airport` ENABLE ANALYTICS;",
#    "pretty":true,
#    "client_context_id":"test"
#}' http://${CB_HOST}:8095/analytics/service

curl --fail -v -u ${CB_USER}:${CB_PSWD} -H "Content-Type: application/json" -d '{
    "statement": "CONNECT LINK Local;",
    "pretty":true,
    "client_context_id":"test"
}' http://${CB_HOST}:8095/analytics/service

echo "sleep 10 to allow stabilization..."
sleep 10

echo
echo "create travel-sample-index"
curl --fail -s -u ${CB_USER}:${CB_PSWD} -X PUT \
    http://${CB_HOST}:8094/api/index/travel-sample-index \
    -H 'cache-control: no-cache' \
    -H 'content-type: application/json' \
    -d @/init-couchbase/travel-sample-index.json

echo
counter=0

until curl -s -u ${CB_USER}:${CB_PSWD} http://${CB_HOST}:8094/api/index/travel-sample-index/count |
    jq -e '.count' | grep 31591 >/dev/null; do # there are 31591 docs to be processed in this index...
    echo -e "Waiting for travel-sample-index to be ready. Current item count: $(curl -s -u ${CB_USER}:${CB_PSWD} http://${CB_HOST}:8094/api/index/travel-sample-index/count | jq -e '.count')/31591 ($counter seconds elapsed)"
    counter=$((counter+20))
    sleep 20
done
echo "travel-sample-index ready"

echo "Done."
