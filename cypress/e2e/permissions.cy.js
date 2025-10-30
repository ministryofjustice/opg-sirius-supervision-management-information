const reportingUser = "1"
const nonReportingUser = "2"

describe("Role-based permissions", () => {
   it("checks permissions for user with no reporting roles", () => {
       cy.setCookie("x-test-user-id", nonReportingUser);
       cy.visit("/downloads");
       cy.get('.govuk-heading-l').should("contain.text", "Forbidden");

       cy.visit("/uploads");
       cy.get('.govuk-heading-l').should("contain.text", "Forbidden");
   });

   it("checks permissions for reporting user role", () => {
       cy.setCookie("x-test-user-id", reportingUser);
       cy.visit("/downloads");
       cy.get('.govuk-heading-l').should("contain.text", "Management information");
       cy.visit("/uploads");
       cy.get('.govuk-heading-l').should("contain.text", "Management information");
   });
});

