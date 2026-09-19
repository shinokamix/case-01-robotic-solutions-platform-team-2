# Инструкции для `apps/api`

Запускайте команды из `apps/api`:

```bash
make dev
make generate
make build
make test
make lint
make check
make migrate
```

## Границы

- Храните бизнес-логику в `internal/<module>`.
- HTTP-типы держите в `internal/httpapi`.
- Изменяйте публичный HTTP-контракт сначала в `packages/contracts/openapi.yaml`.
- Новую схему добавляйте отдельной миграцией в `db/migrations`. Не изменяйте примененные миграции.

## Сгенерированные файлы

Не редактируйте вручную:

- `internal/httpapi/generated/openapi.gen.go`
- `internal/postgres/sqlc/*.go`

После изменений OpenAPI, миграций или SQL-запросов запускайте `make generate`.
