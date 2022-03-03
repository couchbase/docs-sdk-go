#!./test/libs/bats/bin/bats

load 'test_helper'

# Test is a bit flaky on the last assertion if not run first. 
# It seems the search index updates when the other tests run, causing major delay on the
# expected result for the last assertion.
@test "[devguide] - search.go" {
    # TODO: The tags conjunctionquery[] and daterangequery[] don't seem
    # to return any results, needs further investigation.

    runExample $DEVGUIDE_DIR search.go
    assert_success
    assert_output --partial "Document ID: hotel_26223"
    assert_output --partial "fields included in result: map[_\$c:hotel description:Swanky "
    refute_output --partial "Facet field: type, total: 0"
    assert_output --partial "Document ID: a-new-hotel, search score:"
}

@test "[devguide] - analytics-named-placeholders.go" {
    runExample $DEVGUIDE_DIR analytics-named-placeholders.go
    assert_success
    assert_output --partial "Result count: 221"
}

@test "[devguide] - analytics-positional-placeholders.go" {
    runExample $DEVGUIDE_DIR analytics-positional-placeholders.go
    assert_success
    assert_output --partial "Result count: 221"
}

@test "[devguide] - analytics-query-one.go" {
    runExample $DEVGUIDE_DIR analytics-query-one.go
    assert_success
}

@test "[devguide] - analytics-collection-scope.go" {
    runExample $DEVGUIDE_DIR analytics-collection-scope.go
    assert_success
}

@test "[devguide] - analytics-simple-query.go" {
    runExample $DEVGUIDE_DIR analytics-simple-query.go
    assert_success
}

@test "[devguide] - cloud.go" {
    # Not sure how we can test cloud examples at this point in time, skipping for now.
    skip "Example requires a cloud endpoint"

    runExample $DEVGUIDE_DIR cloud.go
    assert_success
}

@test "[devguide] - concurrent-async.go" {
    runExample $DEVGUIDE_DIR concurrent-async.go 
    assert_success
    assert_output --partial "Completed"
}

@test "[devguide] - concurrent-batch.go" {
    runExample $DEVGUIDE_DIR concurrent-batch.go 
    assert_success
    assert_output --partial "Loaded 7303 docs"
    assert_output --partial "Completed"
}

@test "[devguide] - connecting-cca.go" {
      # Not sure how we can test certificates at the moment, skipping for now.
    skip "Example requires certificates"

    runExample $DEVGUIDE_DIR connecting-cca.go 
    assert_success
}

@test "[devguide] - connecting.go" {
     # Not sure how we can test multiple couchbase server nodes at the moment, skipping for now.
    skip "Example requires multiple nodes"

    runExample $DEVGUIDE_DIR connecting.go 
    assert_success
}

@test "[devguide] - custom-logging.go" {
    runExample $DEVGUIDE_DIR custom-logging.go 
    assert_success
}

@test "[devguide] - error-handling.go" {
    runExample $DEVGUIDE_DIR error-handling.go
    assert_success
}

@test "[devguide] - healthcheck.go" {
    runExample $DEVGUIDE_DIR healthcheck.go
    assert_success
}

@test "[devguide] - kv-counter.go" {
    runExample $DEVGUIDE_DIR kv-counter.go
    assert_success
}

@test "[devguide] - kv-crud.go" {
    runExample $DEVGUIDE_DIR kv-crud.go
    assert_success
}

@test "[devguide] - kv-durability-enhanced.go" {
    runExample $DEVGUIDE_DIR kv-durability-enhanced.go
    assert_success
}

@test "[devguide] - kv-durability-observe.go" {
    EXPECTED_OUTPUT=$(cat <<-EOF
Document Retrieved: Durabilty PersistTo Test Value
Document Retrieved: Durabilty ReplicateTo Test Value
Document Retrieved: Durabilty ReplicateTo and PersistTo Test Value
EOF
)
    runExample $DEVGUIDE_DIR kv-durability-observe.go
    assert_success
    assert_output --partial "$EXPECTED_OUTPUT"
}

@test "[devguide] - kv-expiry.go" {
    runExample $DEVGUIDE_DIR kv-expiry.go
    assert_success
}

@test "[devguide] - n1ql-query-consistency.go" {
    runExample $DEVGUIDE_DIR n1ql-query-consistency.go
    assert_success
}

@test "[devguide] - n1ql-query-consistentwith.go" {
    # TODO: Remove `skip` once Couchbase Server 7.0.1 is available.
    skip "BUG: https://issues.couchbase.com/browse/MB-46876"

    runExample $DEVGUIDE_DIR n1ql-query-consistentwith.go
    assert_success
}

@test "[devguide] - n1ql-query-create-index.go" {
    runExample $DEVGUIDE_DIR n1ql-query-create-index.go
    assert_success
}

