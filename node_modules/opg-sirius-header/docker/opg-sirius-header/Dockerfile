FROM node:19.3.0-alpine3.16

WORKDIR /app

COPY sass sass
COPY package.json package.json
COPY public public
# COPY sirius-header.js public/js/sirius-header.js

RUN yarn install
RUN yarn compile-sass
ENTRYPOINT ["yarn", "serve"]