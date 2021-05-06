import log from '../logger/logger';
import {
  getKinopoiskMovieFromDB,
  selectMovieFromDB,
  selectMovies,
  selectMoviesByGenre,
  selectMoviesByLetter,
} from '../db/postgres/movies';
import {
  GenresKeys,
  IDBMovie,
  IDBTorrent,
  IFrontComment,
  IKinopoiskMovie,
  IMovie,
  ITranslatedMovie,
} from '../model/model';
import axios from 'axios';
import { isIMovie } from './utils';

interface IMoviesQuery {
  letter?: string;
  genre?: string;
  limit: number;
  offset: number;
}

export const getTranslatedMovies = async ({
  letter,
  genre,
  limit,
  offset,
}: IMoviesQuery) => {
  const response = await getMovies(limit, offset, letter, genre);
  log.debug('[getTranslatedMovies] getMovies response', response);
  const ens: IMovie[] = response.map((movie) => dbToIMovie(movie));
  log.info(`[getTranslatedMovies] en movies for query ${genre}`, ens);

  const movies: ITranslatedMovie[] = await Promise.all(
    ens.map((en) => getMovieInfo(en.id))
  );
  log.info('[getTranslatedMovies] got translated movies', movies);
  return movies;
};

export const loadKinopoiskTranslation = async (en: IMovie): Promise<IMovie> => {
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

export const getMovies = async (
  limit: number = 5,
  offset: number = 0,
  letter?: string,
  genre?: string
): Promise<IDBMovie[] | null> => {
  try {
    if (letter) return await selectMoviesByLetter(letter, limit, offset);
    else if (genre) return await selectMoviesByGenre(genre, limit, offset);
    else return await selectMovies(limit, offset);
  } catch (e) {
    log.error(e);
    return null;
  }
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
