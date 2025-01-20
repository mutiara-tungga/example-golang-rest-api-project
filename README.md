# golang-rest-api

## Requirements
- [Git](https://git-scm.com/downloads)
- [Go v1.22+](https://go.dev/dl)
- [PostgreSql](https://www.postgresql.org/download/) or [docker](https://hub.docker.com/_/postgres)
- [Golang Migrate](https://github.com/golang-migrate/migrate/tree/master)
- [swagger](https://github.com/swaggo/swag/) `go install github.com/swaggo/swag/cmd/swag@v1.16.4`
- docker and docker-compose

## How to Run
- Install swagger
- swag init -g cmd/api/main.go
- go run cmd/api/main.go

## How to use postgres docker
- run `cp docker/env.example docker/.env`
- fill all variable on docker/.env based on your preference
- run `docker compose -f docker/docker-compose.yml -p golang-rest-api-infra up -d`
- get docker container id `docker ps -a`
- run docker image

## How to run database migration
- create new migration : `make new_migration name={migration_name}`
- run migration:
  - `export DB_URL=postgresql://{db_user}:{db_password}@{db_host}:{db_port}/{db_name}?sslmode=disable`
  - `make migrateup`

## How to seed user for the first time
- run `cp script/seed_user/env.example script/seed_user/.env`
- update the .env
- fill the csv
- run `go run script/seed_user/seed_user.go`


## How to generate RSA key for JWT
- run generate private key first `openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048`
- run generate public key `openssl rsa -pubout -in private_key.pem -out public_key.pem`
