new_migration:
	migrate create -ext sql -dir ./migrations -seq $(name)

sqlc:
	sqlc generate -f ./db/sqlc.yaml

apigen:
	oapi-codegen --config=./configs/oapi.yaml api/task_api.yaml

gogen:
	go generate ./...

unit:
	go test -v ./...

test:
	go test -v -tags=integrations -count=1 ./...

.PHONY: test

lint:
	golangci-lint run

up:
	docker compose up -d --build

stop:
	docker compose stop

PROTO_DIR = api
PROTO_OUT = internal/grpc/pvz/v1
generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTO_OUT) \
		--go-grpc_out=$(PROTO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/pvz.proto

all: sqlc apigen gogen lint generate-proto unit
	docker compose up -d
