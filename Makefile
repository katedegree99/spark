.PHONY: up down restart web schema

up:
	$(MAKE) -j2 web schema

down:
	pkill -f "next dev" || true
	pkill -f "redocly preview" || true

restart: down up

web:
	cd web && bun run dev

schema:
	cd schema && bun run preview
