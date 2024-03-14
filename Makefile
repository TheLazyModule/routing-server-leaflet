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

server-container:
	docker run --name routing-app -p 8080:8080 --network routes-network -e GIN_MODE=release -e  routing-app:latest

postgres:
	 docker run --name $(PG_CONTAINER_NAME) -p $(PG_PORT_MAPPING) -e POSTGRES_PASSWORD=$(PGPASSWORD) -e POSTGRES_USER=$(PGUSER) -d $(IMAGE)

createdb:
	 docker exec -it $(PG_CONTAINER_NAME) createdb --username=$(PGUSER) --owner=$(PGUSER) $(PGDATABASE)

dropdb:
	 docker exec -it $(PG_CONTAINER_NAME) dropdb $(PGDATABASE) -f --username=$(PGUSER)

migrate:
	cd ./geojson-graph/ && python3 main.py && cd ..

migrate_up:
	migrate -path ./db/migrations -database "$(DATABASE_URL)" -verbose up

migrate_down:
	migrate -path ./db/migrations -database "$(DATABASE_URL)" -verbose down

restart_db:
	$(MAKE) migrate_down
	$(MAKE) migrate_up

docker_build:
	docker build -t routing-app:latest .

shuv:
	git add .
	git commit -a
	git push

generate:
	sqlc generate


.PHONY: down, up, up-prod, server, createdb, dropdb, migrate_up, migrate_down, restart_db
