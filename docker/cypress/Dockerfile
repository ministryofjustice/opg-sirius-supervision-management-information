FROM cypress/included:16.3.0@sha256:3efa8d156a094a6c13ddf3fba616e97f7492706685064b75a39823e1147eca97

WORKDIR /root

COPY cypress/package.json .
COPY cypress/package-lock.json .
RUN npm ci --ignore-scripts

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://management-information:7777/management-information

COPY cypress .
