FROM golang:1.7.4

EXPOSE 8080
RUN \
  curl -fsSLO https://get.docker.com/builds/Linux/x86_64/docker-latest.tgz && \
  tar --strip-components=1 -xvzf docker-latest.tgz -C /usr/local/bin
RUN mkdir -p /go/src/go.uber.org/sally
WORKDIR /go/src/go.uber.org/sally
ADD glide.yaml glide.lock /go/src/go.uber.org/sally/
RUN go get -v github.com/Masterminds/glide && glide install
RUN go get -v github.com/golang/lint/golint github.com/kisielk/errcheck honnef.co/go/staticcheck/cmd/staticcheck
ADD . /go/src/go.uber.org/sally/
CMD ["make", "docker-launch-dev-internal"]
