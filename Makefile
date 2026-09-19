.PHONY: install dev build lint check \
        web-install web-dev web-build web-lint web-check web-knip web-preview web-api-generate

# Общие команды репозитория. Сейчас они делегируют задачи веб-приложению.
# Когда появится Go-бэкенд, сюда добавятся api-* цели.
install: web-install

dev: web-dev

build: web-build

lint: web-lint

check: web-check

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
