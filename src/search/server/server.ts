import express from 'express';
import addHandlers, { searchMovies } from './handlers';
import cors from 'cors';
import {
  enableTorrentSearch,
  torrentIndexerSearch,
  YTSsearch,
} from './torrents';
import log from '../logger/logger';
import { initDatabase } from '../db/postgres/config';
import { download, writeFileAsync } from './utils';
const bodyParser = require('body-parser').json();

export default function startServer() {
  const port = process.env.MOVIES_API_PORT || '2222';
  enableTorrentSearch();
  initDatabase();

  const app = express();

  app.use(bodyParser);
  app.use(cors());

  // torrentIndexerSearch("I'm Not There");
  // const movies = searchMovies('Star wars', 'Movies');
  // log.info('got movies', movies);
  // YTSsearch('American Psycho');
  // download(
  //   'https://yts.mx/torrent/download/5AD19028D8880A9AA16B8B4B6E511E31F68B124E'
  // )
  //   .then((res) => {
  //     log.info('[download]', res);
  //     writeFileAsync('americanPsycho.torrent', res);
  //   })
  //   .catch((err) => log.error('[download]', err));

  addHandlers(app);

  const http = require('http').createServer(app);

  http.listen(parseInt(port), () => {
    log.info(`listening on *:${port}`);
  });
}
