.PHONY: install
install:
	glide --version || go get github.com/Masterminds/glide
	glide install


.PHONY: lint
lint:
	go vet .
	golint .


.PHONY: test
test: lint
	go test -race .


.PHONY: run
run:
	go build && ./sally
