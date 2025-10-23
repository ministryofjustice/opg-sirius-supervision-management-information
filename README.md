# opg-sirius-supervision-management-information

### Major dependencies

- [Go](https://golang.org/) (>= 1.22)
- [docker compose](https://docs.docker.com/compose/install/) (>= 2.26.0)

#### Installing dependencies locally:
(This is only necessary if running without docker)

- `yarn install`
- `go mod download`
---

## Local development

The application ran through Docker can be accessed on `localhost:7777/management-information/downloads`.

To enable debugging and hot-reloading of Go files:

`make up`

Hot-reloading is managed independently and should happen seamlessly. Hot-reloading for web assets (JS, CSS, etc.)
is also provided via a Yarn watch command.

-----
## Run the unit/integration tests

`make test`

## Run the Cypress tests

`make cypress`

## Run the Cypress tests in UI
`make up` in one terminal (wait for the app to build)
`yarn && yarn cypress` in another terminal