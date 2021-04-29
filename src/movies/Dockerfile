FROM node:15.7.0-alpine3.10

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY package.json tsconfig.json /usr/src/app/
COPY src/ /usr/src/app/src/
RUN npm install
CMD NODE_ENV=production npm start
