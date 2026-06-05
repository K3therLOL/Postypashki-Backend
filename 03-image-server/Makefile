SHELL := /bin/bash
IMAGE = taskserver

all: up

cert:
	go run $$(go env GOROOT)/src/crypto/tls/generate_cert.go --host localhost

.PHONY: up down api consumer
.ONESHELL:
up:
	docker compose up --build

down:
	docker compose down -v

api:
	docker compose up --build -d api

consumer:
	docker compose restart consumer