@test "[devguide] - n1ql-query-metrics.go" {
    runExample $DEVGUIDE_DIR n1ql-query-metrics.go
    assert_success
}

@test "[devguide] - n1ql-query-named-placeholders.go" {
    runExample $DEVGUIDE_DIR n1ql-query-named-placeholders.go
    assert_success
}

@test "[devguide] - n1ql-query-one.go" {
    runExample $DEVGUIDE_DIR n1ql-query-one.go
    assert_success
}

@test "[devguide] - n1ql-query-positional-placeholders.go" {
    runExample $DEVGUIDE_DIR n1ql-query-positional-placeholders.go
    assert_success
}

@test "[devguide] - n1ql-query-simple.go" {
    runExample $DEVGUIDE_DIR n1ql-query-simple.go
    assert_success
}

@test "[devguide] - orphan-logging.go" {
    runExample $DEVGUIDE_DIR orphan-logging.go
    assert_success
}

@test "[devguide] - provisioning-resources-buckets.go" {
    runExample $DEVGUIDE_DIR provisioning-resources-buckets.go
    assert_success
}

@test "[devguide] - provisioning-resources-views.go" {
    runExample $DEVGUIDE_DIR provisioning-resources-views.go
    assert_success
}

@test "[devguide] - subdoc-counter.go" {
    runExample $DEVGUIDE_DIR subdoc-counter.go
    assert_success
}

@test "[devguide] - subdoc-durability.go" {
    runExample $DEVGUIDE_DIR subdoc-durability.go
    assert_success
}

@test "[devguide] - subdoc-lookupin.go" {
    runExample $DEVGUIDE_DIR subdoc-lookupin.go
    assert_success
}

@test "[devguide] - subdoc-mutatein-arrays.go" {
    runExample $DEVGUIDE_DIR subdoc-mutatein-arrays.go
    assert_success
}

@test "[devguide] - subdoc-mutatein.go" {
    runExample $DEVGUIDE_DIR subdoc-mutatein.go
    assert_success
}

@test "[devguide] - threshold-logging.go" {
    runExample $DEVGUIDE_DIR threshold-logging.go
    assert_success
}

@test "[devguide] - transcoding-custom.go" {
    # Example doesn't seem to be runnable (possibly update in future), 
    # we can be satisfied that it compiles/builds correctly at this point in time.
    # Skipping for now.
    skip "Example is not runnable at this point in time"

    runExample $DEVGUIDE_DIR transcoding-custom.go
    assert_success
}

@test "[devguide] - transcoding-rawbinary.go" {
    runExample $DEVGUIDE_DIR transcoding-rawbinary.go
    assert_success
}

@test "[devguide] - transcoding-rawjson.go" {
    runExample $DEVGUIDE_DIR transcoding-rawjson.go
    assert_success
}

@test "[devguide] - transcoding-rawstring.go" {
    runExample $DEVGUIDE_DIR transcoding-rawstring.go
    assert_success
}

@test "[devguide] - user-management.go" {
    runExample $DEVGUIDE_DIR user-management.go
    assert_success
}

@test "[devguide] - using-cas.go" {
    runExample $DEVGUIDE_DIR using-cas.go
    assert_success
}

@test "[devguide] - views-key.go" {
EXPECTED_OUTPUT=$(cat <<-EOF
Document ID: landmark_26480
Landmark named Circle Bar has value <nil>
EOF
)

    runExample $DEVGUIDE_DIR views-key.go
    assert_success
    assert_output --partial "$EXPECTED_OUTPUT"
}

@test "[devguide] - views-startkey.go" {
EXPECTED_OUTPUT=$(cat <<-EOF
Document ID: landmark_16320
Landmark named United Kingdom has value <nil>
Document ID: landmark_25731
Landmark named United States has value <nil>
Document ID: landmark_26480
Landmark named United States has value <nil>
Document ID: landmark_37519
Landmark named United States has value <nil>
EOF
)

    runExample $DEVGUIDE_DIR views-startkey.go
    assert_success
    assert_output --partial "$EXPECTED_OUTPUT"
}

@test "[devguide] - kv-collection-scope.go" {
    runExample $DEVGUIDE_DIR kv-collection-scope.go
    assert_success
}

@test "[devguide] - fle.go" {
    runExample $DEVGUIDE_DIR fle.go
    assert_success
    assert_output --partial "FirstName:Barry LastName:Sheen Password:bang!"
    assert_output --partial "Addresses:[{HouseName:my house StreetName:my street}"
    assert_output --partial "{HouseName:my other house StreetName:my other street}] Phone:123456}"
}

@test "[devguide] - slow-operations.go" {
    runExample $DEVGUIDE_DIR slow-operations.go
    assert_success
}
