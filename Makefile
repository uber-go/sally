export GOBIN = $(shell pwd)/bin
export PATH := $(GOBIN):$(PATH)

GOLINT = bin/golint
STATICCHECK = bin/staticcheck

TEST_FLAGS ?= -race

.PHONY: all
all: lint install test

.PHONY: lint
lint: golint staticcheck

.PHONY: staticcheck
staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./...

$(STATICCHECK): tools/go.mod
	cd tools && go install honnef.co/go/tools/cmd/staticcheck

.PHONY: golint
golint: $(GOLINT)
	$(GOLINT) ./...

$(GOLINT): tools/go.mod
	cd tools && go install golang.org/x/lint/golint

.PHONY: install
install:
	go install .

.PHONY: test
test:
	go test $(TEST_FLAGS) ./...

.PHONY: cover
cover:
	go test $(TEST_FLAGS) -coverprofile=cover.out -covermode=atomic -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: clean
clean:
	rm -rf _tmp

.PHONY: docker-build-dev
docker-build-dev:
	docker build -t uber/sally-dev .

.PHONY: docker-test
docker-test: docker-build-dev
	docker run uber/sally-dev make test

.PHONY: docker-build-internal
docker-build-internal:
	rm -rf _tmp
	mkdir -p _tmp
	CGO_ENABLED=0 go build -a -o _tmp/sally .
	docker build -t uber/sally -f Dockerfile.scratch .

.PHONY: docker-build
docker-build: docker-build-dev
	docker run -v /var/run/docker.sock:/var/run/docker.sock uber/sally-dev make docker-build-internal

.PHONY: docker-launch-dev
docker-launch-dev: docker-build-dev
	docker run -p 8080:8080 uber/sally-dev

.PHONY: docker-launch
docker-launch: docker-build
	docker run -p 8080:8080 uber/sally

.PHONY: install
run: install
	sally
