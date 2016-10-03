PACKAGES := $(shell glide novendor)

.PHONY: install
install:
	glide --version || go get github.com/Masterminds/glide
	glide install


.PHONY: lint
lint:
	go vet $(PACKAGES)
	golint .


.PHONY: test
test: lint
	go test -race $(PACKAGES)


.PHONY: run
run:
	go build && ./sally
