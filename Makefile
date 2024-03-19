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
new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)
sqlc:
	sqlc generate
test:
	go test -v -cover -short ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/anewgd/simple-bank/db/sqlc Store  
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/anewgd/simple-bank/worker TaskDistributor  


db_docs:
	dbdocs build doc/db.dbml
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple-bank \
    proto/*.proto
	statik --src=./doc/swagger --dest=./doc
evans:
	evans --host localhost --port 8002 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7.2.4-alpine
.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock db_docs db_schema proto evans redis new_migration