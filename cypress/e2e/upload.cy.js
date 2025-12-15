const reportingUser = "1"

describe("Upload page", () => {
    beforeEach(() => {
        cy.setCookie("x-test-user-id", reportingUser);
        cy.visit("uploads");
    });

    describe("Validation", () => {
        it("displays an error message when no upload type is selected", () => {
            cy.contains('.govuk-button', 'Upload file').click();
            cy.get('.govuk-error-summary').contains('Please select a report to upload');
            cy.get('#f-UploadType > .govuk-label').contains('Please select a report to upload');
            cy.get('#f-UploadType').should('have.class', 'govuk-form-group--error');
        });
        it("displays an error message when upload type is bonds and no bond provider is selected", () => {
            cy.get('#upload-type').select('Bonds');
            cy.contains('.govuk-button', 'Upload file').click();
            cy.get('.govuk-error-summary').contains('Please select a bond provider');
            cy.get('#f-BondProvider > .govuk-label').contains('Please select a bond provider');
            cy.get('#f-BondProvider').should('have.class', 'govuk-form-group--error');
        });
    });

    describe("Uploading a file", () => {
        it("displays success message on successful upload", () => {
            cy.get('#upload-type').select('Bonds');
            cy.get('#bond-provider').select('Marsh');
            cy.get('input[type="file"]').selectFile('cypress/fixtures/bonds-without-orders.csv');
            cy.contains('.govuk-button', 'Upload file').click();
            cy.url().should('include', '/uploads?success=upload');
            cy.get('.moj-banner').contains('File successfully uploaded');
        });
    });
});