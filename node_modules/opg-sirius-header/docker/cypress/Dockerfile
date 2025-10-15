FROM cypress/included:12.15.0

WORKDIR /root

RUN yarn add "cypress-failed-log@2.9.2" "axe-core@^4.7.2" "cypress-axe@^1.4.0"

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://opg-sirius-header:8080

COPY cypress.config.js .
COPY cypress cypress
COPY package.json package.json
ENTRYPOINT ["yarn", "cypress-headless"]