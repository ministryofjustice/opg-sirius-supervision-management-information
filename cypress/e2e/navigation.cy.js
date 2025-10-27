describe("Navigation page", () => {
    beforeEach(() => {
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/downloads");
    });


    describe("Tabs", () => {
        it("navigates between tabs correctly", () => {
            cy.url().should("contain", "/downloads");
            cy.contains("Management information");
            cy.contains(".moj-sub-navigation__link", "Downloads").should("be.visible");
            cy.contains(".moj-sub-navigation__link", "Uploads").should("be.visible");
            cy.contains("Uploads").click();
            cy.url().should("contain", "/uploads");
        });
    });
});