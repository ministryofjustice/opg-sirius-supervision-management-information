all: go-lint test build-all scan cypress management-information down

.PHONY: cypress

test-results:
	mkdir -p -m 0777 test-results cypress/screenshots .trivy-cache .go-cache

setup-directories: test-results

go-lint:
	docker compose run --rm go-lint

gosec: setup-directories
	docker compose run --rm gosec

build:
	docker compose build --no-cache --parallel management-information

build-dev:
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build --parallel management-information yarn

build-all:
	docker compose build --parallel management-information yarn cypress

test: setup-directories
	go run gotest.tools/gotestsum@latest --format testname  --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

scan: setup-directories
	docker compose run --rm trivy image --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-management-information:latest
	docker compose run --rm trivy image --format sarif --output /test-results/management-info.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-management-information:latest

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