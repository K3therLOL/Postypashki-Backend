FROM golang:alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main main.go
RUN go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host localhost
CMD ["./main"]
