import { query } from './postgres';
import log from '../../logger/logger';

const POSTGRES_USER = process.env.POSTGRES_USER || 'admin';
const POSTGRES_PASSWORD = process.env.POSTGRES_PASSWORD || 'passwd';
const POSTGRES_HOST = process.env.POSTGRES_HOST || 'localhost';
const POSTGRES_PORT = process.env.POSTGRES_PORT || '5432';
const POSTGRES_DB = process.env.POSTGRES_DB || 'hypertube';
const POSTGRES_SCHEME = process.env.POSTGRES_SCHEME || 'hypertube';

export const DSN = `postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}`;

export const initDatabase = () => {
  const createMoviesTable = `CREATE TABLE IF NOT EXISTS ${POSTGRES_SCHEME}.movies (
    id              VARCHAR(16) PRIMARY KEY,
    title           TEXT NOT NULL,
    image           TEXT NOT NULL,
    rating          VARCHAR(10) DEFAULT '0',
    votes           VARCHAR(10) DEFAULT '0',
    views           VARCHAR(10) DEFAULT '0',
    year            VARCHAR(9) NOT NULL,
    runtimeMins     VARCHAR(10) DEFAULT '0',
    plot            TEXT DEFAULT '',
    genres          TEXT NOT NULL DEFAULT '',
    countries       TEXT,
    contentRating   TEXT,
    imDbRating      VARCHAR(10) NOT NULL,
    directors       TEXT DEFAULT '',
    directorList    TEXT,
    stars           TEXT DEFAULT '',
    actorList       TEXT,
    keywordList     TEXT,
    images          TEXT
  );`;
  const createKinopoiskMoviesTable = `CREATE TABLE IF NOT EXISTS ${POSTGRES_SCHEME}.kinopoisk (
    id              VARCHAR(16) PRIMARY KEY,
    imdbId          VARCHAR(16) REFERENCES ${POSTGRES_SCHEME}.movies(id),
    nameRu          TEXT NOT NULL,
    nameEn          TEXT NOT NULL,
    description     TEXT DEFAULT '',
    posterUrl       TEXT,
    posterUrlPreview TEXT
  );`;

  const createCommentsTable = `CREATE TABLE IF NOT EXISTS ${POSTGRES_SCHEME}.comments (
    id        SERIAL PRIMARY KEY,
    userid    INT4 NOT NULL,
    movieid   VARCHAR(16) NOT NULL REFERENCES ${POSTGRES_SCHEME}.movies(id),
    text      TEXT NOT NULL,
    time      TIMESTAMP DEFAULT current_timestamp
  );`;

  // size in GB
  const createTorrentsTable = `CREATE TABLE IF NOT EXISTS ${POSTGRES_SCHEME}.torrents (
    id        SERIAL PRIMARY KEY,
    movieid   VARCHAR(16) NOT NULL UNIQUE REFERENCES ${POSTGRES_SCHEME}.movies(id),
    torrentname TEXT NOT NULL,
    magnet    TEXT,
    torrent   BYTEA,
    seeds     INTEGER DEFAULT 0,
    peers     INTEGER DEFAULT 0,
    size      NUMERIC,
    CHECK (magnet IS NOT NULL OR torrent IS NOT NULL)
  );`;

  const createRatingTable = `CREATE TABLE IF NOT EXISTS ${POSTGRES_SCHEME}.user_ratings (
    id        SERIAL PRIMARY KEY,
    userid    INTEGER NOT NULL,
    movieid   VARCHAR(16) NOT NULL REFERENCES ${POSTGRES_SCHEME}.movies(id),
    vote    integer NOT NULL,
    UNIQUE (userid, movieid)
  );`;
  log.debug('[initDatabase]');
  // const dropTableQuery = `DROP TABLE IF EXISTS movies`;
  // await query(dropTableQuery);
  query(createMoviesTable);
  setTimeout(() => {
    query(createCommentsTable);
    query(createRatingTable);
    query(createTorrentsTable);
    query(createKinopoiskMoviesTable);
  }, 5000);
};
