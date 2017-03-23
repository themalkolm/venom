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

test: $(TESTS)
	@echo "----------------------------------------------------------------------"
	@echo ""
	@echo "Run $(words $?) tests: OK"

clean:
	rm -rf bin pkg
	rm -rf vendor/src/github.com/themalkolm/venom
