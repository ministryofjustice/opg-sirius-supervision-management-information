FROM cypress/included:15.21.0@sha256:af15ca5f3b851e32360f37e1d446d6e6821fd11928c87a2fd7fad12e8495e63b

WORKDIR /root

COPY cypress/package.json .
COPY cypress/yarn.lock .
RUN yarn install --frozen-lockfile --ignore-scripts

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://management-information:7777/management-information

COPY cypress .