TESTS := $(wildcard tests/*)

all: build test

gb:
	@which gb 2>/dev/null || go get github.com/constabulary/gb/...

fmt:
	@gofmt -l -s -w -- src

build: fmt gb
	@mkdir -p       vendor/src/github.com/themalkolm/venom
	@cp -rv ../*.go vendor/src/github.com/themalkolm/venom
	gb build ...

.PHONY: $(TESTS)
$(TESTS):
	@echo "----------------------------------------------------------------------"
	@echo "Running test: $@ ..."
	@echo "----------------------------------------------------------------------"
	$@

expecttest: $(TESTS)
	@echo "----------------------------------------------------------------------"
	@echo ""
	@echo "Run $(words $?) tests: OK"

.PHONY: gotest
gotest:
	gb test -v ...
	gb test -v github.com/themalkolm/venom

.PHONY: test
test: gotest expecttest

clean:
	rm -rf bin pkg
	rm -rf vendor/bin vendor/pkg
	rm -rf vendor/src/github.com/themalkolm/venom
