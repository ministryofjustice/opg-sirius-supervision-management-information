FROM cypress/included:14.5.4

WORKDIR /root

COPY cypress/package.json .
COPY cypress/yarn.lock .
RUN yarn install

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://management-information:7777/management-information

COPY cypress .