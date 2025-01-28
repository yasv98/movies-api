movies-api:
	docker-compose up

movies-api-build:
	docker-compose up --build

test-db:
	cd test-db && docker-compose up

test-db-build:
	cd test-db && docker-compose up --build

integration-tests:
	cd test-db && docker-compose up -d
	go test ./tests/integration/... -v