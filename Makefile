.PHONY: dev build lint check \
        web-dev web-build web-lint web-check web-knip web-preview

# Общие команды репозитория. Сейчас они делегируют задачи веб-приложению.
# Когда появится Go-бэкенд, сюда добавятся api-* цели.
dev: web-dev

build: web-build

lint: web-lint

check: web-check

web-dev:
	vp -C apps/web run dev

web-build:
	vp -C apps/web run build

web-lint:
	vp -C apps/web run lint

web-check:
	vp -C apps/web run lint && vp -C apps/web run build

web-knip:
	vp -C apps/web run knip

web-preview:
	vp -C apps/web run preview
