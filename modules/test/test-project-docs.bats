# #!./test/libs/bats/bin/bats

load 'test/test_helper.bash'

@test "[project-docs] - migrating.go" {
     # Not sure how we can test certificates at the moment, skipping for now.
    skip "Example requires certificates"

    runExample $PROJECT_DOCS_DIR migrating.go
    assert_success
}