FROM golang:1.12

ENV GO111MODULE on

COPY . /go/consul/

WORKDIR /go/consul

RUN go build -buildmode=c-shared -o consul.so .