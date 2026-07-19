.PHONY: up down restart web api logs schema doc mock generate db db-down db-reset

export PATH := $(HOME)/.local/share/mise/shims:$(PATH)

up: down
	$(MAKE) -j4 api web doc mock

down:
	lsof -ti:3000 | xargs kill -9 2>/dev/null || true
	lsof -ti:3001 | xargs kill -9 2>/dev/null || true
	lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	lsof -ti:8081 | xargs kill -9 2>/dev/null || true

web:
	cd web && bun run dev

api:
	cd api && air

doc:
	cd schema && bun run preview

schema:
	cd schema && bun run preview

mock:
	cd schema && bun run mock

generate:
	cd web && bun run generate
	cd api && oapi-codegen --config=oapi-codegen.yaml ../schema/openapi/openapi.yaml

db:
	cd api && docker compose up -d db

db-down:
	cd api && docker compose down

db-reset:
	cd api && docker compose down -v && docker compose up -d db
