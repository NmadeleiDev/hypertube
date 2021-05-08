import axios from 'axios';
import {
  getMovieFromDB,
  insertTorrentIntoLoadedFiles,
  saveMovieToDB,
  saveTorrentInDB,
} from '../db/postgres/postgres';
import log from '../logger/logger';
import { GenresKeys, IDBMovie, IMovie } from '../model/model';
import { serachIMDBMovie } from './imdbApi';
import { ITranslatedMovie, translateMovie } from './kinopoisk';
import { ITorrent } from './torrents';
const nameToImdb = require('name-to-imdb');

export interface IMDBPerson {
  id: string;
  name: string;
  image?: string;
  asCharacter?: string;
}

interface nameToImdbSearch {
  match: 'imdbFind';
  meta: {
    id: string;
    name: string;
    year: number;
    type: 'feature';
    yearRange: string | number | undefined;
    image: {
      src: string;
      width: number;
      height: number;
    };
    starring: string;
    similarity: number;
  };
}

export interface IMDBMovie {
  id: string;
  title: string;
  type?: string;
  year: string;
  image: string;
  runtimeMins?: string;
  plot: string;
  genres: string;
  countries?: string;
  contentRating?: string;
  imDbRating: string;
  directors: string;
  directorList?: IMDBPerson[];
  stars: string;
  actorList?: IMDBPerson[];
  keywordList?: string[];
  images?: string[];
  similars?: IMDBMovie[];
  errorMessage?: string;
}

const isInfoFull = ({
  id,
  title,
  year,
  image,
  plot,
  imDbRating,
  directors,
  stars,
}: IMDBMovie) => {
  return !!(
    id?.length &&
    title?.length &&
    year?.length &&
    image?.length &&
    plot?.length &&
    imDbRating?.length &&
    directors?.length &&
    stars?.length
  );
};

export const removeDuplicates = (movies: ITranslatedMovie[]) => {
  return movies.reduce((acc: ITranslatedMovie[], cur) => {
    return !cur || acc.find((movie) => movie.en.id === cur.en.id)
      ? acc
      : [...acc, cur];
  }, []);
};

export const getIMDBInfo = async (
  torrent: ITorrent
): Promise<IMDBMovie | null> => {
  try {
    const movie = await searchMovieInIMDB(torrent);
    log.debug('[getIMDBInfo] search movie result', movie);
    if (movie) {
      if (movie.similars?.length) {
        movie.similars.forEach((item) => saveMovieToDB(item));
      }
      return movie;
    }
    return await getNameToIMDBInfo(torrent.movieTitle);
  } catch (e) {
    log.error(e);
    throw new Error('Error getting IMDB info');
  }
};

const isIMDBMovie = (movie: IMDBMovie | unknown): movie is IMDBMovie => {
  return !!(
    movie &&
    (movie as IMDBMovie).id &&
    (movie as IMDBMovie).title &&
    (movie as IMDBMovie).year &&
    (movie as IMDBMovie).image &&
    (movie as IMDBMovie).plot &&
    (movie as IMDBMovie).genres
  );
};

const searchMovieInIMDB = async (torrent: ITorrent) => {
  try {
    let id = torrent.torrent.imdb || '';
    if (id === '') {
      const search = await axios(
        `https://imdb-api.com/en/API/Search/${process.env.IMDB_API_KEY}/${torrent.movieTitle} ${torrent.year}`,
        { validateStatus: (status) => status >= 200 && status < 500 }
      );
      log.debug('[searchMovieInIMDB] SEARCH result', search.data);
      if (search.data.results && search.data.results.length) {
        id = search.data.results[0]?.id;
      }
      if (isInfoFull(search.data.results[0]))
        return search.data.results[0] as IMDBMovie;
    }
    const movie = await axios(
      `https://imdb-api.com/en/API/Title/${process.env.IMDB_API_KEY}/${id}`,
      { validateStatus: (status) => status >= 200 && status < 500 }
    );
    log.debug('[searchMovieInIMDB] TITLE result', movie.data);
    if (!movie.data?.id)
      throw new Error('[searchMovieInIMDB] IMDB movie not found');
    if (isIMDBMovie(movie.data)) return movie.data as IMDBMovie;
    log.error('malformed imdb-api.com response', movie.data);
    throw new Error('[searchMovieInIMDB] malformed imdb-api.com response');
  } catch (e) {
    log.error(e);
    return null;
  }
};

