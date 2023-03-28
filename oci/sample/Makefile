JAR_FILE_PATH = app/HelloWorldServer.jar
DOCKER_IMAGE_NAME = java-app-base 
DOCKER_CONTAINER_NAME = java-app


.PHONY: build
build:
	docker run --rm -v $(shell pwd)/app:/app -w /app openjdk:8-jdk-alpine sh -c 'javac HelloWorldServer.java && jar cfe HelloWorldServer.jar HelloWorldServer .'
	@echo "JAR file created successfully."

.PHONY: clean
clean:
	rm -f $(JAR_FILE_PATH)
	rm -f app/*.class

.PHONY: run
run: 
	docker run --rm -p 8080:8080 -v $(shell pwd)/$(JAR_FILE_PATH):/app.jar --name $(DOCKER_CONTAINER_NAME) $(DOCKER_IMAGE_NAME)

image:
	docker build -t $(DOCKER_IMAGE_NAME) .

