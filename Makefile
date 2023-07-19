MAKEFLAGS=--no-builtin-rules --no-builtin-variables --always-make
ROOT := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))
BIN_DIR:=$(ROOT)/bin

INTERNAL_TOKEN := ""
GO_DIST := grpc/proto/go
TS_DIST := grpc/web/src/rpc
PROTOC_GEN_TS_PATH := "grpc/web/node_modules/.bin/protoc-gen-ts"

vendor:
	go mod tidy

gen:
	@$(BIN_DIR)/gqlgen
	cp di/wire_gen.default.go di/wire_gen.go
	go generate di/wire_gen.go

gen-wire:
	cd di && $(BIN_DIR)/wire

gen-proto:
	mkdir -p $(GO_DIST)
	rm -f $(GO_DIST)/*
	protoc --proto_path=grpc/proto/. \
           --go-grpc_opt require_unimplemented_servers=false,paths=source_relative \
           --go-grpc_out $(GO_DIST) \
           --go_opt paths=source_relative \
           --go_out $(GO_DIST) \
           grpc/proto/*.proto

gen-proto-client:
	mkdir -p $(TS_DIST)
	rm -f $(TS_DIST)/*
	protoc --proto_path=grpc/proto/. \
    	   --plugin="protoc-gen-ts=$(PROTOC_GEN_TS_PATH)" \
           --js_out=import_style=commonjs,binary:$(TS_DIST) \
           --ts_out=service=grpc-web:$(TS_DIST) \
           grpc/proto/*.proto
	find grpc/web/src/rpc -type f -name "*_pb.js" | xargs gsed -i -e "1i /* eslint-disable */"
	find grpc/web/src/rpc -type f -name "*_pb_service.js" | xargs gsed -i -e "1i /* eslint-disable */"

build:
	GOOS=linux GOARCH=amd64 go build -o .tmp/main ./entrypoint/

run-local:
	go run cmd/merge_yaml/main.go app.yaml entrypoint/default/app.template.yaml env.yaml
	docker-compose up
	# IS_LOCAL=true GOOGLE_APPLICATION_CREDENTIALS=gcp-cred.json go run entrypoint/default/main.go

deploy-gae: vendor gen build
	go run cmd/merge_yaml/main.go app.yaml entrypoint/default/app.template.yaml env.yaml
	gcloud app deploy --quiet --version 1 --project canvas-329810 app.yaml
	go run cmd/merge_yaml/main.go app.yaml entrypoint/grpc/app.template.yaml env.yaml
	gcloud app deploy --quiet --version 1 --project canvas-329810 app.yaml

deploy-config:
	gcloud app deploy --quiet --project canvas-329810 index.yaml
	gcloud app deploy --quiet --project canvas-329810 cron.yaml
	gcloud app deploy --quiet --project canvas-329810 queue.yaml

clean-index:
	gcloud datastore indexes cleanup --quiet --project canvas-329810 index.yaml

deploy-functions:
	firebase deploy --only functions

deploy-functions-env:
	firebase functions:config:set token.internal=$(INTERNAL_TOKEN)

gen-gcp-credential-pem:
	openssl pkcs12 -in key.p12 -passin pass:notasecret -out key.pem -nodes
	cat key.pem | base64

install-tools:
	@GOBIN=$(BIN_DIR) go install github.com/99designs/gqlgen@v0.17.19
	@GOBIN=$(BIN_DIR) go install github.com/google/wire/cmd/wire@latest
