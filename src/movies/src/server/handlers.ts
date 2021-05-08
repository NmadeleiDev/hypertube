import { Express } from 'express';
import { createSuccessResponse, createErrorResponse } from './utils';
import log from '../logger/logger';
import { IComment } from '../model/model';
import {
  deleteComment,
  insertComment,
  selectCommentById,
  selectCommentsByMovieID,
  updateComment,
} from '../db/postgres/comments';
import {
  IRating,
  updateMovieRating,
  updateMovieViews,
  updateUserMovieViews,
} from '../db/postgres/movies';
import { getMovieInfo, getTranslatedMovies } from './movies';

export function addMoviesHandlers(app: Express) {
  app.get('/movies', async (req, res) => {
    log.trace(req);
    const limit = +req.query.limit || 5;
    const offset = +req.query.offset || 0;
    const id = (req.query.id as string) || null;
    const headers = req.headers;

    try {
      const token = (headers.accesstoken as string) || '';
      if (id) {
        const movie = await getMovieInfo(id, token);
        log.debug('[GET /movies] movieById:', movie);
        if (movie) res.status(200).json(createSuccessResponse([movie]));
        else res.status(404).json(createErrorResponse(null));
        return;
      }
      const movies = await getTranslatedMovies({ token, limit, offset });
      log.info('[/movies] got translated movies', movies);

      if (movies) res.status(200).json(createSuccessResponse(movies));
      else res.status(404).json(createErrorResponse(null));
    } catch (e) {
      log.error(`Error getting movies: ${e}`);
      res.status(500).json(createErrorResponse('Error getting movies'));
    }
  });
  app.get('/bygenre', async (req, res) => {
    const genre = (req.query.genre as string) || undefined;
    const limit = +req.query.limit || 5;
    const offset = +req.query.offset || 0;
    const headers = req.headers;
    log.debug('genre:', genre);
    try {
      const token = (headers.accesstoken as string) || '';
      const movies = await getTranslatedMovies({ token, limit, offset, genre });
      log.info('[/bygenre] got translated movies', movies);

      if (movies) res.status(200).json(createSuccessResponse(movies));
      else res.status(404).json(createErrorResponse(null));
    } catch (e) {
      log.error(`Error getting movies: ${e}`);
      res.status(500).json(createErrorResponse('Error getting movies'));
    }
  });
  app.get('/byname', async (req, res) => {
    const limit = +req.query.limit || 5;
    const offset = +req.query.offset || 0;
    const letter = (req.query.letter as string) || undefined;
    const headers = req.headers;
    log.debug('letter:', letter);
    try {
      const token = (headers.accesstoken as string) || '';
      const movies = await getTranslatedMovies({
        token,
        limit,
        offset,
        letter,
      });
      log.info('[/byname] got translated movies', movies);

      if (movies) res.status(200).json(createSuccessResponse(movies));
      else res.status(404).json(createErrorResponse(null));
    } catch (e) {
      log.error(`Error getting movies: ${e}`);
      res.status(500).json(createErrorResponse('Error getting movies'));
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
      const token = (headers.accesstoken as string) || '';
      log.debug(`rating: ${rating}, movieid: ${movieid}, token: ${token}`);

      if (rating === null || !movieid || !token)
        return res
          .status(400)
          .json(createErrorResponse('accessToken, ID and RATING are required'));

      const newRating = await updateMovieRating(movieid, rating, token);
      res.status(200).json(createSuccessResponse(newRating));
    } catch (e) {
      log.error(`Error getting movies: ${e}`);
      res.status(500).json(createErrorResponse('Error getting movies'));
    }
  });
  app.patch('/views', async (req, res) => {
    const headers = req.headers;
    const cookies = req.cookies;
    log.debug(req.body, headers, cookies);
    try {
      const token = (headers.accesstoken as string) || '';
      const movieid = (req.body?.movieId as string) || null;
      log.debug(`movieid: ${movieid}`);
      if (!movieid || !token)
        return res
          .status(400)
          .json(createErrorResponse('accessToken and movieId are required'));

      updateUserMovieViews(movieid, token);
      updateMovieViews(movieid).then((views) => {
        res.status(200).json(createSuccessResponse(views));
      });
    } catch (e) {
      log.error(`Error getting movies: ${e}`);
      res.status(500).json(createErrorResponse('Error getting movies'));
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
      res.status(200).json(createSuccessResponse(comments));
    } catch (e) {
      log.error(`Error getting comments: ${e}`);
      res.status(500).json(createErrorResponse('Error getting comments'));
    }
  });
  app.post('/comment', async (req, res) => {
    log.trace(req);
    try {
      const comment = req.body as IComment;
      const newComment = (await insertComment(comment)).rows[0] as IComment;
      if (!newComment)
        return res
          .status(500)
          .json(createErrorResponse('Error posting comment'));
      const result = await selectCommentById(newComment.id);
      log.debug(result);
      res.status(200).json(createSuccessResponse(result));
    } catch (e) {
      log.error(`Error posting comment: ${e}`);
      res.status(500).json(createErrorResponse('Error posting comment'));
    }
  });
  app.patch('/comment', async (req, res) => {
    log.trace(req);
    try {
      const comment = req.body as IComment;
      const result = await updateComment(comment);
      if (result) res.status(200).json(createSuccessResponse(result.rows[0]));
      else res.status(404).json(createErrorResponse('Comment not found'));
    } catch (e) {
      log.error(`Error posting comment: ${e}`);
      res.status(500).json(createErrorResponse('Error posting comment'));
    }
  });
  app.delete('/comment', async (req, res) => {
    log.trace(req.query, req.params);
    try {
      const { id } = req.query;
      if (typeof id === 'string') {
        const result = await deleteComment(parseInt(id));
        if (result.rowCount)
          res.status(200).json(createSuccessResponse(result.rows[0]));
        else res.status(404).json(createErrorResponse('Comment not found'));
      }
    } catch (e) {
      log.error(`Error posting comment: ${e}`);
      res.status(500).json(createErrorResponse('Error posting comment'));
    }
  });
}
