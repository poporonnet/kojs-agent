FROM golang:1.19.3-alpine AS Builder
WORKDIR /go/src

ENV GO111MODULE=on
COPY jkojs-worker/go.mod jkojs-worker/go.sum ./
COPY jkojs-worker ./jkojs-worker
WORKDIR ./jkojs-worker
RUN CGO_ENABLED=0 GOOS=linux go build -o ojs-worker -ldflags "-s -w"

FROM debian:bullseye AS Runner

RUN <<EOF
useradd worker
mkdir -p /home/worker

mkdir /home/worker/work
mkdir /home/worker/built
touch /home/worker/out.json
chmod -R 777 /home/worker
chown -R worker:worker /home/worker
EOF

WORKDIR /home/worker
COPY --from=Builder /go/src/jkojs-worker/ojs-worker ./ojs-worker
RUN <<EOF
apt-get update && apt-get install -y build-essential clang
EOF
USER worker
