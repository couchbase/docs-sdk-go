.PHONY: build-all-examples

build-all-examples: build-docs build-devguide build-hello-world build-howtos build-project-docs 

build-docs:
	$(info ************ BUILDING DOC EXAMPLES ************)

build-hello-world:
	@./build-docs.sh modules/hello-world/examples

build-howtos:
	@./build-docs.sh modules/howtos/examples

build-project-docs:
	@./build-docs.sh modules/project-docs/examples

build-devguide:
	@./build-docs.sh modules/devguide/examples/go
