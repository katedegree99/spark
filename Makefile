.PHONY: up down restart web api logs schema doc mock generate db db-down db-reset

NVM_NODE := $(HOME)/.nvm/versions/node/v24.16.0/bin

up: down
	$(MAKE) -j4 api web doc mock

down:
	lsof -ti:3000 | xargs kill -9 2>/dev/null || true
	lsof -ti:3001 | xargs kill -9 2>/dev/null || true
	lsof -ti:8081 | xargs kill -9 2>/dev/null || true
	cd api && docker compose down 2>/dev/null || true

web:
	cd web && bun run dev

api:
	cd api && docker compose up --build -d app

logs:
	cd api && docker compose logs -f app

doc:
	cd schema && PATH="$(NVM_NODE):$$PATH" bun run preview

schema:
	cd schema && PATH="$(NVM_NODE):$$PATH" bun run preview

mock:
	cd schema && PATH="$(NVM_NODE):$$PATH" bun run mock

generate:
	cd web && bun run generate
	cd api && PATH="$$PATH:$$(go env GOPATH)/bin" oapi-codegen --config=oapi-codegen.yaml ../schema/openapi/openapi.yaml

db:
	cd api && docker compose up -d db

db-down:
	cd api && docker compose down

db-reset:
	cd api && docker compose down -v && docker compose up -d db
