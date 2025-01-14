build:
	docker build . -t marvik-api

clean:
	docker rmi -f marvik-api

run:
	docker run -d --name marvik-api-container -p 8080:8080 marvik-api

exec:
	docker exec -it marvik-api-container /bin/bash
