# Go Migrate Cheat Sheet

## Source
https://github.com/golang-migrate/migrate/

## Установка
```sh
# Установить go-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Создание новой миграции
```sh
migrate create -seq -ext=.sql -dir=./migrations name
```

### В этой команде:
- Флаг -seq указывает, что мы хотим использовать последовательную нумерацию типа 0001, 0002, ...
  для файлов миграции (вместо временной метки Unix, которая используется по умолчанию).
- Флаг -ext указывает, что мы хотим присвоить файлам миграции расширение .sql.
- Флаг -dir указывает, что мы хотим хранить файлы миграции в каталоге ./migrations
  (который будет создан автоматически, если он еще не существует).


## Применение миграций
```sh
migrate -path=./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up
```

## Откат последней миграции
```sh
migrate -path=./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down 1
```

## Откат всех миграций
```sh
migrate -path=./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down
```

## Просмотр статуса миграций
```sh
migrate -path=./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" version
```

## Пропуск миграций без их выполнения
```sh
migrate -path=./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" force <версия>
```

## Применение определенного количества миграций
```sh
migrate -path=./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up 2
```

psql --host=localhost --dbname=orderservice --username=temut