const getNameToIMDBInfo = async (title: string): Promise<IMDBMovie | null> => {
  try {
    const res = await nameToImdbPromise(title);
    log.debug('[getIMDBInfo] IMDBMovie', res);
    return {
      id: res.meta.id,
      title: res.meta.name,
      year: `${res.meta.year}`,
      image: res.meta.image.src,
      plot: '',
      genres: '',
      imDbRating: '',
      directors: '',
      stars: '',
    } as IMDBMovie;
  } catch (e) {
    return null;
  }
};

const nameToImdbPromise = (title: string): Promise<nameToImdbSearch> => {
  return new Promise(function (resolve, reject) {
    nameToImdb(title, (err, res, inf) => {
      if (!res) reject(err);
      resolve(inf);
    });
  });
};

export const imdbToIMovie = (
  movie: IMDBMovie | null,
  torrent: ITorrent
): IMovie => {
  log.trace('[imdbToIMovie]', movie);
  if (!movie) return null;
  return {
    id: movie.id,
    title: movie.title,
    img: movie.image,
    src: '',
    isViewed: false,
    info: {
      avalibility: torrent?.torrent.seeds || 0,
      year: +movie.year,
      genres: movie.genres?.split(', ') as GenresKeys[],
      rating: 0,
      views: 0,
      imdbRating: +movie.imDbRating,
      length: +movie.runtimeMins || 0,
      pgRating: movie.contentRating || 'N/A',
      countries: movie.countries?.split(', '),
      description: movie.plot || '',
      directors: movie.directors,
      directorList: movie.directorList,
      stars: movie.stars,
      cast: movie.actorList,
      keywords: movie.keywordList,
      photos: movie.images || undefined,
    },
  };
};

export const dbToIMovie = (row: IDBMovie, torrent?: ITorrent): IMovie => {
  log.trace('[dbToIMovie]', row);
  return {
    id: row.id,
    title: row.title,
    img: row.image,
    src: '',
    isViewed: !!row.isviewed,
    info: {
      avalibility: +row.avalibility || torrent?.torrent.seeds || 0,
      year: +row.year,
      genres: (row.genres.split(', ') as unknown) as GenresKeys[],
      rating: +row.rating,
      views: +row.views,
      imdbRating: +row.imdbrating,
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
    },
  };
};

export const loadMoviesInfo = (
  torrents: ITorrent[]
): Promise<ITranslatedMovie[] | null> => {
  return Promise.all(
    torrents.map(async (torrent) => {
      let en: IMovie = null;
      const dbInfo = await getMovieFromDB(
        torrent.movieTitle,
        torrent.torrent.imdb,
        torrent.year
      );
      if (dbInfo) {
        en = dbToIMovie(dbInfo, torrent);
        log.debug('movie found in DB:', en);
      } else {
        log.info(
          `${torrent.movieTitle} was not found in database, now search in IMDB`
        );
        let IMDBInfo = await serachIMDBMovie(
          `${torrent.movieTitle} ${torrent.year}`
        );
        if (!IMDBInfo?.id) {
          IMDBInfo = await getIMDBInfo(torrent);
          // saveMovieToDB(IMDBInfo);
          // en = imdbToIMovie(IMDBInfo, torrent);
        }
        await saveMovieToDB(IMDBInfo);
        en = imdbToIMovie(IMDBInfo, torrent);
      }
      if (!en) return null;
      torrent.torrent.imdb = en.id;
      saveTorrentInDB(torrent);
      insertTorrentIntoLoadedFiles(torrent);
      const movie = await translateMovie(en);
      log.info('Translated movie', movie);
      return movie;
    })
  );
};
