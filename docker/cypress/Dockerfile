FROM cypress/included:15.20.1@sha256:327d62ee6d1c9bd7b4b490bb2772c1b70d37d56647953a4ec3105a939b7c61ff

WORKDIR /root

COPY cypress/package.json .
COPY cypress/yarn.lock .
RUN yarn install --frozen-lockfile --ignore-scripts

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://management-information:7777/management-information

COPY cypress .