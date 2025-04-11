new_migration:
	migrate create -ext sql -dir ./migrations -seq $(name)

sqlc:
	sqlc generate -f ./db/sqlc.yaml

codegen:
	oapi-codegen --config=./configs/oapi.yaml api/task_api.yaml

all: sqlc
	docker compose up -d
	go run ./cmd/main.go
