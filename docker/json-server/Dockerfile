## build json-server for mock API
FROM node:24.22.1-alpine3.24@sha256:1200d815de24f8dcca745fbf7651d166034a35adbcad8860a8538ffff44f1420 AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN npm ci --ignore-scripts

ENTRYPOINT ["npm", "run", "serve"]
