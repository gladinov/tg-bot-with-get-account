# =========================
# Variables
# =========================

COMPOSE_FILE = deployments/docker-compose.yml
ENV_FILE = deployments/env/prod.env


# =========================
# Phony targets
# =========================

.PHONY: up build-up down 
        
# =========================
# Docker: production
# =========================

up:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d

build-up:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d --build

down:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) down -v

