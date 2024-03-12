postgres:
	docker run -p 5433:5432 --name postgres16 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine 
createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres16 dropdb simple_bank
migrateup:
	migrate --path db/migration --database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" --verbose up
migratedown:
	migrate --path db/migration --database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" --verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test