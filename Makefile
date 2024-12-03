FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
EVENT_BINARY=eventApp
NOTIFICATION_BINARY=notificationApp
LISTENER_BINARY=listenerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_event build_notification build_listener
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the authentication binary as a linux executable
build_auth:
	@echo "Building authentication binary..."
	cd ./authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_event: builds the event binary as a linux executable
build_event:
	@echo "Building event binary..."
	cd ./event-service && env GOOS=linux CGO_ENABLED=0 go build -o ${EVENT_BINARY} ./cmd/api
	@echo "Done!"

build_notification:
	@echo "Building notification binary..."
	cd ./notification-service && env GOOS=linux CGO_ENABLED=0 go build -o ${NOTIFICATION_BINARY} ./cmd/api
	@echo "Done!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd ./listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Done!"

## build_front: builds the front end binary
build_front:
	@echo "Building React front end..."
		cd ./client && npm install && npm run build
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting React front end..."
	cd ./client && npm run dev

## stop: stop the front end
stop:
	@echo "Stopping React front end..."
	@pkill -f "vite"
	@echo "Stopped React front end!"
