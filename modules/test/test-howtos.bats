# #!./test/libs/bats/bin/bats

load 'test/test_helper.bash'

@test "[howtos] - analytics.go" {
    runExample $HOWTOS_DIR analytics.go
    assert_success
}

@test "[howtos] - query.go" {
   runExample $HOWTOS_DIR query.go
   assert_success
}

@test "[howtos] - subdoc.go" {
    echo "create example document"
    cbimport json --verbose \
    -c localhost -u Administrator -p password \
    -b travel-sample \
    -f lines \
    -d file:///modules/test/fixtures/customer123.json \
    -g customer123
    echo

    runExample $HOWTOS_DIR subdoc.go
    assert_success
}
