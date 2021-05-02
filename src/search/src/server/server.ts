import express from 'express';
import addHandlers from './handlers';
import cors from 'cors';
import { enableTorrentSearch } from './torrents';
import log from '../logger/logger';
import { initDatabase } from '../db/postgres/config';
const bodyParser = require('body-parser').json();

export default function startServer() {
  const port = '2222';
  enableTorrentSearch();
  initDatabase();

  const app = express();

  app.use(bodyParser);
  app.use(cors());

  addHandlers(app);

  const http = require('http').createServer(app);

  http.listen(parseInt(port), () => {
    log.info(`listening on *:${port}`);
  });
}
