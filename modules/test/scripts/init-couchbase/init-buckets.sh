#!/bin/bash

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

echo "create ariports dataset"
curl --fail -s -u ${CB_USER}:${CB_PSWD} -H "Content-Type: application/json" -d '{
    "statement": "CREATE DATASET airports ON `travel-sample` WHERE `type`=\"airport\";",
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
echo "Waiting for travel-sample-index to be ready..."
until curl --fail -s -u ${CB_USER}:${CB_PSWD} http://${CB_HOST}:8094/api/index/travel-sample-index/count |
    jq -e '.count' | grep 31591 >/dev/null; do # there are 31591 docs to be processed in this index...
    echo "Waiting for travel-sample-index to be ready. Trying again in 10 seconds."
    sleep 10
done

echo "Done."
