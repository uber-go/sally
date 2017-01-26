PKGS := $(shell go list ./... | grep -v go.uber.org/sally/vendor)
SRCS := $(wildcard *.go)

.PHONY: all
all: test

.PHONY: vendor-update
vendor-update:
	go get -v github.com/Masterminds/glide
	glide update

.PHONY: vendor-install
vendor-install:
	go get -v github.com/Masterminds/glide
	glide install

.PHONY: build
build:
	go build $(PKGS)

.PHONY: install
install:
	go install $(PKGS)

.PHONY: lint
lint:
	go install ./vendor/github.com/golang/lint/golint
	for file in $(SRCS); do \
		golint $$file; \
		if [ -n "$$(golint $$file)" ]; then \
			exit 1; \
		fi; \
	done

.PHONY: vet
vet:
	go vet $(PKGS)

.PHONY: errcheck
errcheck:
	go install ./vendor/github.com/kisielk/errcheck
	errcheck $(PKGS)

.PHONY: staticcheck
staticcheck:
	go install ./vendor/honnef.co/go/staticcheck/cmd/staticcheck
	staticcheck $(PKGS)

.PHONY: pretest
pretest: lint vet errcheck staticcheck

.PHONY: test
test: pretest
	go test -race $(PKGS)

.PHONY: clean
clean:
	go clean -i $(PKGS)
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
	CGO_ENABLED=0 go build -a -installsuffix cgo -o _tmp/sally $(SRCS)
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
