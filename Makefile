include app.env
export

up:
	docker compose up -d

up-prod:
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up

down:
	docker compose down

server:
	go run main.go

postgres:
	 docker run --name $(PG_CONTAINER_NAME) -p $(PG_PORT_MAPPING) -e POSTGRES_PASSWORD=$(PGPASSWORD) -e POSTGRES_USER=$(PGUSER) -d $(IMAGE)
#	 docker start $(DB_CONTAINER_NAME)

createdb:
	 docker exec -it $(PG_CONTAINER_NAME) createdb --username=$(PGUSER) --owner=$(PGUSER) $(PGDATABASE)

dropdb:
	 docker exec -it $(PG_CONTAINER_NAME) dropdb $(PGDATABASE) --username=$(PGUSER)

migrate_up:
	migrate -path ./db/migrations -database "$(DATABASE_URL)" -verbose up

migrate_down:
	migrate -path ./db/migrations -database "$(DATABASE_URL)" -verbose down

restart_db:
	$(MAKE) migrate_down
	$(MAKE) migrate_up

shuv:
	git add .
	git commit -a
	git push


.PHONY: down, up, up-prod, server, createdb, dropdb, migrate_up, migrate_down, restart_db
