# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-alpine AS build

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . /app/

RUN go build -o /app/front-tracing-go

## Deploy
FROM alpine

RUN mkdir -p /app
WORKDIR /

COPY --from=build /app/front-tracing-go /front-tracing-go
COPY --from=build /app/config.yaml /config.yaml

EXPOSE 3000

ENTRYPOINT ["/front-tracing-go"]