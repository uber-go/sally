export GOBIN = $(shell pwd)/bin
export PATH := $(GOBIN):$(PATH)

GOLINT = bin/golint
GO_FILES = $(shell find . \
	   -path '*/.*' -prune -o \
	   '(' -type f -a -name '*.go' ')' -print)

STATICCHECK = bin/staticcheck

TEST_FLAGS ?= -race

.PHONY: all
all: lint install test

.PHONY: lint
lint: gofmt golint staticcheck

.PHONY: gofmt
gofmt:
	$(eval FMT_LOG := $(shell mktemp -t gofmt.XXXXX))
	@gofmt -e -s -l $(GO_FILES) > $(FMT_LOG) || true
	@[ ! -s "$(FMT_LOG)" ] || \
		(echo "gofmt failed. Please reformat the following files:" | \
		cat - $(FMT_LOG) && false)

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

.PHONY: install
run: install
	sally
