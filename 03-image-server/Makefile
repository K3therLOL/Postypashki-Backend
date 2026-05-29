SHELL := /bin/bash
IMAGE = taskserver

all: up

cert:
	go run $$(go env GOROOT)/src/crypto/tls/generate_cert.go --host localhost

docker_run:
	export $$(cat .env | xargs) && envsubst < .env > .env.resolved
	docker run --env-file .env.resolved -it $(IMAGE) /bin/sh
	rm .env.resolved

.PHONY: clean
.ONESHELL:
up:
	docker compose up --build

down:
	docker compose down -v
