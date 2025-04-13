new_migration:
	migrate create -ext sql -dir ./migrations -seq $(name)

sqlc:
	sqlc generate -f ./db/sqlc.yaml

apigen:
	oapi-codegen --config=./configs/oapi.yaml api/task_api.yaml

gogen:
	go generate ./...

coverage:
	@go test -coverprofile=coverage.out ./... && \
	grep -v -E '/(mocks|sqlc|openapi|pkg)/' coverage.out > filtered.out && \
	COVERAGE=$$(go tool cover -func=filtered.out | grep total | awk '{print $$3}') && \
	echo "--------------------------------" && \
	echo "Coverage: $${COVERAGE}" && \
	rm coverage.out filtered.out

all: sqlc apigen gogen
	docker compose up -d
	go run ./cmd/main.go
