# Инструкции для `apps/api`

Запускайте команды из корня репозитория через Moon:

```bash
moon run api:dev
moon run api:generate
moon run api:build
moon run api:test
moon run api:lint
moon run api:check
moon run api:migrate
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

После изменений OpenAPI, миграций или SQL-запросов запускайте `moon run api:generate`.
