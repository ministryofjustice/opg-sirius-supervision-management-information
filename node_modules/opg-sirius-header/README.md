# opg-sirius-header

This is a small template repo imported into other Go repos. <br>
There is custom CSS within the sirius-header.scss that will get pulled through into these other repos. <br>

To add a new nav link, include it in the `navLinks` array and it will appear in order. Current page highlighting is done via URLs.

If the service you are adding this to wants to display the Finance tab, you can add the `show-finance="true"` attribute to the component (see `public/supervision/index.html` for example). It is the service's responsibility to ensure the user has the correct permissions to view the Finance page.

## Versioning:

Other repos are set to pull in a specific tag of this repo (within the package.json, e.g. "opg-sirius-header": "ministryofjustice/opg-sirius-header#semver:v0.9.0") <br>
When you make a change to this repo and merge it into main, you will need to update the tag used within the other repos package.json in order to pull this new version in.

### Development:

If you want to test these changes in the other repos while you are developing you can amend the other repos package.json to see changes made. <br>
(E.g. if developing against workflow amend the package.json dependency in workflow from: <br> `"opg-sirius-header": "ministryofjustice/opg-sirius-header"` <br>
to `"opg-sirius-header": "ministryofjustice/opg-sirius-header#commitId"`. <br>
Then rebuild the css/ js styling in workflow). <br> It should pull through the latest commit from sirius-header.

#### To build locally:

- Download dependencies: `yarn install` <br>
- Build Sass/ Css: `yarn compile-sass` <br>
- Build page: `yarn serve-local` <br>

This will then host with http-server, it's usually on 8080 but the console will tell you which port it's been hosted on.

Alternatively, you can use Docker: `docker-compose up opg-sirius-header`

##### Testing:

Due to the requirements for testing conditional rendering, the header needs to be served via Docker in order to use the `public/` directory structure:

- Run `docker-compose up opg-sirius-header`
- Run `yarn cypress` in another console window

(NB: Cypress expects the app to be running on 8080 which is the default port,
if this is taken and the app hosts on another port Cypress will fail)

#### To import into a new app that isn't currently using it:

Add dependency to package.json in repo you want to import it into `"opg-sirius-header": "ministryofjustice/opg-sirius-header"` <br>
Import the SCSS from sirius-header into the repo's main.scss file `@import "node_modules/opg-sirius-header/sass/sirius-header"` <br>
Import module into repo's main.js file `import "opg-sirius-header/sirius-header.js"` <br>
Run `Yarn install` and build the CSS locally

To incorporate the header into a page, use this HTML:

```
<sirius-header></sirius-header>
```

This inserts a header with the default Sirius Supervision links in the primary navigation. These links can be customised
with `<sirius-header-nav>` elements as follows:

```
<sirius-header>
    <sirius-header-nav url="/lpa">Case list</sirius-header-nav>
    <sirius-header-nav url="/another/path">Another link</sirius-header-nav>
</sirius-header>
```
