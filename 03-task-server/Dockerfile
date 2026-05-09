FROM golang:alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main

FROM golang:alpine
WORKDIR /app
RUN go install github.com/rubenv/sql-migrate/...@latest
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/dbconfig.yml .
COPY --from=builder /build/main .
