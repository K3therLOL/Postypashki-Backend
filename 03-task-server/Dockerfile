FROM golang:latest AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/rubenv/sql-migrate/...@latest
COPY dbconfig.yml ./


