## build json-server for mock API
FROM node:24-alpine3.23@sha256:838d276a37a1e42da2a5d01e0ecbb4401ca4147728f62ca2eba6d023a352f88f AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN npm ci --ignore-scripts

ENTRYPOINT ["npm", "run", "serve"]
