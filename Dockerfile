FROM golang:1.8.1

EXPOSE 8080
RUN \
  curl -fsSLO https://get.docker.com/builds/Linux/x86_64/docker-latest.tgz && \
  tar --strip-components=1 -xvzf docker-latest.tgz -C /usr/local/bin
WORKDIR /go/src/go.uber.org/sally
ADD glide.yaml glide.lock /go/src/go.uber.org/sally/
RUN go get -v github.com/Masterminds/glide && glide install
ADD . /go/src/go.uber.org/sally/
CMD ["make", "run"]
