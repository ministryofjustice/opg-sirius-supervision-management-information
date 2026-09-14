## build json-server for mock API
FROM node:24-alpine3.23@sha256:0388af2af070cd4736a1567cfed02469ba117848845b4165d87a333edb53d2ca AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN yarn install --frozen-lockfile --ignore-scripts

ENTRYPOINT ["yarn", "serve"]
