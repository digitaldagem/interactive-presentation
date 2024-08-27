SRC_DOCKER_IMAGES := $(shell docker images -q interactive-presentation-src)
POSTGRES_DOCKER_IMAGES := $(shell docker images -q postgres)

up:
	docker-compose up -d --build --remove-orphans --timeout 60

up-local:
	docker-compose up --build --remove-orphans

down:
	docker-compose down -v --remove-orphans

	if [ -n "$(SRC_DOCKER_IMAGES)" ]; then docker rmi $(SRC_DOCKER_IMAGES); fi
	if [ -n "$(POSTGRES_DOCKER_IMAGES)" ]; then docker rmi $(POSTGRES_DOCKER_IMAGES); fi

integration-tests:
	go test ./integration_tests -v

integration-tests-local:
	go test ./integration_tests/backend_test.go -v

unit-tests:
	go test ./src/utilities -v

.PHONY: up up-local down integration-tests integration-tests-local unit-tests