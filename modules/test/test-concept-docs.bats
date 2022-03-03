#!./test/libs/bats/bin/bats

load 'test_helper'

@test "[concept-docs] - collections.go" {
    runExample $CONCEPT_DOCS_DIR collections.go
    assert_success
    assert_output --partial "Done"
}

@test "[concept-docs] - buckets-and-clusters.go" {
    runExample $CONCEPT_DOCS_DIR buckets-and-clusters.go
    assert_success
}

@test "[concept-docs] - documents.go" {
    runExample $CONCEPT_DOCS_DIR documents.go
    assert_success
    assert_output --partial "Current value: 0"
    assert_output --partial "RESULT: 5"
}

@test "[concept-docs] - n1ql-query.go" {
    runExample $CONCEPT_DOCS_DIR n1ql-query.go
    assert_success
}

@test "[concept-docs] - xattr.go" {
    runExample $CONCEPT_DOCS_DIR xattr.go
    assert_success
}
