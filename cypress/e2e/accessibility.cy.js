const reportingUser = "1"
const nonReportingUser = "2"

describe("Navigation page", () => {
    it("should have no accessibility violations on pages", () => {
        cy.setCookie("x-test-user-id", reportingUser);
        cy.visit("/downloads");
        cy.get('.govuk-heading-l').should("contain.text", "Management information");
        cy.checkAccessibility();
        cy.visit("/uploads");
        cy.checkAccessibility();
    });

    it("should have no accessibility violations on error page", () => {
        cy.setCookie("x-test-user-id", nonReportingUser);
        cy.visit("/downloads");
        cy.get('.govuk-heading-l').should("contain.text", "Forbidden");
        cy.checkAccessibility();
    });
});