PROJECT := offheap
PACKAGE := github.com/remerge/$(PROJECT)

# http://stackoverflow.com/questions/322936/common-gnu-makefile-directory-path#comment11704496_324782
TOP := $(dir $(CURDIR)/$(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST)))

GOFMT=gofmt -w -s

GOSRCDIR=$(GOPATH)/src/$(PACKAGE)
GOPATHS=$(shell glide novendor)
GOFILES=$(shell git ls-files | grep '\.go$$')
MAINGO=$(wildcard main/*.go)
MAIN=$(patsubst main/%.go,%,$(MAINGO))

.PHONY: build run clean lint test bench fmt dep init up gen release deploy

all: build

build: fmt
	cd $(GOSRCDIR) && \
		CGO_ENABLED=0 \
		go build $(MAINGO)

run: build
	./$(MAIN)

clean:
	go clean
	rm -f $(MAIN)
	rm -rf $(TOP)/vendor/

lint:
	cd $(GOSRCDIR) && \
		gometalinter --vendor --errors --fast --deadline=60s -D gotype -D vet $(GOPATHS)

test: build lint
	cd $(GOSRCDIR) && \
		go test -timeout 60s $(GOPATHS)

bench:
	cd $(GOSRCDIR) && \
		go test -bench=. -cpu 4 $(GOPATHS)

fmt:
	$(GOFMT) $(GOFILES)

dep:
	go get -u github.com/Masterminds/glide
	go get -u github.com/smartystreets/goconvey
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update
	cd $(GOSRCDIR) && glide install

init:
	cd $(GOSRCDIR) && \
		glide init

up:
	cd $(GOSRCDIR) && \
		glide update

gen:
	cd $(GOSRCDIR) && \
		go generate $(GOPATHS)
	$(GOFMT) $(GOFILES)

release:
	git push origin master master:production
