.PHONY: up down restart web schema mock

NVM_NODE := $(HOME)/.nvm/versions/node/v24.16.0/bin

up:
	$(MAKE) -j3 web schema mock

down:
	pkill -f "next dev" || true
	pkill -f "redocly preview-docs" || true
	pkill -f "prism mock" || true

restart: down up

web:
	cd web && bun run dev

schema:
	cd schema && PATH="$(NVM_NODE):$$PATH" bun run preview

mock:
	cd schema && PATH="$(NVM_NODE):$$PATH" bun run mock
