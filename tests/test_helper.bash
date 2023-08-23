setup() {
    DEVGUIDE_DIR=../modules/devguide/examples/go
    HOWTOS_DIR=../modules/howtos/examples
    HELLO_WORLD_DIR=../modules/hello-world/examples
    CONCEPT_DOCS_DIR=../modules/concept-docs/examples

    BATS_TEST_RETRIES=3

    load 'node_modules/bats-support/load'
    load 'node_modules/bats-assert/load'
}

function runExample() {
    cd $1
    run go run $2
}