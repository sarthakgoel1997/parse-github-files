APP_NAME=parse-github-files
DOCKER_IMAGE=$(APP_NAME)
CONTAINER_NAME=$(APP_NAME)
EXTERNAL_PORT=9000
INTERNAL_PORT=8000
PERSONAL_ACCESS_TOKEN=<YOUR_PERSONAL_ACCESS_TOKEN>

.PHONY: build run dev stop logs clean

build:
	docker build -t $(DOCKER_IMAGE) .

run:
	docker run --name $(CONTAINER_NAME) -e PORT=$(INTERNAL_PORT) -e PERSONAL_ACCESS_TOKEN=$(PERSONAL_ACCESS_TOKEN) -p $(EXTERNAL_PORT):$(INTERNAL_PORT) -d $(DOCKER_IMAGE)

stop:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

dev: build stop run

logs:
	docker logs -f $(CONTAINER_NAME)

clean: stop
	docker rmi -f $(DOCKER_IMAGE)

query-db:
	docker exec -it $(CONTAINER_NAME) sh -c "cd /app && sqlite3 database.db"

test:
	@echo "Running unit tests with coverage..."
	PERSONAL_ACCESS_TOKEN=$(PERSONAL_ACCESS_TOKEN) go test ./... -coverprofile=coverage.out -v
	go tool cover -func=coverage.out

coverage:
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html
