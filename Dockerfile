# This image provides the sally binary.
# It does not include a sally configuration.
# A /sally.yaml file is required for this to run.

FROM golang:1.22-alpine

COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o sally go.uber.org/sally

FROM scratch
COPY --from=0 /build/sally /sally
EXPOSE 8080
WORKDIR /
CMD ["/sally"]
