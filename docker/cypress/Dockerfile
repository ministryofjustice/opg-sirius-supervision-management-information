FROM cypress/included:15.19.0@sha256:fe59f2cac5c7451644c8780cb2fb6735fde512d06642117199a7479269702f03

WORKDIR /root

COPY cypress/package.json .
COPY cypress/yarn.lock .
RUN yarn install --frozen-lockfile --ignore-scripts

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://management-information:7777/management-information

COPY cypress .