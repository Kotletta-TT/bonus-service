include .env

run:
	go run cmd/gophermart/main.go -d "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}"


accural:


build:
	go build -o cmd/gophermart/gophermart cmd/gophermart/main.go

db-init:
	docker run \
	--name bonus-service-db \
	-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
	-e POSTGRES_DB=${POSTGRES_DB} \
	-e POSTGRES_USER=${POSTGRES_USER} \
	-v /tmp/bonus-service/db_data:/var/lib/postgresql/data \
	-p 5432:5432 \
	-d \
	postgres:latest

db-down:
	docker container rm bonus-service-db

db-clean: db-down
	rm -rf /tmp/bonus-service/db_data

test: build
	./gophermarttest \
    -test.v -test.run=^TestGophermart$$ \
    -gophermart-binary-path=cmd/gophermart/gophermart \
    -gophermart-host=localhost \
    -gophermart-port=8080 \
    -gophermart-database-uri="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" \
	-accrual-binary-path=cmd/accrual/accrual_darwin_amd64 \
	-accrual-host=localhost \
	-accrual-database-uri="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" \
    -accrual-port=8082