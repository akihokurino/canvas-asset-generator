MAKEFLAGS=--no-builtin-rules --no-builtin-variables --always-make
ROOT := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))

deploy-functions:
	firebase deploy --only functions

vendor:
	go mod tidy

deploy-gae:
	gcloud app deploy --quiet --version 1 --project canvas-329810 app.yaml