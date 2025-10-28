const { defineConfig } = require('cypress')

module.exports = defineConfig({
  e2e: {
    // We've imported your old cypress plugins here.
    // You may want to clean this up later by importing these.
    setupNodeEvents(on, config) {
      return require('./cypress/plugins/index.js')(on, config)
    },
    baseUrl: 'http://localhost:7777/management-information',
    specPattern: "cypress/e2e/**/*.cy.{js,ts}",
    screenshotsFolder: "cypress/screenshots",
    supportFile: "cypress/support/e2e.ts",
    modifyObstructiveCode: false,
  },
  viewportWidth: 1000,
  viewportHeight: 1000,
})