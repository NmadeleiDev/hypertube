FROM node:15.7.0-alpine3.10

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY package.json /usr/src/app
# RUN npm install
COPY . /usr/src/app

CMD NODE_ENV=production npm start
