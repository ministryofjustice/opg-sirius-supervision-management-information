import cypress_failed_log from "cypress-failed-log/src/failed";
import {defineConfig} from "cypress";

export default defineConfig({
  e2e: {
    setupNodeEvents(on, config) {
      on("task", {
        log(message) {
          console.log(message);
          return null
        },
        table(message) {
          console.table(message);
          return null
        },
        failed: cypress_failed_log()
      });
    },
    baseUrl: 'http://localhost:7777/management-information',
    specPattern: "e2e/**/*.cy.{js,ts}",
    screenshotsFolder: "screenshots",
    supportFile: "support/e2e.ts",
    modifyObstructiveCode: false,
  },
  viewportWidth: 1000,
  viewportHeight: 1000,
});