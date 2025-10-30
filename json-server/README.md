# json-server

json-server is the mock API we use for local development and testing. It is a simple Node Express app reading from a
JSON file, which makes it flexible for our needs.

If you need data, add it into `db.json` in the format you require it, making sure to include an `id` if you are using a
plural route (i.e. a route that could return many different entries).
json-server provides functionality for nested routes and parent/child relationships but if you require custom routing
(e.g. you always want the same data returned, regardless of the id), you can add these to `routes.json`.

If you want to inspect the data, json-server is served on port 3000 using the `docker-compose.dev.yml` so you can visit it in a browser.

For more advanced customisation, you can create your own Express middleware and include it in the `serve` script in `package.json`.

## Middleware

### User Permissions

To allow us to test different permissions validation the `switch-user`
middleware enables requests to be rerouted to the users listed in `db.json`.
We only expect users with a "Reporting User" permission to be able to view the pages.
To add custom permissions use the `x-test-user-id` cookie, the default is set to be a user with the Reporting User permissions.

e.g.:
`cy.setCookie("x-test-user-id", reportingUser);`
(in this example reportingUser is set to string "1")