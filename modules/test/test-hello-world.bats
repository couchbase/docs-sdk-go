# #!./test/libs/bats/bin/bats

load 'test_helper'

@test "[hello-world] - startusing.go" {
    runExample $HELLO_WORLD_DIR startusing.go
    assert_success
    assert_output --partial "User: {Jade jade@test-email.com [Swimming Rowing]}"
}
