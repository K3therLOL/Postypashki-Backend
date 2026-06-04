FROM golang:alpine

WORKDIR /app
RUN go install github.com/rubenv/sql-migrate/...@latest
COPY . .
