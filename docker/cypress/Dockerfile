FROM cypress/included:16.2.1@sha256:f1f257405a482bdf95358cd30202d9c3ac69f907d5e96342ae377b8111ff7a15

WORKDIR /root

COPY cypress/package.json .
COPY cypress/package-lock.json .
RUN npm ci --ignore-scripts

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://management-information:7777/management-information

COPY cypress .
