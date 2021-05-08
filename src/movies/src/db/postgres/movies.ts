import { query } from './postgres';
import log from '../../logger/logger';
import { IDBMovie, IKinopoiskMovie } from '../../model/model';
import { getUserIdFromToken } from '../../server/utils';
const POSTGRES_SCHEME = process.env.POSTGRES_SCHEME || 'hypertube';

export type IRating = 1 | 2 | 3 | 4 | 5;

export const selectMovieFromDB = async ({
  movieId,
  title = '',
  token,
}: {
  movieId: string;
  title?: string;
  token: string;
}): Promise<IDBMovie | null> => {
  try {
    const userId = getUserIdFromToken(token);
    log.debug(`userId: ${userId}, token: ${token}`);
    if (!movieId && !title)
      throw new Error('[selectMovieFromDB] Both movieID and title are missing');
    const res = movieId
      ? await query(
          `with c as (select count(c.id) maxcomments from ${POSTGRES_SCHEME}.comments c where movieid=$1),
          v as (select coalesce(count(userid), 0) isviewed from ${POSTGRES_SCHEME}.views where userid='${userId}' and movieid=$1)
          SELECT m.id as imdbid, c.*, m.*, v.isviewed FROM ${POSTGRES_SCHEME}.movies m, c, v
          WHERE m.id=$1;`,
          [movieId]
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

export const selectMoviesByLetter = async ({
  letter,
  token,
  limit = 5,
  offset = 0,
}: {
  letter: string;
  limit: number;
  offset: number;
  token: string;
}): Promise<IDBMovie[] | null> => {
  try {
    const userId = getUserIdFromToken(token);
    log.debug(`userId: ${userId}, token: ${token}`);
    if (!letter) throw new Error('[selectMoviesByLetter] letter is missing');
    const res = await query(
      `with v as (
        select movieid, coalesce(count(userid), 0) from ${POSTGRES_SCHEME}.views where userid='${userId}' group by movieid
      ) select
          m.id as imdbid, m.*, count(c.*) maxComments, coalesce(count(v.*), 0) isviewed
        from
        ${POSTGRES_SCHEME}.movies m
        join ${POSTGRES_SCHEME}.torrents t on
          t.movieid = m.id
        left join ${POSTGRES_SCHEME}.comments c on
          m.id = c.movieid 
        left join v on
          m.id = v.movieid 
        where
          upper(m.title) like '${letter}%' and 
        t.torrent is not null
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

export const selectMoviesByGenre = async ({
  genre,
  token,
  limit = 5,
  offset = 0,
}: {
  genre: string;
  limit: number;
  offset: number;
  token: string;
}): Promise<IDBMovie[] | null> => {
  try {
    const userId = getUserIdFromToken(token);
    log.debug(`userId: ${userId}, token: ${token}`);
    const res = await query(
      `with v as (
        select movieid, coalesce(count(userid), 0) from ${POSTGRES_SCHEME}.views where userid='${userId}' group by movieid
      ) select
        m.id as imdbid, m.*, count(c.*) maxComments, coalesce(count(v.*), 0) isviewed
      from
      ${POSTGRES_SCHEME}.movies m
      join ${POSTGRES_SCHEME}.torrents t on
        t.movieid = m.id
      left join ${POSTGRES_SCHEME}.comments c on
        m.id = c.movieid
      left join v on
        m.id = v.movieid 
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
    log.debug(`userId: ${userId}, token: ${token}`);
    const movie = await selectMovieFromDB({ movieId, token });
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

export const updateMovieViews = async (movieId: string) => {
  try {
    const res = await query(
      `UPDATE ${POSTGRES_SCHEME}.movies m SET views=views::INTEGER + 1
      WHERE m.id=$1 RETURNING id as imdbId, views;`,
      [movieId]
    );
    if (!res.rowCount) log.warn('No movies found');
    return res.rows;
  } catch (e) {
    log.error(e);
    return null;
  }
};

export const updateUserMovieViews = async (movieId: string, token: string) => {
  try {
    const userId = getUserIdFromToken(token);
    log.debug(`userId: ${userId}, token: ${token}`);

    const res = await query(
      `INSERT INTO ${POSTGRES_SCHEME}.views (userid, movieid) VALUES ($1, $2) RETURNING userid, movieid;`,
      [userId, movieId]
    );
    return res.rows;
  } catch (e) {
    log.error(e);
    return null;
  }
};

export const selectMovies = async (
  userId: number,
  limit: number = 5,
  offset: number = 0
) => {
  try {
    const res = await query(
      `with v as (
        select movieid, coalesce(count(userid), 0) isviewed from ${POSTGRES_SCHEME}.views where userid='${userId}' group by movieid
      ) SELECT
        m.id as imdbid, m.*, count(c.*) maxComments, coalesce(count(v.*), 0)
      FROM
      ${POSTGRES_SCHEME}.movies m
      JOIN ${POSTGRES_SCHEME}.torrents t on
        t.movieid = m.id
      left join ${POSTGRES_SCHEME}.comments c on
        m.id = c.movieid
      left join v on
          m.id = v.movieid 
      where
        t.torrent is not null
      group by m.id
      order by
        m.rating
      limit $1 offset $2;`,
      [limit, offset]
    );
    if (!res.rowCount) log.warn('No movies with saved torrents found');
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
