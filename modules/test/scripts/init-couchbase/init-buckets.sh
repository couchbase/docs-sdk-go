#!/bin/bash

# expand variables and print commands
set -o xtrace

# exit immediately if a command fails or if there are unset vars
set -euo pipefail

CB_USER="${CB_USER:-Administrator}"
CB_PSWD="${CB_PSWD:-password}"

CB_BUCKET_RAMSIZE="${CB_BUCKET_RAMSIZE:-128}"

echo "couchbase-cli bucket-create travel-sample..."
/opt/couchbase/bin/couchbase-cli bucket-create \
-c localhost -u ${CB_USER} -p ${CB_PSWD} \
--bucket travel-sample \
--bucket-type couchbase \
--bucket-ramsize ${CB_BUCKET_RAMSIZE} \
--bucket-replica 0 \
--bucket-priority low \
--bucket-eviction-policy fullEviction \
--enable-flush 1 \
--enable-index-replica 0 \
--wait

echo "cbimport travel-sample..."
/opt/couchbase/bin/cbimport json --format sample --verbose \
-c localhost -u ${CB_USER} -p ${CB_PSWD} \
-b travel-sample \
-d file:///opt/couchbase/samples/travel-sample.zip

echo "curl cmd - create ariports dataset"
curl -u Administrator:password -d 'statement=CREATE DATASET airports ON `travel-sample` WHERE `type`="airport"' http://localhost:8095/analytics/service

echo "sleep 10 to allow stabilization..."
sleep 10
