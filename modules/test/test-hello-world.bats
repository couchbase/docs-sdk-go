# #!./test/libs/bats/bin/bats

load 'test/test_helper.bash'

@test "[hello-world] - startusing.go" {
    runExample $HELLO_WORLD_DIR startusing.go
    assert_success
    assert_output --partial "map[name:mike]"
}
