all: go-lint test build-all scan cypress down

.PHONY: cypress

test-results:
	mkdir -p -m 0777 test-results cypress/screenshots .go-cache

setup-directories: test-results

go-lint:
	docker compose run --rm go-lint

build:
	docker compose build --no-cache --parallel management-information

build-dev:
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build --parallel management-information yarn

build-all:
	docker compose build --parallel management-information yarn cypress

test: setup-directories
	go run gotest.tools/gotestsum@latest --format testname  --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

scan: setup-directories
	docker compose run --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-management-information:latest
	docker compose run --format sarif --output /test-results/hub.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-management-information:latest

clean:
	docker compose down
	docker compose run --rm yarn

up: clean build-dev
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up management-information yarn

down:
	docker compose down

compile-assets:
	docker compose run --rm yarn build

cypress: setup-directories clean
	docker compose run --build cypress