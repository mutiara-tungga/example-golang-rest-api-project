new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down $(step)

run-rest-api:
	go run cmd/api/main.go

swaggo:
	swag init -g cmd/api/main.go --pd