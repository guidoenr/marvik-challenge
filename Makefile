# some vars
CONTAINER_NAME=marvik-api-container
IMAGE_NAME=marvik-api
PORT=8080

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


