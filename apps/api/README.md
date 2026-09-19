# Go API

## Запуск

Скопируйте `.env.example` в окружение и запустите PostgreSQL:

```bash
make migrate
make dev
```

## Генерация

```bash
make generate
```

Команда генерирует Go-код по `packages/contracts/openapi.yaml`, SQL-запросам и миграциям. Не редактируйте файлы в `internal/httpapi/generated` и `internal/postgres/sqlc` вручную.

Контракт API хранится в `packages/contracts/openapi.yaml`.
