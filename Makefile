postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up


migrateupaws:
	migrate -path db/migration -database "postgresql://root:secret@simplebank.cucmeqddxrhf.us-east-2.rds.amazonaws.com:5432/simple_bank" -verbose up
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/techschool/simple-bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateupaws