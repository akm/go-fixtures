.PHONY: up
up:
	docker compose up -d
	docker compose exec mysql mysqladmin --wait --count 5 ping

.PHONY: down
down:
	docker compose down
