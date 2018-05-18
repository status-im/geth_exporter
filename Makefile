.PHONY: setup lint-install dep-install lint build
.DEFAULT_GOAL := build

setup: dep-install lint-install ##@other Prepare project for first build

build: test
	go build

lint:
	@echo "lint"
	@gometalinter ./...

lint-install:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

test: lint
	@echo "test"
	@go test -v ./...

dep-install: ##@dependencies Install vendoring tool
	go get -u github.com/golang/dep/cmd/dep

dep-ensure: ##@dependencies Dep ensure
	@dep ensure
