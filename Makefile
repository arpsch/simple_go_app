# set of env variables that you need for testing
export SERVER=localhost
export PORT=8888
export SID=123
# to make it multi sector uncomment bellow and comment above line
#export SID=
export TAG=latest

SERVICE := ha_dns

server_start:
	@echo "Starting the server..."
	docker build -t ha:$(TAG) . 
	docker run -dt --name $(SERVICE) -p $(PORT):$(PORT) ha:$(TAG) -port=$(PORT) -sid=$(SID)

docker_build:
	@echo "Building the docker image..."
	docker build -t ha:$(TAG) .

docker_stop:
	@echo "Stopping the service $(SERVICE)..."
	docker stop $(SERVICE)

docker_remove:
	@echo "Stopping the service $(SERVICE) and Removing it..."
	docker stop $(SERVICE); docker rm $(SERVICE)

unit_test:
	go test -tags=unit ./...

integration_test:
	@echo "Building the docker image..."
	docker build -t ha:$(TAG) . 
	@echo "Starting the server..."
	docker run -dt --name $(SERVICE) -p $(PORT):$(PORT) ha:$(TAG) -port=$(PORT) -sid=$(SID)
	@echo "Running integration tests..."
	go test -tags=integration ./... -v
	@echo "Stopping the service $(SERVICE) and Removing it..."
	docker stop $(SERVICE); docker rm $(SERVICE)
