TEST?="unit"

default:
	$(MAKE) deps
	$(MAKE) all
deps:
	bash -c "./scripts/deps.sh"
test:
	bash -c "./scripts/test.sh $(TEST)"
check:
	$(MAKE) test
all:
	bash -c "./scripts/build.sh $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))"
