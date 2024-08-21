# Makefile for managing the Docker application

# Default target
.PHONY: all
all: build up

# Build the Docker images
.PHONY: build
build:
	docker-compose build

# Start the services in detached mode
.PHONY: up
up:
	docker-compose up -d

# Stop and remove the containers
.PHONY: down
down:
	docker-compose down

# View logs for the app service
.PHONY: logs
logs:
	docker-compose logs -f app

# Run tests inside the container
.PHONY: test
test:
	docker-compose run app go test ./...

# Clean up Docker volumes and networks
.PHONY: clean
clean:
	docker-compose down -v --rmi all
