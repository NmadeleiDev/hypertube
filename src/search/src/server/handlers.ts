import { Express } from 'express';
import * as utils from './utils';
import log from '../logger/logger';
import {
  groupTorrentsByTitle,
  searchTorrents,
  torrentIndexerSearch,
  YTSsearch,
} from './torrents';
import { dbToIMovie, loadMoviesInfo, removeDuplicates } from './imdb';
import {
  getKinopoiskMovieByImdbid,
  ITranslatedMovie,
  translateMovie,
} from './kinopoisk';
import { selectMoviesFromDB } from '../db/postgres/postgres';

export const searchMovies = async (search: string, category: string) => {
  // torrentIndexerSearch(search);
  try {
    const torrents = await searchTorrents(search, category, {
      limit: 20,
      retries: 3,
    });
    const grouppedTorrents = groupTorrentsByTitle(torrents);

    log.debug('[searchMovies] grouppedTorrents', grouppedTorrents);
    log.info('loading movies');
    const movies = await loadMoviesInfo(grouppedTorrents);
    log.trace('loaded movies: ', movies);
    if (movies && movies.length) {
      const unduplicated = removeDuplicates(movies);
      log.debug('Removed duplicates: ', unduplicated);
      return unduplicated;
    }
    return [];
  } catch (e) {
    log.error(`Error getting torrents: ${e}`);
    return null;
  }
};

export default function addHandlers(app: Express) {
  app.get('/find', async (req, res) => {
    log.trace(req);
    const category = req.query['category'].toString();
    const search = req.query['search'].toString();
    log.info(`[GET /find] category: ${category}, serach: ${search}`);

    try {
      let movies: ITranslatedMovie[] = null;
      let dbMovies = await selectMoviesFromDB(search);
      if (dbMovies) {
        const ens = dbMovies.map((movie) => dbToIMovie(movie));
        const promises = ens.map((en) => translateMovie(en));
        movies = await Promise.all(promises);
      }
      log.info('[GET /find] movies from database', movies);
      if (!movies || !movies.length) {
        movies = await YTSsearch(search);
        log.info('[GET /find] YTSsearch movies', movies);
      }
      if (!movies || !movies.length) {
        movies = await searchMovies(search, category);
        log.info('[GET /find] RARBG and TPB movies', movies);
      }
      if (movies && movies.length) {
        log.info(
          `[GET /find] found movies for query: ${search}, results: `,
          movies
        );
        res.json(utils.createSuccessResponse(movies)).status(200);
      } else {
        log.info('[GET /find] No movies found');
        res
          .json(utils.createErrorResponse('Could not find movies'))
          .status(404);
      }
    } catch (e) {
      res.json(utils.createErrorResponse('Error getting torrents')).status(500);
    }
  });
  app.get('/translate', async (req, res) => {
    log.trace(req);
    const imdbid = req.query.imdbid.toString();
    const title = req.query.title.toString();

    try {
      const ru = await getKinopoiskMovieByImdbid(title, imdbid);
      log.info('[GET /translate] ru', ru);
      if (!ru)
        res
          .json(utils.createErrorResponse('Could not find movies'))
          .status(404);
      else res.json(utils.createSuccessResponse(ru)).status(200);
    } catch (e) {
      log.error(e);
      res
        .json(utils.createErrorResponse('Error translating movie'))
        .status(500);
    }
  });
}
