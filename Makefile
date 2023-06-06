# DB_URL=postgresql://couple_admin:admin_password@localhost:5432/couple_db?sslmode=disable
DB_URL=postgresql://root:root@localhost:5432/couple_db?sslmode=disable

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

server:
	go run main.go	

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/leemingen1227/couple-server/db/sqlc Store

.PHONY: migrateup migrateup1 migratedown migratedown1 sqlc server mock

