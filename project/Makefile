SHELL=cmd.exe
FRONT_END_BINARY=frontApp.exe
BACKEND_BINARY=backendApp


up:
	@echo "Starting docker images..."
	docker-compose up -d
	@echo "Docker started!"

## docker_up_build: Build all projects and start docker compose
up_build: build_backend
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

build_backend:
	@echo Building backend end binary...
	cd ..\backend && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${BACKEND_BINARY} ./cmd/api
	@echo Done!



build_front:
	@echo "Building front end binary..."
	cd ..\front-end && set CGO_ENABLED=0 && go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"


start: build_front
	@echo "Starting front end"
	cd ../front-end && go build -o ${FRONT_END_BINARY} ./cmd/web
	cd ../front-end && ${FRONT_END_BINARY}

## docker_down: Stop docker compose
down:
	@echo "Stopping docker images..."
	docker-compose down
	@echo "Docker stopped!"
