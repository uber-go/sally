PKGS := $(shell go list ./... | grep -v go.uber.org/sally/vendor)

all: test

vendor:
	rm -rf vendor
	go get -d -v -t -u -f ./...
	go get -v github.com/Masterminds/glide
	glide create
	glide update

build:
	go build $(PKGS)

install:
	go install $(PKGS)

lint:
	go get -v github.com/golang/lint/golint
	for file in $$(find . -name '*.go' | grep -v '^\./vendor'); do \
		golint $$file; \
		if [ -n "$$(golint $$file)" ]; then \
			exit 1; \
		fi; \
	done

vet:
	go vet $(PKGS)

errcheck:
	go get -v github.com/kisielk/errcheck
	errcheck $(PKGS)

pretest: lint vet errcheck

test: pretest
	go test -race $(PKGS)

clean:
	go clean -i $(PKGS)
	rm -rf _tmp

docker-build-dev:
	docker build -t uber/sally-dev .

docker-test: docker-build-dev
	docker run uber/sally-dev make test

docker-build-internal:
	rm -rf _tmp
	mkdir -p _tmp
	CGO_ENABLED=0 go build -a -installsuffix cgo -o _tmp/sally cmd/sally/main.go
	docker build -t uber/sally -f Dockerfile.sally .

docker-build: docker-build-dev
	docker run -v /var/run/docker.sock:/var/run/docker.sock uber/sally-dev make docker-build-internal

docker-launch-dev-internal: install
	sally

docker-launch-dev: docker-build-dev
	docker run -p 8080:8080 uber/sally-dev

docker-launch: docker-build
	docker run -p 8080:8080 uber/sally

launch:
	GO_TEMPLATE=etc/template/go.html INDEX_TEMPLATE=etc/template/index.html CONFIG=etc/config/config.yaml go run cmd/sally/main.go

.PHONY: \
	all \
	vendor \
	build \
	install \
	lint \
	vet \
	errcheck \
	pretest \
	test \
	clean \
	docker-build-dev \
	docker-test \
	docker-build-internal \
	docker-build \
	docker-launch-dev-internal \
	docker-launch-dev \
	docker-launch \
	launch
