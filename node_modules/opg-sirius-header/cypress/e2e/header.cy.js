import "cypress-axe";

describe("header spec", () => {
  it("has all of the expected information within the header", () => {
    cy.visit("index.html");
    cy.contains(".govuk-header__link--homepage", "OPG");
    cy.contains(".govuk-header__service-name", "Sirius");
  });

  const expectedTitle = [
    "Power of Attorney",
    "Supervision",
    "Admin",
    "Sign out",
  ];
  const expectedUrl = ["/lpa", "/supervision", "/admin", "/logout"];

  it("has working nav links within header banner", () => {
    cy.visit("index.html");
    cy.get("#header-navigation")
      .children()
      .each(($el, index) => {
        cy.wrap($el).should("contain", expectedTitle[index]);
        const $expectedLinkName = expectedUrl[index];
        cy.wrap($el)
          .find("a")
          .should("have.attr", "href")
          .and("contain", `${$expectedLinkName}`);
      });
  });

  it("has all the expected links within the secondary nav list", () => {
    cy.visit("index.html");
    const expectedTitle = ["Create client", "Workflow", "Guidance", "Finance"];
    const expectedUrl = [
      "/supervision/#/clients/search-for-client",
      "/supervision/workflow",
      "https://wordpress.sirius.opg.service.justice.gov.uk",
      "/supervision/#/finance-hub/reporting",
    ];
    cy.get("#header-navigation")
      .get(".moj-primary-navigation__list")
      .children()
      .each(($el, index) => {
        cy.wrap($el).should("contain", expectedTitle[index]);
        const $expectedLinkName = expectedUrl[index];
        cy.wrap($el)
          .find("a")
          .should("have.attr", "href")
          .and("contain", `${$expectedLinkName}`);
      });
  });

  it("has the feedback link", () => {
    cy.visit("index.html");
    cy.get("#feedback-span").find("a").should("contain", "feedback");
  });

  it("highlights the current page nav", () => {
    cy.visit("/supervision/workflow");
    cy.get(".moj-primary-navigation__list > .moj-primary-navigation__item")
      .contains("Workflow")
      .should("have.attr", "aria-current", "page");

    cy.visit("/supervision");
    cy.get(".moj-primary-navigation__list > .moj-primary-navigation__item")
      .contains("Workflow")
      .should("not.have.attr", "aria-current");
  });

  it("hides the finance tab when not set to show", () => {
    // no show-finance passed in (default)
    cy.visit("/supervision/workflow");
    cy.get(".moj-primary-navigation__list")
      .contains("Finance")
      .should("be.hidden");

    // show-finance is "false"
    cy.visit("/supervision");
    cy.get(".moj-primary-navigation__list")
      .contains("Finance")
      .should("be.hidden");
  });

  it("meets accessibility standards", () => {
    cy.visit("/index.html");
    cy.injectAxe();
    cy.checkA11y(null, {
      rules: {
        "landmark-one-main": { enabled: false },
        "page-has-heading-one": { enabled: false },
      },
    });
  });

  it("can be customised with sirius-header-nav elements", () => {
    cy.visit("/lpa");
    let linkTextContents = [];

    const navList = cy
      .get(".moj-primary-navigation__link")
      .each((item) => {
        linkTextContents.push(item.text().trim());
      })
      .then(() => {
        cy.wrap(linkTextContents).should("eql", [
          "Case list",
          "Another nav item",
        ]);
      });
  });

  it("meets accessibility standards", () => {
    cy.visit("/index.html");
    cy.injectAxe();
    cy.checkA11y(null, {
      rules: {
        "landmark-one-main": { enabled: false },
        "page-has-heading-one": { enabled: false },
      },
    });

    cy.visit("/lpa");
    cy.injectAxe();
    cy.checkA11y(null, {
      rules: {
        "landmark-one-main": { enabled: false },
        "page-has-heading-one": { enabled: false },
      },
    });
  });
});
