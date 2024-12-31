FROM golang:1.22 as builder

WORKDIR /app

ARG GITHUB_ACCESS_TOKEN

ENV GO111MODULE on
ENV GOBIN=/usr/local/bin/go/bin

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./

RUN go build -o ./bin/driver-location-api ./cmd

FROM ubuntu:22.04

COPY --from=builder /app/coordinates.csv /etc/driver-location-api/data/coordinates.csv
COPY --from=builder /app/app.yaml /etc/driver-location-api/config/app.yaml
COPY --from=builder /app/bin/driver-location-api /opt/app/driver-location-api

ENTRYPOINT ["/opt/app/driver-location-api","--config","/etc/driver-location-api/config/app.yaml", "--coordinates", "/etc/driver-location-api/data/coordinates.csv"]