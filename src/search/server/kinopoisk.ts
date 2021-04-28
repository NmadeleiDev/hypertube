/**
 * API токен: fae28196-78fe-48e1-a552-99b006827fff

Для доступа, запрос должен содержать header
name: X-API-KEY
value: fae28196-78fe-48e1-a552-99b006827fff

(\w+)\*\s(string|integer)
$1: $2;

(\w+)\s(string|integer)
$1?: $2;

example: .*\n
 */

import axios from 'axios';
import { InsertOneWriteOpResult } from 'mongodb';
import {
  getKinopoiskMovieFromDB,
  saveKinopoiskMovieToDB,
} from '../db/postgres/postgres';
import log from '../logger/logger';
import { IKinopoiskMovie, IMovie } from '../model/model';
import { getIMDBInfo } from './imdb';
import { ITorrent } from './torrents';

const api = axios.create({
  baseURL: 'https://kinopoiskapiunofficial.tech/api/v2.1/films/',
  withCredentials: true,
  headers: {
    'X-API-KEY': 'fae28196-78fe-48e1-a552-99b006827fff',
  },
  validateStatus: (status) => status >= 200 && status < 500,
});

interface IKinopoiskGenre {
  genre: string;
}
interface IKinopoiskCountry {
  country: string;
}

interface IKinopoiskEpisodes {
  seasonNumber: number;
  episodeNumber: number;
  nameRu: string;
  nameEn: string;
  synopsis: string;
  releaseDate: string;
}

interface IKinopoiskSeason {
  number: number;
  episodes: IKinopoiskEpisodes[];
}

interface IKinopoiskRating {
  rating: number;
  ratingVoteCount: number;
  ratingImdb: number;
  ratingImdbVoteCount: number;
  ratingFilmCritics: string;
  ratingFilmCriticsVoteCount: number;
  ratingAwait?: string;
  ratingAwaitCount?: number;
  ratingRfCritics?: string;
  ratingRfCriticsVoteCount?: number;
}

interface IKinopoiskBudget {
  grossRu?: number;
  grossUsa?: number;
  grossWorld?: number;
  budget?: string;
  marketing?: number;
}

interface IKinopoiskReview {
  reviewsCount?: number;
  ratingGoodReview?: string;
  ratingGoodReviewVoteCount?: number;
}

interface IKinopoiskExternalId {
  imdbId: string;
}

interface IKinopoiskImage {
  language: string;
  url: string;
  height: number;
  width: number;
}

interface IKinopoiskImages {
  posters: IKinopoiskImage[];
  backdrops: IKinopoiskImage[];
}

export interface IKinopoiskCommonFilmData {
  filmId: number;
  nameRu: string;
  nameEn: string;
  webUrl: string;
  posterUrl: string;
  posterUrlPreview: string;
  year: string;
  filmLength: string;
  slogan: string;
  description: string;
  type?: 'FILM' | 'TV_SHOW' | 'UNKNOWN';
  ratingMpaa: string;
  ratingAgeLimits: number;
  premiereRu: string;
  distributors: string;
  premiereWorld: string;
  premiereDigital?: string;
  premiereWorldCountry: string;
  premiereDvd: string;
  premiereBluRay: string;
  distributorRelease: string;
  countries: IKinopoiskCountry[];
  genres: IKinopoiskGenre[];
  facts: string[];
  seasons: IKinopoiskSeason[];
}

interface IKinopoiskSearchFilmData {
  filmId: number;
  nameRu: string;
  nameEn: string;
  type: string;
  year: string;
  description: string;
  filmLength: string;
  countries: IKinopoiskCountry[];
  genres: IKinopoiskGenre[];
  rating: string;
  ratingVoteCount: number;
  posterUrl: string;
  posterUrlPreview: string;
}

export interface IKinopoiskFilmResponse {
  data: IKinopoiskCommonFilmData;
  rating: IKinopoiskRating;
  budget: IKinopoiskBudget;
  review: IKinopoiskReview;
  externalId: IKinopoiskExternalId;
  images: IKinopoiskImages;
}

