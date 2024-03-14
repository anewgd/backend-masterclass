postgres:
	docker run --network bank-network -p 5433:5432 --name postgres16 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine 
createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres16 dropdb simple_bank
migrateup:
	migrate --path db/migration --database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" --verbose up
migrateup1:
	migrate --path db/migration --database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" --verbose up 1
migratedown:
	migrate --path db/migration --database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" --verbose down
migratedown1:
	migrate --path db/migration --database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" --verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/anewgd/simple-bank/db/sqlc Store  
.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock