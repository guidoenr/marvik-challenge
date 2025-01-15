# some vars
CONTAINER_NAME=marvik-api-container
IMAGE_NAME=marvik-api
PORT=8080
API_URL=http://localhost:$(PORT)


# build the Docker image
build:
	docker build . -t $(IMAGE_NAME)


# run the container
run:
	docker run --rm -it --name $(CONTAINER_NAME) -p $(PORT):$(PORT) $(IMAGE_NAME)


# clean up trash
clean:
	docker kill $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)
	docker rmi -f $(IMAGE_NAME)


# enter to the container to drop commands
exec:
	docker exec -it $(CONTAINER_NAME) /bin/bash


# check if the server is up, exit with error if not
check-server-up:
	@echo "checking if the server is up..."
	@curl -s $(API_URL) | grep '"status":"ready"' > /dev/null || (echo "server is not up, skipping tests!" && exit 1)


# test (server must be up to run the tests)
test: check-server-up
	go test -v


