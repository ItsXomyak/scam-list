compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f


compose-down: ### Down docker-compose
	docker-compose down --remove-orphans

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci


## Database
DB_URL=postgres://scam-user:scam-password@localhost:5432/scam-db?sslmode=disable

## Создание новой миграции: make migrate-create name=название
migrate-create:
	@echo "Creating new migration: $(name)"
	migrate create -seq -ext=.sql -dir=./migrations/postgres $(name)

## Применить все миграции
migrate-up:
	migrate -path=./migrations/postgres -database "$(DB_URL)" up

## Применить N миграций: make migrate-upn n=2
migrate-upn:
	migrate -path=./migrations/postgres -database "$(DB_URL)" up $(n)

## Откатить одну миграцию
migrate-down1:
	migrate -path=./migrations/postgres -database "$(DB_URL)" down 1

## Откатить все миграции
migrate-down:
	migrate -path=./migrations/postgres -database "$(DB_URL)" down

## Посмотреть текущую версию миграций
migrate-version:
	migrate -path=./migrations/postgres -database "$(DB_URL)" version


test: ### run test
	go test -v ./...

coverage-html: ### run test with coverage and open html report
	go test -coverprofile=cvr.out ./...
	go tool cover -html=cvr.out
	rm cvr.out

coverage: ### run test with coverage
	go test -coverprofile=cvr.out ./...
	go tool cover -func=cvr.out
	rm cvr.out
