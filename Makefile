DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name simplebank-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it simplebank-postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simplebank-postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	docker run --rm -v "$(CURDIR):/src" -w //src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server
