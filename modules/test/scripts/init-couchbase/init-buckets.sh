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

sleep 5

echo "cbimport travel-sample..."
/opt/couchbase/bin/cbimport json --format sample --verbose \
-c localhost -u ${CB_USER} -p ${CB_PSWD} \
-b travel-sample \
-d file:///opt/couchbase/samples/travel-sample.zip

echo "create ariports dataset"
curl --fail -v -u ${CB_USER}:${CB_PSWD} -H "Content-Type: application/json" -d '{
    "statement": "CREATE DATASET airports ON `travel-sample` WHERE `type`=\"airport\";",
    "pretty":true,
    "client_context_id":"test"
}' http://localhost:8095/analytics/service

echo "sleep 10 to allow stabilization..."
sleep 10

echo "create travel-sample-index"
curl --fail -v -u Administrator:password -X PUT \
http://localhost:8094/api/index/travel-sample-index \
-H 'cache-control: no-cache' \
-H 'content-type: application/json' \
-d '{
        "name": "travel-sample-index",
        "type": "fulltext-index",
        "params": {
            "doc_config": {
                "docid_prefix_delim": "",
                "docid_regexp": "",
                "mode": "type_field",
                "type_field": "type"
            },
            "mapping": {
                "default_analyzer": "standard",
                "default_datetime_parser": "dateTimeOptional",
                "default_field": "_all",
                "default_mapping": {
                    "dynamic": false,
                    "enabled": true,
                    "properties": {
                        "description": {
                            "enabled": true,
                            "dynamic": false,
                            "fields": [
                                {
                                    "docvalues": true,
                                    "include_in_all": true,
                                    "include_term_vectors": true,
                                    "index": true,
                                    "name": "description",
                                    "store": true,
                                    "type": "text"
                                }
                            ]
                        },
                        "type": {
                            "enabled": true,
                            "dynamic": false,
                            "fields": [
                                {
                                    "docvalues": true,
                                    "include_in_all": true,
                                    "include_term_vectors": true,
                                    "index": true,
                                    "name": "type",
                                    "store": true,
                                    "type": "text"
                                }
                            ]
                        }
                    }
                },
                "default_type": "_default",
                "docvalues_dynamic": true,
                "index_dynamic": true,
                "store_dynamic": false,
                "type_field": "_type"
            },
            "store": {
                "indexType": "scorch",
                "segmentVersion": 15
            }
        },
        "sourceType": "gocbcore",
        "sourceName": "travel-sample",
        "sourceParams": {},
        "planParams": {
            "maxPartitionsPerPIndex": 1024,
            "indexPartitions": 1,
            "numReplicas": 0
        }
}'

echo "sleep 20 to allow search index to load..."
sleep 20