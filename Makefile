.PHONY: install
install:
	glide --version || go get github.com/Masterminds/glide
	glide install


.PHONY: test
test:
	go test .


.PHONY: run
run:
	go build && ./sally
