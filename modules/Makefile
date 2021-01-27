.PHONY: build

GO_VERSION ?="1.13"
GOCB_VERSION ?="v2.2.1"

CB_IMAGE = build-docker.couchbase.com/smallcb-sdks
CB_VERSION = cb7-sdk3-ub20

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

# Run couchbase server+sdks container. Note that this runs with the `-rm` option, 
# which will ensure the container is deleted when stopped.
cb-run:
	@docker run -it --rm -v ${PWD}:/modules -d --name cb-test -p 8091-8096:8091-8096 ${CB_IMAGE}:${CB_VERSION}
	@docker exec cb-test bin/bash -c "cp -r /modules/test/scripts/init-couchbase ."
	@docker exec cb-test bin/bash -c "/init-couchbase/init.sh"
	@docker exec cb-test bin/bash -c "/init-couchbase/init-buckets.sh"

# Run tests
tests:
	@docker exec -it cb-test bin/bash -c "cd modules && bats ./test"

# Run a single test.
# i.e make TEST_NAME="analytics-named-placeholders.go" single-test
single-test:
	@docker exec -it cb-test bin/bash -c "cd modules && bats ./test -f ${TEST_NAME}"
	
# Stop the server container
cb-stop:
	@docker stop cb-test
