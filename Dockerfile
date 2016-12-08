FROM golang:1.7.4
MAINTAINER pedge@uber.com

EXPOSE 8080
ENV LOG_APP_NAME sally
ENV GO_TEMPLATE /go/src/go.uber.org/sally/etc/template/go.html
ENV INDEX_TEMPLATE /go/src/go.uber.org/sally/etc/template/index.html
ENV CONFIG /go/src/go.uber.org/sally/etc/config/config.yaml
RUN \
  curl -sSL https://get.docker.com/builds/Linux/x86_64/docker-1.12.3 > /bin/docker && \
  chmod +x /bin/docker
RUN \
  go get -v \
    github.com/golang/lint/golint \
    github.com/kisielk/errcheck \
    github.com/Masterminds/glide
RUN mkdir -p /go/src/go.uber.org/sally
ADD glide.yaml /go/src/go.uber.org/sally/
ADD glide.lock /go/src/go.uber.org/sally/
RUN glide install
ADD . /go/src/go.uber.org/sally/
WORKDIR /go/src/go.uber.org/sally
CMD ["make", "docker-dev-launch-internal"]
