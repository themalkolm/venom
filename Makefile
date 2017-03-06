all: build

fmt:
	@gofmt -l -s -w -- *.go

build: fmt
	$(MAKE) -C _example

clean:
	$(MAKE) -C _example clean
