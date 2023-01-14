DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
proto:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
           --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
           --grpc-gateway_out=pb  --grpc-gateway_opt=paths=source_relative \
           --openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank\
           proto/*.proto
	statik -src=./docs/swagger -dest=./docs

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up


migrateupaws:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...

server:
	go run main.go

db_docs:
	dbdocs build doc/db.dbml


db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

mock:
	mockgen -build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/techschool/simple-bank/db/sqlc Store

evans:
	evans --host localhost --port 9090 -r redis:7-alpine

redis:
	docker run --name redis -p 6379:6379 -d redis:7.0.7-alpine

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test \
 server mock migrateupaws db_docs db_schema proto evans redis