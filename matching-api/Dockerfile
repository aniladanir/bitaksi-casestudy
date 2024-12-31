ARG GO_PATH
FROM golang:1.22 as builder

WORKDIR /app

ARG GITHUB_ACCESS_TOKEN

ENV GO111MODULE on
ENV GOBIN=/usr/local/bin/go/bin

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./

RUN go build -o ./bin/matching-api ./cmd/

FROM ubuntu:22.04

COPY --from=builder /app/app.yaml /etc/matching-api/config/app.yaml
COPY --from=builder /app/bin/matching-api /opt/app/matching-api

ENTRYPOINT ["/opt/app/matching-api","--config","/etc/matching-api/config/app.yaml"]