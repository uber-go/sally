.PHONY: install
install:
	glide --version || go get github.com/Masterminds/glide
	glide install
	go get -u github.com/jteeuwen/go-bindata/...


.PHONY: run
run:
	go install . && sally
