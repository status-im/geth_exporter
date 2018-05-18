# Build status-go in a Go builder container
FROM golang:1.10-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

RUN mkdir -p /go/src/github.com/status-im/geth_exporter
ADD . /go/src/github.com/status-im/geth_exporter
RUN cd /go/src/github.com/status-im/geth_exporter && make setup && make

# Copy the binary to the second image
FROM alpine:latest

RUN apk add --no-cache ca-certificates bash

COPY --from=builder /go/src/github.com/status-im/geth_exporter/build/bin/geth_exporter /usr/local/bin/

EXPOSE 9200
