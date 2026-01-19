# opg-sirius-supervision-management-information

### Major dependencies

- [Go](https://golang.org/) (>= 1.22)
- [docker compose](https://docs.docker.com/compose/install/) (>= 2.26.0)

#### Installing dependencies locally:
(This is only necessary if running without docker, you will need to be logged into sirius with a user who has `reporting user` permissions)

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

Running without docker with: <br>
`yarn install && yarn build` <br>
`go run main.go`

This will host at:
`localhost:7777/downloads`.

-----
## Run the unit tests
`make test`

## Run *one* Cypress test headless (i.e. not in UI)
`make cypress-single SPEC=upload.cy.js`

## Run *all* the Cypress tests headless
`make build-all` (optional) <br>
`make cypress`

## Run the Cypress tests in UI
`make up` in one terminal (wait for the app to build) <br>
`cd cypress`
`npx cypress open baseUrl=http://localhost:7777/downloads` in another terminal