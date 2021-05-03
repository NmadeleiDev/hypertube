import { query } from './postgres';
import log from '../../logger/logger';
import {
  GenresKeys,
  IMDBMovie,
  IDBMovie,
  IFrontComment,
  IMovie,
  IKinopoiskMovie,
  ITranslatedMovie,
  IDBTorrent,
} from '../../model/model';
import { selectCommentsByMovieID } from './comments';
import axios from 'axios';
const POSTGRES_SCHEME = process.env.POSTGRES_SCHEME || 'hypertube';

export type IRating = 1 | 2 | 3 | 4 | 5;

const selectMovieFromDB = async (
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

const selectMoviesByLetter = async (
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
      // `SELECT id as imdbid, * FROM movies WHERE title LIKE '${letter}%'
      // order by title limit $1 offset $2`,
      [limit, offset]
    );
    log.debug('[selectMoviesByLetter] selected movies', res.rows);
    return res.rowCount ? res.rows : [];
  } catch (e) {
    log.error(e);
    return null;
  }
};

const updateUserVoteList = async (
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

const getTorrentFromDB = async (movieId: string) => {
  log.debug(movieId);
  const res = await query(
    `SELECT * FROM ${POSTGRES_SCHEME}.torrents WHERE movieid=$1`,
    [movieId]
  );
  log.debug(res.rowCount);
  return res.rowCount > 0;
};

const isUserVoted = async (movieId: string, userId: number) => {
  if (Number.isNaN(userId)) throw new Error('userId is not valid');
  log.debug(userId);
  const res = await query(
    `SELECT * FROM ${POSTGRES_SCHEME}.user_ratings WHERE userid=$1 and movieid=$2`,
    [userId, movieId]
  );
  log.debug(res.rowCount);
  return res.rowCount > 0;
};

const getUserIdFromToken = (token: string) => {
  return +JSON.parse(
    Buffer.from(
      Buffer.from(token, 'base64').toString().split('.')[0],
      'base64'
    ).toString()
  ).userId;
};

const isIMovie = (movie: IMovie | unknown): movie is IMovie => {
  return (
    (movie as IMovie).id !== undefined && (movie as IMovie).info !== undefined
  );
};

const loadKinopoiskTranslation = async (en: IMovie): Promise<IMovie> => {
  try {
    const host = process.env.SEARCH_API_HOST || 'localhost';
    const res = await axios(
      `http://${host}:${process.env.SEARCH_API_PORT}/translate`,
      {
        params: {
          imdbid: en.id,
          title: en.title,
        },
      }
    );
    log.info('[getMovieInfo] got russian translation', res.data);
    if (!res.data.status) return en;
    if (!isIMovie(res.data.data))
      throw new Error('[loadKinopoiskTranslation] wrong return type');
    return res.data.data;
  } catch (e) {
    log.error(e);
    return null;
  }
};

// ========== exported =============

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

export const saveMovieToDB = async (movie: IMDBMovie) => {
  log.debug('[saveMovieToDB]', movie);
  try {
    query(
      `INSERT INTO ${POSTGRES_SCHEME}.movies (id, title, image, year, genres, rating, views,
runtimeMins, contentRating, countries, plot, directors, directorList, stars, actorList, keywordList, images, imDbRating)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
ON CONFLICT (id) DO UPDATE SET (title, image, year, genres,
runtimeMins, contentRating, countries, plot, directors, directorList, stars, actorList, keywordList, images, imDbRating)
= ($2, $3, $4, $5, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) WHERE movies.id = $1`,
      [
        movie.id,
        movie.title,
        movie.image,
        movie.year,
        movie.genres,
        '0',
        '0',
        movie.runtimeMins || '0',
        movie.contentRating || null,
        movie.countries || null,
        movie.plot,
        movie.directors,
        movie.directorList ? JSON.stringify(movie.directorList) : null,
        movie.stars,
        movie.actorList ? JSON.stringify(movie.actorList) : null,
        movie.keywordList?.join(', '),
        movie.images ? JSON.stringify(movie.images) : null,
        movie.imDbRating,
      ]
    );
  } catch (e) {
    log.debug(movie);
    log.error(e);
  }
};

export const updateMovieInDB = async (movie: IMDBMovie) => {
  log.debug('[updateMovieInDB]', movie);
  try {
    const res = await query(`SELECT * FROM movies WHERE id=${movie.id}`);
    const prevValue = res.rows[0];
    query(
      `UPDATE ${POSTGRES_SCHEME}.movies SET (title, image, year, genres,
runtimeMins, contentRating, countries, plot, directors, directorList, stars, actorList, keywordList, images, imDbRating)
= ($2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
WHERE movies.id = $1`,
      [
        movie.id,
        movie.title,
        movie.image,
        movie.year,
        movie.genres,
        movie.runtimeMins || prevValue.runtimeMins,
        movie.contentRating || prevValue.contentRating,
        movie.countries || prevValue.countries,
        movie.plot,
        movie.directors,
        movie.directorList
          ? JSON.stringify(movie.directorList)
          : prevValue.directorList,
        movie.stars,
        movie.actorList ? JSON.stringify(movie.actorList) : prevValue.actorList,
        movie.keywordList
          ? movie.keywordList.join(', ')
          : prevValue.keywordList,
        movie.images ? JSON.stringify(movie.images) : prevValue.images,
        movie.imDbRating,
      ]
    );
  } catch (e) {
    log.debug(movie);
    log.error(e);
  }
};

export const getMovieInfo = async (
  movieid: string
): Promise<ITranslatedMovie> => {
  log.debug('[getMovieInfo]', movieid);
  if (!movieid) throw new Error('comment id is missing');
  try {
    let movie = null;
    const enMovie = await selectMovieFromDB(movieid);
    if (!enMovie) throw new Error('En movie not found');
    const en = dbToIMovie(enMovie);
    let ruMovie = await getKinopoiskMovieFromDB(movieid);
    if (ruMovie) {
      const ru = KinoDBToIMovie(ruMovie);
      movie = { en, ru };
    } else {
      const ru = await loadKinopoiskTranslation(en);
      log.info('[getMovieInfo] loaded russian translation', ru);

      if (!ru) movie = { en, ru: en };
      movie = { en, ru: ru };
    }
    // const comments = await selectCommentsByMovieID(movieid);
    log.info('[getMovieInfo] full movie', movie);
    return movie;
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

export const getMovies = async (
  limit: number = 5,
  offset: number = 0,
  letter?: string
): Promise<IDBMovie[] | null> => {
  try {
    if (letter) return await selectMoviesByLetter(letter, limit, offset);
    else return await selectMovies(limit, offset);
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

export const KinoDBToIMovie = (movie: IKinopoiskMovie): IMovie => {
  return {
    id: movie.kinoid,
    title: movie.nameru,
    img: movie.posterurlpreview,
    src: '',
    info: {
      avalibility: 0,
      year: 0,
      genres: [],
      rating: 0,
      views: 0,
      imdbRating: 0,
      length: 0,
      pgRating: null,
      countries: [],
      description: movie.description || '',
    },
  };
};

export const dbToIMovie = (
  row: IDBMovie,
  comments?: IFrontComment[],
  torrent?: IDBTorrent
): IMovie => {
  log.trace('[dbToIMovie]', row);
  const defaultNumberOfCommentsToLoad = 5;
  const avalibility = torrent ? torrent.seeds + torrent.peers * 0.1 : 0;
  return {
    id: row.imdbid,
    title: row.title,
    img: row.image,
    src: '',
    info: {
      avalibility: avalibility,
      year: +row.year,
      genres: (row.genres.split(', ') as unknown) as GenresKeys[],
      imdbRating: +row.imdbrating,
      rating: +row.rating,
      views: +row.views,
      length: +row.runtimemins,
      pgRating: row.contentrating || 'N/A',
      countries: row.countries?.split(', ') || [],
      description: row.plot || '',
      directors: row.directors,
      directorList: JSON.parse(row.directorlist),
      stars: row.stars,
      cast: JSON.parse(row.actorlist),
      keywords: row.keywordlist?.split(','),
      photos: row.images ? JSON.parse(row.images) : undefined,
      comments: comments?.slice(0, defaultNumberOfCommentsToLoad),
      maxComments: comments ? comments.length : +row.maxcomments || 0,
    },
  };
};
