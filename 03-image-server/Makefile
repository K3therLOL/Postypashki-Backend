SHELL := /bin/bash
IMAGE = taskserver

all: run

cert:
	go run $$(go env GOROOT)/src/crypto/tls/generate_cert.go --host localhost

docker_run:
	export $$(cat .env | xargs) && envsubst < .env > .env.resolved
	docker run --env-file .env.resolved -it $(IMAGE) /bin/sh
	rm .env.resolved

.ONESHELL:
run:
	docker compose up
	#set -a; source .env; set +a; 
	#sql-migrate up
	#go run main.go

clean:
	docker compose down -v

down:
	set -a; source .env; set +a; 
	sql-migrate down
