DB_URL=postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable
postgres:
	docker run --network bank-network -p 5433:5432 --name postgres16 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine 
createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres16 dropdb simple_bank
migrateup:
	migrate --path db/migration --database "$(DB_URL)" --verbose up
migrateup1:
	migrate --path db/migration --database "$(DB_URL)" --verbose up 1
migratedown:
	migrate --path db/migration --database "$(DB_URL)" --verbose down
migratedown1:
	migrate --path db/migration --database "$(DB_URL)" --verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/anewgd/simple-bank/db/sqlc Store  
db_docs:
	dbdocs build doc/db.dbml
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
    proto/*.proto
evans:
	evans --host localhost --port 8002 -r repl
.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock db_docs db_schema proto evans