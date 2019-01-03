GOLINT = go run github.com/golang/lint/golint
STATICCHECK = go run honnef.co/go/tools/cmd/staticcheck
GOBINDATA = go run github.com/go-bindata/go-bindata/go-bindata

.PHONY: all
all: test

.PHONY: build
build:
	go build

.PHONY: generate
generate: bindata.go

bindata.go: templates/*
	$(GOBINDATA) templates

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
