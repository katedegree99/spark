.PHONY: up down restart web schema mock

up:
	$(MAKE) -j3 web schema mock

down:
	pkill -f "next dev" || true
	pkill -f "redocly preview" || true
	pkill -f "redocly mock" || true

restart: down up

web:
	cd web && bun run dev

schema:
	cd schema && bun run preview

mock:
	cd schema && bun run mock
