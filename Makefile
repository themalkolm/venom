all: build

fmt:
	@gofmt -l -s -w -- *.go

build: fmt
	$(MAKE) -C _example build

test:
	$(MAKE) -C _example test

clean:
	$(MAKE) -C _example clean
