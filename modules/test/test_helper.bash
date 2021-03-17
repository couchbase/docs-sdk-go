setup() {
    DEVGUIDE_DIR=../modules/devguide/examples/go
    HOWTOS_DIR=../modules/howtos/examples
    PROJECT_DOCS_DIR=../modules/project-docs/examples
    HELLO_WORLD_DIR=../modules/hello-world/examples
    
    load 'node_modules/bats-support/load'
    load 'node_modules/bats-assert/load'
}

function runExample() {
    cd $1
    run go run $2
}
