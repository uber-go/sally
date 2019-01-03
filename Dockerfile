FROM golang:1.11

EXPOSE 8080
RUN \
  curl -fsSLO https://get.docker.com/builds/Linux/x86_64/docker-latest.tgz && \
  tar --strip-components=1 -xvzf docker-latest.tgz -C /usr/local/bin
ENV GO111MODULE=on
RUN mkdir -p /go/src/go.uber.org/sally
WORKDIR /go/src/go.uber.org/sally
ADD . /go/src/go.uber.org/sally/
RUN go mod vendor
CMD ["make", "run"]
