.PHONY: setup lint-install dep-install lint

setup: dep-install lint-install ##@other Prepare project for first build

lint:
	@echo "lint"
	gometalinter ./...

lint-install:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

dep-install: ##@dependencies Install vendoring tool
	go get -u github.com/golang/dep/cmd/dep

test: lint
	go test -v ./...
