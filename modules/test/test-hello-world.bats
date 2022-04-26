# #!./test/libs/bats/bin/bats

load 'test_helper'

@test "[hello-world] - startusing.go" {
    skip "Requires TLS certificate, skipping for now"

    runExample $HELLO_WORLD_DIR startusing.go
    assert_success
    assert_output --partial "map[name:mike]"
}
