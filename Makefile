COMPOSE = docker compose -p renal-relay-demo --env-file .env -f deploy/compose.yaml -f deploy/compose.demo.yaml

.PHONY: up build down logs status seed

up:
	$(COMPOSE) up -d --build web

build:
	$(COMPOSE) build web api

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

status:
	$(COMPOSE) ps -a

seed:
	$(COMPOSE) --profile dev-seed run --rm seed
