FROM golang:1.21.7 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOFLAGS="-mod=vendor"

ARG version="undefined"

WORKDIR /build

ADD . /build

RUN go build -o /app -ldflags "-X main.version=${version} -s -w"  ./cmd/main

CMD ["/app"]