all: build

gb:
	@which gb 2>/dev/null || go get github.com/constabulary/gb/...

fmt:
	@gofmt -l -s -w -- src

build: fmt gb
	@mkdir -p       vendor/src/github.com/themalkolm/venom
	@cp -rv ../*.go vendor/src/github.com/themalkolm/venom
	gb build ...

clean:
	rm -rf bin pkg
