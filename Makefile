up:
	docker compose up -d --build

down:
	docker compose down

restart: down up

rebuild: clean up

clean:
	docker compose down -v

logs:
	docker compose logs -f app

test: mocks
	go test -v ./...

mocks:
	go generate ./...

fmt:
	go fmt ./...

migrate-up:
	docker exec -i warehouse-control-postgres psql -U warehouse_control_user -d warehouse_control_db < migrations/init/001_init.sql

migrate-down:
	docker exec -i warehouse-control-postgres psql -U warehouse_control_user -d warehouse_control_db < migrations/manual/002_cleanup.sql

migrate-reset: migrate-down migrate-up