import { Express } from 'express';
import { createSuccessResponse, createErrorResponse } from './utils';
import log from '../logger/logger';
import { IComment, IMovie, ITranslatedMovie } from '../model/model';
import {
  deleteComment,
  insertComment,
  selectCommentById,
  selectCommentsByMovieID,
  updateComment,
} from '../db/postgres/comments';
import {
  dbToIMovie,
  getKinopoiskMovieFromDB,
  getMovieInfo,
  getMovies,
  IRating,
  KinoDBToIMovie,
  updateMovieRating,
} from '../db/postgres/movies';

export function addMoviesHandlers(app: Express) {
  app.get('/movies', async (req, res) => {
    log.trace(req);
    const limit = +req.query.limit || 5;
    const offset = +req.query.offset || 0;
    const id = (req.query.id as string) || null;

    try {
      if (id) {
        const movie = await getMovieInfo(id);
        log.debug('[GET /movies] movieById:', movie);
        if (movie) res.json(createSuccessResponse([movie])).status(200);
        else res.json(createErrorResponse(null)).status(404);
        return;
      }

      const ens: IMovie[] = (await getMovies(limit, offset)).map((movie) =>
        dbToIMovie(movie)
      );
      log.info('[addMoviesHandlers] got en movies', ens);
      const movies: ITranslatedMovie[] = await Promise.all(
        ens.map((en) => getMovieInfo(en.id))
      );
      log.info('[addMoviesHandlers] got translated movies', movies);

      if (movies) res.json(createSuccessResponse(movies)).status(200);
      else res.json(createErrorResponse(null)).status(404);
    } catch (e) {
      log.error(`Error getting movies: ${e}`);
      res.json(createErrorResponse('Error getting movies')).status(500);
    }
  });
  app.get('/genres', async (req, res) => {
    log.debug(req);
    res.json(createErrorResponse(null)).status(404);
  });
  app.get('/byname', async (req, res) => {
    const limit = +req.query.limit || 5;
    const offset = +req.query.offset || 0;
    const letter = (req.query.letter as string) || undefined;
    log.debug('letter:', req.params.letter);
    try {
      const response = await getMovies(limit, offset, letter);
      log.debug('[/byname] got getMovies response', response);
      const ens: IMovie[] = response.map((movie) => dbToIMovie(movie));
      log.info(`[/byname] en movies for ${letter}`, ens);

      const movies: ITranslatedMovie[] = await Promise.all(
        ens.map((en) => getMovieInfo(en.id))
      );
      log.info('[/byname] got translated movies', movies);

      if (movies) res.json(createSuccessResponse(movies)).status(200);
      else res.json(createErrorResponse(null)).status(404);
    } catch (e) {
      log.error(`Error getting movies: ${e}`);
      res.json(createErrorResponse('Error getting movies')).status(500);
    }
  });
}

export function addRatingHandlers(app: Express) {
  app.patch('/rating', async (req, res) => {
    const headers = req.headers;
    const cookies = req.cookies;
    log.debug(req.body, headers, cookies);

    try {
      const rating = (+req.body?.rating as IRating) || null;
      const movieid = (req.body?.movieId as string) || null;
      const token = (req.headers.accesstoken as string) || '';
      log.debug(`rating: ${rating}, movieid: ${movieid}, token: ${token}`);

      if (rating === null || !movieid || !token)
        return res
          .json(createErrorResponse('accessToken, ID and RATING are required'))
          .status(400);

      const newRating = await updateMovieRating(movieid, rating, token);
      res.json(createSuccessResponse(newRating)).status(200);
    } catch (e) {
      log.error(`Error getting movies: ${e}`);
      res.json(createErrorResponse('Error getting movies')).status(500);
    }
  });
}

export function addCommentsHandlers(app: Express) {
  app.get('/comments', async (req, res) => {
    log.debug(req.query);
    const limit = +req.query.limit || 5;
    const offset = +req.query.offset || 0;
    const movieId = req.query.movieId as string;
    try {
      const comments = await selectCommentsByMovieID(movieId, limit, offset);
      res.json(createSuccessResponse(comments)).status(200);
    } catch (e) {
      log.error(`Error getting comments: ${e}`);
      res.json(createErrorResponse('Error getting comments')).status(500);
    }
  });
  app.post('/comment', async (req, res) => {
    log.trace(req);
    try {
      const comment = req.body as IComment;
      const newComment = (await insertComment(comment)).rows[0] as IComment;
      if (!newComment)
        return res
          .json(createErrorResponse('Error posting comment'))
          .status(500);
      const result = await selectCommentById(newComment.id);
      log.debug(result);
      res.json(createSuccessResponse(result)).status(200);
    } catch (e) {
      log.error(`Error posting comment: ${e}`);
      res.json(createErrorResponse('Error posting comment')).status(500);
    }
  });
  app.patch('/comment', async (req, res) => {
    log.trace(req);
    try {
      const comment = req.body as IComment;
      const result = await updateComment(comment);
      if (result) res.json(createSuccessResponse(result.rows[0])).status(200);
      else res.json(createErrorResponse('Comment not found')).status(404);
    } catch (e) {
      log.error(`Error posting comment: ${e}`);
      res.json(createErrorResponse('Error posting comment')).status(500);
    }
  });
  app.delete('/comment', async (req, res) => {
    log.trace(req.query, req.params);
    try {
      const { id } = req.query;
      if (typeof id === 'string') {
        const result = await deleteComment(parseInt(id));
        if (result.rowCount)
          res.json(createSuccessResponse(result.rows[0])).status(200);
        else res.json(createErrorResponse('Comment not found')).status(404);
      }
    } catch (e) {
      log.error(`Error posting comment: ${e}`);
      res.json(createErrorResponse('Error posting comment')).status(500);
    }
  });
}
