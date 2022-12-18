GOLINT = go run golang.org/x/lint/golint
STATICCHECK = go run honnef.co/go/tools/cmd/staticcheck

.PHONY: all
all: test

.PHONY: build
build:
	go build

.PHONY: install
install:
	go install .

.PHONY: lint
lint:
	$(GOLINT) ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	$(STATICCHECK) -tests=false ./...

.PHONY: pretest
pretest: lint vet staticcheck

.PHONY: test
test: pretest
	go test -race ./...

.PHONY: cover
cover:
	go test -coverprofile=cover.out -covermode=atomic -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: clean
clean:
	rm -rf _tmp

.PHONY: install
run: install
	sally
