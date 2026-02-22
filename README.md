# Subscriptions API

### Миграции

Для миграций использовался `golang-migrate`

```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Запустить миграции:

```
migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up
```