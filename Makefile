all: build

fmt:
	@gofmt -l -s -w -- *.go _example/src/cmd/example/*.go

build: fmt
	$(MAKE) -C _example build

gotest: build
	$(MAKE) -C _example gotest

test: build
	$(MAKE) -C _example test

clean:
	$(MAKE) -C _example clean
