const { defineConfig } = require("cypress")

module.exports = defineConfig({
    fixturesFolder: false,
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
                failed: require("cypress-failed-log/src/failed")()
            });
        },
        baseUrl: "http://localhost:7777/management-information",
        modifyObstructiveCode: false,
    },
    viewportWidth: 1000,
    viewportHeight: 1000,
});