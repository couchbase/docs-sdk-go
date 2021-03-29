.PHONY: build

GO_VERSION ?="1.13"
GOCB_VERSION ?="v2.2.1"

CB_IMAGE=couchbase
CB_VERSION=cb7-sdk3-go-ub20

TEST_NAME=""

# --------------------------
# BUILD
# --------------------------
build:
	$(info ************ BUILDING DOC EXAMPLES ************)
	@./test/scripts/build-examples.sh hello-world/examples
	@./test/scripts/build-examples.sh howtos/examples
	@./test/scripts/build-examples.sh project-docs/examples
	@./test/scripts/build-examples.sh devguide/examples/go

# Utility script to update the go and gocb versions for all code examples.
# Avoids having to manually change each go.mod and go.sum file in every examples directory.
update-versions:
	$(info ************ UPDATING DOC EXAMPLE VERSIONS ************)
	@./test/scripts/update-versions.sh devguide/examples/go ${GO_VERSION} ${GOCB_VERSION}
	@./test/scripts/update-versions.sh hello-world/examples ${GO_VERSION} ${GOCB_VERSION}
	@./test/scripts/update-versions.sh howtos/examples ${GO_VERSION} ${GOCB_VERSION}
	@./test/scripts/update-versions.sh project-docs/examples ${GO_VERSION} ${GOCB_VERSION}


# -------------------------------------------
# DOCKER
# -------------------------------------------
cb-build:
	@docker build --network host -t ${CB_IMAGE}:${CB_VERSION} -f test/Dockerfile .

# Run couchbase server+sdk container. Note that this runs with the `-rm` option, 
# which will ensure the container is deleted when stopped.
cb-start:
	@docker run -t --rm -v ${PWD}:/modules -d --name cb-test -p 8091-8096:8091-8096 ${CB_IMAGE}:${CB_VERSION}
	@docker exec -t cb-test bin/bash -c "/init-couchbase/init.sh"
	@docker exec -t cb-test bin/bash -c "/init-couchbase/init-buckets.sh"

# Run all tests
tests:
	@cd test && npm install
	@docker exec -t cb-test bin/bash -c "cd modules && bats -T ./test"

# Run a single test.
# i.e make TEST_NAME="analytics-named-placeholders.go" single-test
single-test:
	@cd test && npm install
	@docker exec -t cb-test bin/bash -c "cd modules && bats -T ./test -f ${TEST_NAME}"
	
# Stop the server container
cb-stop:
	@docker stop cb-test