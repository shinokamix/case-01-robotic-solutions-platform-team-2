# Go API

## Запуск

Запустите локальное окружение из корня репозитория:

```bash
moon run root:dev
```

Для отдельной работы с API используйте задачи Moon:

```bash
moon run api:dev
moon run api:migrate
```

## Генерация

```bash
moon run api:generate
```

Команда генерирует Go-код по `packages/contracts/openapi.yaml`, SQL-запросам и миграциям. Не редактируйте файлы в `internal/httpapi/generated` и `internal/postgres/sqlc` вручную.

Контракт API хранится в `packages/contracts/openapi.yaml`.
