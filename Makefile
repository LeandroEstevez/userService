DB_URL=postgresql://root:userMicroServiceDB@usermicroservicedb.cviqqzopm7zr.us-east-2.rds.amazonaws.com:5432/userMicroServiceDB

newPostgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=userMicroServiceDB -d postgres:latest

postgres:
	docker start postgresUser

createdb:
	docker exec -it postgresUser createdb --username=root --owner=root userMicroServiceDB

dropdb:
	docker exec -it postgresUser dropdb userMicroServiceDB

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination userMicroService/db/mock/store.go userMicroService/db/sqlc Store

.PHONY: network newPostgres postgres createdb dropdb migrateup migratedown sqlc server mock