new_migration:
	migrate create -ext sql -dir ./migrations -seq $(name)

sqlc:
	sqlc generate -f ./db/sqlc.yaml

apigen:
	oapi-codegen --config=./configs/oapi.yaml api/task_api.yaml

gogen:
	go generate ./...

all: sqlc apigen gogen
	docker compose up -d
	go run ./cmd/main.go