export interface IKinopoiskSearchFilmResponse {
  keyword: string;
  pagesCount: number;
  films: IKinopoiskSearchFilmData[];
  searchFilmsCountResult: number;
}

export interface ITranslatedMovie {
  en: IMovie;
  ru: IMovie;
}

const searchKinopoiskMovie = async (
  keyword: string
): Promise<IKinopoiskSearchFilmData[] | null> => {
  try {
    const res = await api('search-by-keyword', {
      params: {
        keyword,
        page: 1,
      },
    });
    if (!res || !res.data) throw new Error('Kinopoisk fetch error');
    else if (res.status === 401) throw new Error('Empty access token');
    else if (res.status === 404) throw new Error('Movie not found');
    else if (res.status === 429) {
      // wait a sec, too many requests
      setTimeout(() => {
        return searchKinopoiskMovie(keyword);
      }, 1000);
    }
    const data = res.data as IKinopoiskSearchFilmResponse;
    log.debug('searchKinopoiskMovie', data);
    if (!data.searchFilmsCountResult) throw new Error('Zero movies found');
    return data.films;
  } catch (e) {
    log.error(e);
    return null;
  }
};

const getKinopoiskMovieById = async (
  filmId: number
): Promise<IKinopoiskFilmResponse | null> => {
  try {
    log.debug('[getKinopoiskMovieById] filmId', filmId);
    const res = await api(`${filmId}`);
    if (!res || !res.data) throw new Error('Kinopoisk fetch error');
    else if (res.status === 401) throw new Error('Empty access token');
    else if (res.status === 404) throw new Error('Movie not found');
    // else if (res.status === 429) {
    //   // wait a sec, too many requests
    //   setTimeout(() => {
    //     return getKinopoiskMovieById(filmId);
    //   }, 1000);
    // }
    const data = res.data as IKinopoiskFilmResponse;
    log.debug('getKinopoiskMovieById', data);
    return data;
  } catch (e) {
    log.error(e);
    return null;
  }
};

export const kinopoiskToIMovie = (
  movie: IKinopoiskFilmResponse,
  torrent?: ITorrent
): IMovie => {
  return {
    id: movie.data.filmId + '',
    title: movie.data.nameRu,
    img: movie.data.posterUrlPreview,
    src: '',
    info: {
      avalibility: torrent?.torrent.seeds || 0,
      year: +movie.data.year,
      genres: [],
      rating: 0,
      views: 0,
      imdbRating: 0,
      length: 0,
      pgRating: null,
      countries: [],
      description: movie.data.description || '',
    },
  };
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

export const getKinopoiskMovieByImdbid = async (
  title: string,
  imdbId: string
) => {
  try {
    const movies = await searchKinopoiskMovie(title);
    if (!movies || !movies.length) return null;
    log.debug(
      '[translateMovie] found similar movies for a title',
      title,
      movies
    );
    for (let movie of movies) {
      const kinopoiskMovie = await getKinopoiskMovieById(movie.filmId);
      if (!kinopoiskMovie) return null;
      if (kinopoiskMovie.externalId.imdbId !== imdbId) continue;
      saveKinopoiskMovieToDB(kinopoiskMovie);
      return kinopoiskToIMovie(kinopoiskMovie);
    }
    return null;
  } catch (e) {
    log.error(e);
    return null;
  }
};

export const translateMovie = async (
  en: IMovie
): Promise<ITranslatedMovie | null> => {
  try {
    const title = en.title;
    const imdbId = en.id;

    const movie = await getKinopoiskMovieFromDB(imdbId);
    log.info('getKinopoiskMovieFromDB result:', movie);
    if (movie) {
      return { en, ru: KinoDBToIMovie(movie) };
    }
    log.info('Cannot find matching movie in DB, searching in web');
    let ru = await getKinopoiskMovieByImdbid(title, imdbId);
    if (!ru) ru = en;
    return { en, ru };
  } catch (e) {
    log.error(e);
    return null;
  }
};
