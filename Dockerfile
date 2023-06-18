FROM golang:1.19.3-alpine AS build
WORKDIR /go/src

ENV GO111MODULE=on
COPY ./jkojs-worker/go.mod ./jkojs-worker/go.sum ./
RUN go mod download

COPY ./jkojs-worker/ /jkojs-worker
WORKDIR /jkojs-worker

RUN CGO_ENABLED=0 GOOS=linux go build -o jkworker -ldflags "-s -w" && cp ./jkworker /jkworkers

FROM ubuntu:20.04 AS run

ENV DEBIAN_FRONTEND nointeractive

RUN apt-get update && apt-get install -y  build-essential gcc clang golang-go --no-install-recommends && mkdir /work && chmod 777 /work && mkdir /built && chmod 777 /built && touch out.json && chmod 666 out.json && mkdir /home/worker

COPY --from=build /jkworker /jkworker

RUN useradd worker
USER worker