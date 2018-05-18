.PHONY: setup lint-install dep-install lint build
.DEFAULT_GOAL := build

DOCKER_IMAGE_NAME ?= statusteam/geth_exporter

setup: dep-install lint-install ##@other Prepare project for first build

build: test
	go build -o build/bin/geth_exporter -v .

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

docker-image:
	@echo "Building docker image..."
	docker build --file _assets/Dockerfile -t $(DOCKER_IMAGE_NAME):latest .
