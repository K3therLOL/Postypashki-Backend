FROM golang:latest AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o cryptoserver

FROM golang:latest

WORKDIR /app

COPY --from=build /build .

CMD ["./cryptoserver"]
