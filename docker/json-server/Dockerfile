## build json-server for mock API
FROM node:24.22.1-alpine3.24@sha256:398bcf38afb035ae76feb9ad3c7493a93feff4751351c65310811761e37bbfc7 AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN npm ci --ignore-scripts

ENTRYPOINT ["npm", "run", "serve"]
