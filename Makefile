.PHONY: install dev dev-down dev-reset dev-logs dev-ps build lint test check generate generate-check auth-test admin-create \
        api-dev api-generate api-build api-test api-lint api-check api-migrate \
        web-install web-dev web-build web-lint web-check web-knip web-preview web-api-generate

COMPOSE = docker compose

install: web-install

dev:
	$(COMPOSE) up --build --watch

dev-down:
	$(COMPOSE) down --remove-orphans

dev-reset:
	$(COMPOSE) down --volumes --remove-orphans
	$(COMPOSE) up --build --watch

dev-logs:
	$(COMPOSE) logs -f

dev-ps:
	$(COMPOSE) ps

build: api-build web-build

lint: api-lint web-lint

test: api-test

check: api-check web-check

generate: api-generate web-api-generate

generate-check: generate
	@test -z "$$(git status --porcelain --untracked-files=all -- \
		apps/api/internal/httpapi/generated \
		apps/api/internal/postgres/sqlc \
		apps/web/src/shared/api)" || { \
		git status --short --untracked-files=all -- \
			apps/api/internal/httpapi/generated \
			apps/api/internal/postgres/sqlc \
			apps/web/src/shared/api; \
		exit 1; \
	}

auth-test:
	@set -eu; \
	trap '$(COMPOSE) -f compose.test.yaml down --volumes --remove-orphans' EXIT; \
	$(COMPOSE) -f compose.test.yaml down --volumes --remove-orphans; \
	$(COMPOSE) -f compose.test.yaml up --build --abort-on-container-exit --exit-code-from auth-test

admin-create:
	$(COMPOSE) --profile tools run --rm admin create

api-dev:
	$(MAKE) -C apps/api dev

api-generate:
	$(MAKE) -C apps/api generate

api-build:
	$(MAKE) -C apps/api build

api-test:
	$(MAKE) -C apps/api test

api-lint:
	$(MAKE) -C apps/api lint

api-check:
	$(MAKE) -C apps/api check

api-migrate:
	$(MAKE) -C apps/api migrate

web-install:
	vp -C apps/web install

web-dev:
	vp -C apps/web run dev

web-build:
	vp -C apps/web run build

web-api-generate:
	vp -C apps/web run api:generate

web-lint:
	vp -C apps/web run lint

web-check:
	vp -C apps/web run check && vp -C apps/web run build

web-knip:
	vp -C apps/web run knip

web-preview:
	vp -C apps/web run preview
