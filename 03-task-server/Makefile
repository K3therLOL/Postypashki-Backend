SHELL := /bin/bash
IMAGE = taskserver

all: run

cert:
	go run $$(go env GOROOT)/src/crypto/tls/generate_cert.go --host localhost

docker_run:
	docker run --env-file .env $(IMAGE)

.ONESHELL:
run:
	set -a; source .env; set +a; 
	sql-migrate up
	go run main.go

down:
	set -a; source .env; set +a; 
	sql-migrate down
