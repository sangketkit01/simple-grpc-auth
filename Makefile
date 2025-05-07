DB_URL=postgres://root:secret@localhost:5433/simple_grpc_auth?sslmode=disable

server:
	go run main.go

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_grpc_auth

dropdb:
	docker exec -it postgres dropdb --username=root --owner=root simple_grpc_auth

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc \
	--proto_path=proto \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	proto/*.proto

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc: 
	sqlc generate

.PHONY: proto new_migration sqlc serer migrateup migratedown createdb dropdb