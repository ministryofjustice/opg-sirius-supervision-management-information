const reportingUser = "1"

describe("Navigation page", () => {
    beforeEach(() => {
        cy.setCookie("x-test-user-id", reportingUser);
        cy.visit("/downloads");
    });

    describe("Tabs", () => {
        it("navigates between tabs correctly", () => {
            cy.url().should("contain", "/downloads");
            cy.contains("Management information");
            cy.contains(".moj-sub-navigation__link", "Downloads").should("be.visible");
            cy.get('[data-cy="downloads"] > .moj-sub-navigation__link').should("have.attr", "aria-current", "page")
            cy.get('[data-cy="uploads"] > .moj-sub-navigation__link').should("not.have.attr", "aria-current", "page");

            cy.contains(".moj-sub-navigation__link", "Uploads").should("be.visible");
            cy.contains("Uploads").click();
            cy.url().should("contain", "/uploads");
            cy.get('[data-cy="downloads"] > .moj-sub-navigation__link').should("not.have.attr", "aria-current", "page")
            cy.get('[data-cy="uploads"] > .moj-sub-navigation__link').should("have.attr", "aria-current", "page");
        });
    });
});