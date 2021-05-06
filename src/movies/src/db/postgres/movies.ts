import { query } from './postgres';
import log from '../../logger/logger';
import { IDBMovie, IKinopoiskMovie } from '../../model/model';
import { getUserIdFromToken } from '../../server/utils';
const POSTGRES_SCHEME = process.env.POSTGRES_SCHEME || 'hypertube';

export type IRating = 1 | 2 | 3 | 4 | 5;

export const selectMovieFromDB = async (
  id: string,
  title: string = ''
): Promise<IDBMovie | null> => {
  try {
    if (!id && !title)
      throw new Error('[selectMovieFromDB] Both movieID and title are missing');
    const res = id
      ? await query(
          `with c as (select count(c.id) maxcomments from ${POSTGRES_SCHEME}.comments c where movieid=$1)
          SELECT m.id as imdbid, c.*, m.* FROM ${POSTGRES_SCHEME}.movies m, c WHERE m.id=$1`,
          [id]
        )
      : await query(
          `SELECT id as imdbid, * FROM ${POSTGRES_SCHEME}.movies WHERE title LIKE '%${title}%'`
        );
    log.debug('[selectMovieFromDB]', res.rows);
    return res.rowCount ? res.rows[0] : [];
  } catch (e) {
    log.error(e);
    return null;
  }
};

export const selectMoviesByLetter = async (
  letter: string,
  limit: number = 5,
  offset: number = 0
): Promise<IDBMovie[] | null> => {
  try {
    if (!letter) throw new Error('[selectMoviesByLetter] letter is missing');
    const res = await query(
      `select
        m.id as imdbid, m.*, count(c.*) maxComments
      from
      ${POSTGRES_SCHEME}.movies m
      join ${POSTGRES_SCHEME}.torrents t on
        t.movieid = m.id
      left join ${POSTGRES_SCHEME}.comments c on
        m.id = c.movieid 
      where
        m.title like '${letter}%' and 
        ( t.magnet is not null
        or t.torrent is not null )
      group by m.id
      order by
        m.rating
      limit $1 offset $2;`,
      [limit, offset]
    );
    log.debug('[selectMoviesByLetter] selected movies', res.rows);
    return res.rowCount ? res.rows : [];
  } catch (e) {
    log.error(e);
    return null;
  }
};

export const selectMoviesByGenre = async (
  genre: string,
  limit: number = 5,
  offset: number = 0
): Promise<IDBMovie[] | null> => {
  try {
    const res = await query(
      `select
        m.id as imdbid, m.*, count(c.*) maxComments
      from
      ${POSTGRES_SCHEME}.movies m
      join ${POSTGRES_SCHEME}.torrents t on
        t.movieid = m.id
      left join ${POSTGRES_SCHEME}.comments c on
        m.id = c.movieid 
      where
        m.genres like '%${genre}%' and 
        ( t.magnet is not null
        or t.torrent is not null )
      group by m.id
      order by
        m.rating
      limit $1 offset $2;`,
      [limit, offset]
    );
    log.debug('[selectMoviesByLetter] selected movies', res.rows);
    return res.rowCount ? res.rows : [];
  } catch (e) {
    log.error(e);
    return null;
  }
};

export const updateUserVoteList = async (
  movieId: string,
  vote: IRating,
  userId: number
) => {
  log.debug(userId);
  const res = await query(
    `INSERT INTO ${POSTGRES_SCHEME}.user_ratings (movieid, userid, vote) VALUES ($1, $2, $3)
    ON CONFLICT (movieid, userid) DO UPDATE SET (movieid, userid, vote) = ($1, $2, $3) WHERE user_ratings.userid=$2`,
    [movieId, userId, vote]
  );
  log.debug(res.rowCount);
  return res.rowCount > 0;
};

export const isUserVoted = async (movieId: string, userId: number) => {
  if (Number.isNaN(userId)) throw new Error('userId is not valid');
  log.debug(userId);
  const res = await query(
    `SELECT * FROM ${POSTGRES_SCHEME}.user_ratings WHERE userid=$1 and movieid=$2`,
    [userId, movieId]
  );
  log.debug(res.rowCount);
  return res.rowCount > 0;
};

export const updateMovieRating = async (
  movieId: string,
  vote: IRating,
  token: string
): Promise<string | null> => {
  try {
    const userId = getUserIdFromToken(token);
    if (!userId)
      throw new Error(`Can't parse userId: ${userId}, token: ${token}`);
    const movie = await selectMovieFromDB(movieId);
    if (!movie) throw new Error(`MovieID ${movieId} not found`);
    // const userVoted = await isUserVoted(movieId, userId);
    // log.debug('userVoted:', userVoted, userId);
    // if (userVoted) return null;
    updateUserVoteList(movieId, vote, userId);
    const newRating = (
      (+movie.rating * +movie.votes + vote) /
      (+movie.votes + 1)
    ).toFixed(8);
    query(
      `UPDATE ${POSTGRES_SCHEME}.movies SET (rating, votes) = ($2, $3) WHERE id = $1 RETURNING *`,
      [movieId, newRating, +movie.votes + 1]
    );
    return newRating;
  } catch (e) {
    log.error(e);
    return null;
  }
};

export const selectMovies = async (limit: number = 5, offset: number = 0) => {
  try {
    const res = await query(
      `SELECT
        m.id as imdbid, m.*, count(c.*) maxComments
      FROM
      ${POSTGRES_SCHEME}.movies m
      JOIN ${POSTGRES_SCHEME}.torrents t on
        t.movieid = m.id
      left join ${POSTGRES_SCHEME}.comments c on
        m.id = c.movieid 
      where
        t.magnet is not null
        or t.torrent is not null
      group by m.id
      order by
        m.rating
      limit $1 offset $2;`,
      [limit, offset]
    );
    if (!res.rowCount) log.info('No movies with saved torrents found');
    return res.rows;
  } catch (e) {
    log.error(e);
    return null;
  }
};

export const getKinopoiskMovieFromDB = async (
  id: string
): Promise<IKinopoiskMovie | null> => {
  try {
    const res = await query(
      `SELECT id as kinoid, * FROM ${POSTGRES_SCHEME}.kinopoisk WHERE imdbid=$1`,
      [id]
    );
    if (res.rowCount) return res.rows[0] as IKinopoiskMovie;
  } catch (e) {
    log.error(e);
  }
  return null;
};
