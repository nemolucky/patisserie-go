include .env

DATABASE_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

MIGRATION_PATH=./migrations


migration-up:
	migrate \
		-path $(MIGRATION_PATH) \
		-database "$(DATABASE_URL)" \
		up


migration-down:
	migrate \
		-path $(MIGRATION_PATH) \
		-database "$(DATABASE_URL)" \
		down

migration-create:
	migrate create \
		-ext sql \
		-dir $(MIGRATION_PATH) \
		-seq $(name)

reset:
	docker compose down
	sudo rm -rf /home/nemolucky/docker/postgres
	sudo rm -rf /home/nemolucky/docker/minio
	sudo mkdir -p /home/nemolucky/docker/postgres
	sudo mkdir -p /home/nemolucky/docker/minio
	docker compose up --build