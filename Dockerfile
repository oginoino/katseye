# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder
WORKDIR /app

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY internal ./internal
COPY cmd ./cmd

RUN go build -ldflags="-s -w" -o /app/bin/katseye ./cmd/api

FROM alpine:3.20 AS runtime
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bin/katseye /usr/local/bin/katseye

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/katseye"]
