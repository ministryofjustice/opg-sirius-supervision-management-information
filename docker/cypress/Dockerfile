FROM cypress/included:16.2.0@sha256:8286fe42a37ca6766d2987732e89f046997a59e5387baffd940a1e99501411f7

WORKDIR /root

COPY cypress/package.json .
COPY cypress/package-lock.json .
RUN npm ci --ignore-scripts

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://management-information:7777/management-information

COPY cypress .