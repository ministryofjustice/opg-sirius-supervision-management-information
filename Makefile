all: go-lint test build-all scan cypress management-information down

.PHONY: cypress

test-results:
	mkdir -p -m 0777 test-results cypress/screenshots .go-cache

setup-directories: test-results

go-lint:
	docker compose run --rm go-lint

gosec: setup-directories
	docker compose run --rm gosec

build:
	docker compose build --no-cache --parallel management-information management-information-api

build-dev:
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build --no-cache --parallel management-information management-information-api yarn json-server

build-all:
	docker compose build --parallel management-information management-information-api json-server cypress

hub-tests: setup-directories
	docker compose run --rm hub-test-runner

api-tests: setup-directories
	go run gotest.tools/gotestsum@latest --format testname  --junitfile test-results/api-unit-tests.xml -- -p 1 ./finance-api/... -coverprofile=test-results/api-coverage.txt

combine-coverage:
	cat test-results/hub-coverage.txt > test-results/coverage.txt
	tail -n +2 test-results/api-coverage.txt >> test-results/coverage.txt

test: hub-tests api-tests combine-coverage

clean:
	docker compose down
	docker compose run --rm yarn

dev-up: clean build-dev
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up management-information management-information-api localstack yarn

up: clean compile-assets build-all
	docker compose -f docker-compose.yml up -d --wait management-information

down:
	docker compose down

compile-assets:
	docker compose run --rm yarn build

cypress: setup-directories clean
	docker compose run --build cypress

cypress-single: setup-directories clean
	docker compose run --rm cypress run --spec e2e/$(SPEC)
