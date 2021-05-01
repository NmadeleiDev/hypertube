FROM node:15.7.0-alpine3.10

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY package.json tsconfig.json ./
RUN npm install
COPY src ./src/
CMD NODE_ENV=production npm start
