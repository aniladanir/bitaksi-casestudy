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

FROM debian:bullseye-slim

COPY --from=builder /app/app.yaml /etc/matching-api/config
COPY --from=builder /app/bin/matching-api /opt/app

ENTRYPOINT ["/opt/app/matching-api","--config","/etc/matching-api/config/app.yaml"]