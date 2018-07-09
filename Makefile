.PHONY: setup lint-install dep-install lint build
.DEFAULT_GOAL := build

DOCKER_IMAGE_NAME ?= statusteam/geth_exporter

setup: dep-install lint-install ##@other Prepare project for first build

build: test
	go build -o build/bin/geth_exporter -v .

lint:
	@echo "lint"
	@golangci-lint run ./...

lint-install:
	@# The following installs a specific version of golangci-lint, which is appropriate for a CI server to avoid different results from build to build
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $(GOPATH)/bin v1.9.1

test: lint
	@echo "test"
	@go test -v ./...

dep-install: ##@dependencies Install vendoring tool
	go get -u github.com/golang/dep/cmd/dep

dep-ensure: ##@dependencies Dep ensure
	@dep ensure

docker-image:
	@echo "Building docker image..."
	docker build -t $(DOCKER_IMAGE_NAME):latest .
