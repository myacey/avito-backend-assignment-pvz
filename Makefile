new_migration:
	migrate create -ext sql -dir ./migrations -seq $(name)

sqlc:
	sqlc generate

compose_up:
	docker compose -p avito-pvz -f ./deployments/docker-compose.yaml up --build