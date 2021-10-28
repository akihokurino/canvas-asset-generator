MAKEFLAGS=--no-builtin-rules --no-builtin-variables --always-make
ROOT := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))

vendor:
	go mod tidy

gen:
	gqlgen
	cp di/wire_gen.default.go di/wire_gen.go
	go generate di/wire_gen.go

build:
	GOOS=linux GOARCH=amd64 go build -o .tmp/main ./entrypoint/

run-local:
	docker-compose up

deploy-gae:
	gcloud app deploy --quiet --version 1 --project canvas-329810 app.yaml

deploy-index:
	gcloud app deploy --quiet index.yaml

deploy-functions:
	firebase deploy --only functions