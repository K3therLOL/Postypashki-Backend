SHELL := /bin/bash

all: run

cert:
	go run $$(go env GOROOT)/src/crypto/tls/generate_cert.go --host localhost

run:
	set -a; source .env; set +a; go run main.go
