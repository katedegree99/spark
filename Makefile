.PHONY: up down restart web schema mock

NVM_NODE := $(HOME)/.nvm/versions/node/v24.16.0/bin

up: down
	$(MAKE) -j3 web schema mock

down:
	lsof -ti:3000 | xargs kill -9 2>/dev/null || true
	lsof -ti:3001 | xargs kill -9 2>/dev/null || true
	lsof -ti:8081 | xargs kill -9 2>/dev/null || true

restart: down up

web:
	cd web && bun run dev

schema:
	cd schema && PATH="$(NVM_NODE):$$PATH" bun run preview

mock:
	cd schema && PATH="$(NVM_NODE):$$PATH" bun run mock
