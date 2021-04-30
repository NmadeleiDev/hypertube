import express from 'express';
import {
  addCommentsHandlers,
  addMoviesHandlers,
  addRatingHandlers,
} from './handlers';
import cors from 'cors';
import log from '../logger/logger';
const bodyParser = require('body-parser').json();

export default function startServer() {
  const port = "2222";

  const app = express();

  app.use(bodyParser);
  app.use(cors());

  addMoviesHandlers(app);
  addCommentsHandlers(app);
  addRatingHandlers(app);

  const http = require('http').createServer(app);

  http.listen(parseInt(port), () => {
    log.info(`listening on *:${port}`);
  });
}
