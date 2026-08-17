FROM cypress/included:15.20.0@sha256:5baa6e2f1831e7b9823246d22b5914fe77187a9bab030aa2f44fceabde17c9d1

WORKDIR /root

COPY cypress/package.json .
COPY cypress/yarn.lock .
RUN yarn install --frozen-lockfile --ignore-scripts

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://management-information:7777/management-information

COPY cypress .