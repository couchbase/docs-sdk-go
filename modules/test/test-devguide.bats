#!./test/libs/bats/bin/bats

load 'test/test_helper.bash'

@test "[devguide] - analytics-named-placeholders.go" {
    runExample $DEVGUIDE_DIR analytics-named-placeholders.go
    assert_success
}

@test "[devguide] - analytics-positional-placeholders.go" {
    runExample $DEVGUIDE_DIR analytics-positional-placeholders.go
    assert_success
}

@test "[devguide] - analytics-query-one.go" {
    runExample $DEVGUIDE_DIR analytics-query-one.go
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
}

@test "[devguide] - concurrent-batch.go" {
    runExample $DEVGUIDE_DIR concurrent-batch.go 
    assert_success
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
