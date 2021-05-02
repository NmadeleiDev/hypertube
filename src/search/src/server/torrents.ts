import axios from 'axios';
import { writeFileSync } from 'fs';
import parse from 'node-html-parser';
import TorrentSearchApi, { Torrent } from 'torrent-search-api';
import log from '../logger/logger';
import { loadMoviesInfo } from './imdb';
import { ITranslatedMovie } from './kinopoisk';
import { download } from './utils';
const Magnet2torrent = require('magnet2torrent-js');
const TorrentIndexer = require('torrent-indexer');
const torrentIndexer = new TorrentIndexer();

export interface ITorrent {
  torrentTitle: string; // 'name.of.a.movie.year.format.row'
  croppedTorrentTitle: string; // 'name.of.a.movie'
  movieTitle: string; // 'name of a movie'
  year: number; // 2005
  torrent: FullTorrent;
}

export interface FullTorrent extends Torrent {
  seeds?: number;
  peers?: number;
  torrent?: Buffer;
  imdb?: string;
}

const m2t = new Magnet2torrent({ timeout: 60 });

const parseDottedFormat = (torrent: Torrent): ITorrent => {
  const titleRE = /^([.\w'’:-]+?)(?=\.(\d{4})|(\d{3,4}p))/;
  const match = torrent.title.match(titleRE);
  log.debug(match);
  const name = match ? match[1] : torrent.title;
  const seriesName = name.match(/^([\w.-]+)\.s\d{1,2}\./i);
  const title = seriesName ? seriesName[1] : name;
  return {
    torrentTitle: torrent.title,
    croppedTorrentTitle: title,
    movieTitle: formatTitle(title.replace(/\./g, ' ')),
    year: match ? +match[2] : 0,
    torrent,
  };
};

const parseSpacedFormat = (torrent: Torrent): ITorrent => {
  const titleRE = /^([.\w "'’,:-]+?)\.? \(?(\d{4})\)?|(\d{3,4}p)/;
  const match = torrent.title.match(titleRE);
  log.debug('[parseSpacedFormat]', match);
  const name = match ? match[1] : torrent.title;
  const seriesName = name.match(/^([\w\s"'’,:-]+?) \(?(s\d{1,2})\)?/i);
  const title = seriesName ? seriesName[1] : name;
  return {
    torrentTitle: torrent.title,
    croppedTorrentTitle: title,
    movieTitle: formatTitle(title),
    year: match ? +match[2] : 0,
    torrent,
  };
};

export const parseSize = (size: string) => {
  return size.search('GB') !== -1
    ? Number.parseFloat(size)
    : Number.parseFloat(size) / 1000;
};

export const enableTorrentSearch = () => {
  TorrentSearchApi.enableProvider('Rarbg');
  TorrentSearchApi.enableProvider('ThePirateBay');
  // TorrentSearchApi.enablePublicProviders();
  const activeProviders = TorrentSearchApi.getActiveProviders();
  log.debug('activeProviders', activeProviders);
};

export const getMovieInfo = (torrent: Torrent): ITorrent => {
  // log.debug('[getMovieInfo]', torrent);
  if (torrent.title.search(' ') !== -1) return parseSpacedFormat(torrent);
  else return parseDottedFormat(torrent);
};

export const downloadTorrent = async (torrent: FullTorrent) => {
  try {
    const torrentFile = await m2t.getTorrent(torrent.magnet);
    const file = torrentFile.toTorrentFile();
    writeFileSync(torrent.title, file);
    log.trace(file);
    return file;
  } catch (e) {
    log.error(e);
    return null;
  }
};

/**
 * Loops through given serach results and groups movies by title
 * @param torrents FullTorrent[] torrents search result
 * @param options params for filter
 */
export const groupTorrentsByTitle = (
  torrents: FullTorrent[],
  options = {
    maxSize: 3,
    sizeStep: 1,
  }
): ITorrent[] => {
  log.debug('[groupTorrentsByTitle] start');
  const sortedTorrents = torrents.sort(
    (a, b) => parseSize(a.size) - parseSize(b.size)
  );

  log.debug('[groupTorrentsByTitle] sortedTorrents', sortedTorrents);

  const grouppedTorrents = torrents.reduce(
    (acc: ITorrent[], cur: FullTorrent) => {
      log.trace('acc:', acc, '\ncur:', cur);
      try {
        const currentMovie = getMovieInfo(cur);
        log.info('[groupTorrentsByTitle] got parsed ITorrent', currentMovie);
        const isAdded = acc.find(
          (item: ITorrent) =>
            item.croppedTorrentTitle === currentMovie.croppedTorrentTitle
        );
        if (!isAdded) {
          log.trace('adding: ', currentMovie);
          return [...acc, currentMovie];
        }
        return acc;
      } catch (e) {
        log.error(e);
        return acc;
      }
    },
    []
  );

  log.debug('[groupTorrentsByTitle] grouppedTorrents', grouppedTorrents);

  // try to find lightest torrents first
  // if failed - increase limit by 1 GB and try again
  if (grouppedTorrents.length === 0 && torrents.length > 0)
    return groupTorrentsByTitle(torrents, {
      maxSize: options.maxSize + options.sizeStep,
      sizeStep: options.sizeStep,
    });
  log.debug('[groupTorrentsByTitle] end');
  return grouppedTorrents;
};

/**
 * searches HD movies
 * @param search query sting
 * @param category string (default Movie)
 * @param options limit and retries
 * @returns FullTorrent[]
 */

const searchFn = async (
  search: string,
  category: string = 'Movie',
  options = { limit: 20, retries: 3 }
) => {
  const { limit, retries } = options;
  if (retries <= 0) throw new Error(`No torrents found after all retries`);

  let torrents = (await TorrentSearchApi.search(
    search,
    category,
    limit
  )) as FullTorrent[];
  log.info(`Found ${torrents.length} torrents`);

  let filteredTorrents = torrents
    .filter((torrent) => !torrent.title.match(/xxx/gi))
    .filter((torrent) => torrent.title.match(/(720p)|(1080p)|(dvdrip)/gi));

  while (!filteredTorrents.length) {
    log.warn('No HD movies found');
    log.warn(`retires left: ${retries - 1}`);
    filteredTorrents = await searchFn(search, category, {
      limit: limit + 10,
      retries: retries - 1,
    });
  }
  log.info(`Got ${filteredTorrents.length} HD movies torrents`);
  log.debug(filteredTorrents);
  return filteredTorrents;
};

export const searchTorrents = async function (
  search: string,
  category: string = 'All',
  options = { limit: 20, retries: 3 }
) {
  log.debug('[Torrents] searchTorrents start', options);
  try {
    const torrents = await searchFn(search, category, options);
    const resolvedTorrents = await Promise.all(
      torrents.map(async (torrent) => {
        try {
          let file, magnet;
          try {
            file = await downloadTorrent(torrent);
          } catch (e) {
            log.error(e);
            magnet = await TorrentSearchApi.getMagnet(torrent);
          }

          log.debug('Downloaded torrent and magnet', file, magnet);
          const newTorrent = { ...torrent };
          if (file) newTorrent.torrent = file;
          if (magnet) newTorrent.magnet = magnet;
          log.info('Updated torrent:', newTorrent);
          return newTorrent;
        } catch (e) {
          log.error(e);
          return torrent;
        }
      })
    );
    log.debug('[Torrents] searchTorrents end');
    log.debug('Filtered torrents', resolvedTorrents);
    return resolvedTorrents;
  } catch (e) {
    log.error(e);
    throw new Error(`Error getting torrents`);
  }
};

interface IIndexerTorrent {
  fileName: string;
  seeders: number;
  leechers: number;
  uploaded: string;
  size: string;
  length: number;
  link?: string;
  site?: string;
  title?: string;
  score?: number;
  year?: number;
  resolution?: string;
  sourceName?: string;
  codec?: string;
  group?: string;
  season?: string;
  episode?: string;
}

const formatTitle = (title: string) => {
  if (!title) return null;
  const trimmed = title.trim().replace(/`|’/g, "'");
  if (trimmed.endsWith('.') && !trimmed.endsWith('...'))
    return trimmed.slice(0, -1);
  return trimmed;
};

export const torrentIndexerSearch = async (search: string) => {
  try {
    const torrents = (await torrentIndexer.search(
      search,
      'movie'
    )) as IIndexerTorrent[];

    log.debug('fetched torretns', torrents);
    if (torrents.length) {
      const filteredTorrents = torrents.filter(
        (torrent) =>
          !torrent.fileName.match(/xxx/gi) && parseSize(torrent.size) > 0.5
      );
      log.debug('filteredTorrents', filteredTorrents);

      const formattedTorrents = filteredTorrents.map<ITorrent>((torrent) => {
        const parcedTorrent = parseSpacedFormat({
          title: torrent.fileName,
          time: new Date().toISOString(),
          size: torrent.size,
          magnet: torrent.link || torrent.site,
          desc: '',
          provider: torrent.sourceName,
          seeds: torrent.seeders,
          peers: torrent.leechers,
        } as Torrent);
        return parcedTorrent;
      });

      log.debug('formattedTorrents', formattedTorrents);

      const unduplicated = [];
      formattedTorrents
        .sort(
          (a, b) =>
            b.torrent.seeds +
            b.torrent.peers -
            a.torrent.seeds -
            a.torrent.peers
        )
        .forEach((torrent) => {
          if (
            torrent.torrentTitle.match(/(720)|(1080)p/) &&
            !unduplicated.find(
              (el: ITorrent) => el.movieTitle === torrent.movieTitle
            )
          )
            unduplicated.push(torrent);
        });

      unduplicated.forEach(async (el) => {
        try {
          el.torrent.torrent = await torrentIndexer.torrent(el.torrent.magnet);
          log.debug('unduplicated el', el);
        } catch (e) {
          log.error(e);
        }
      });

      return unduplicated;
    }
  } catch (e) {
    log.error(e);
  }
};

interface IPage {
  url: string;
  img: string;
  title: string;
  year: string;
}

/**
 * Parses page and returns basic torrent info
 * also loads info from IMDB and saves torrent to database
 * @param page IPage
 */
const getTorrentsFromPage = async (page?: IPage): Promise<ITorrent | null> => {
  const title = page.title;

  const html = await axios(page.url);
  if (!html.data) throw new Error('Cannot fetch page ' + page.url);
  const root = parse(html.data);

  const movieInfo = root.querySelector('#movie-info');
  log.trace(movieInfo.toString());
  const hiddenSM = movieInfo.querySelector('.hidden-sm');
  const hd = hiddenSM
    .querySelectorAll('a')
    .find((el) => el.rawAttributes.title.search(/(720)|(1080)p/) !== -1);
  // nothing's found
  if (!hd) return null;

  // getMovieInfo
  const techQuality = root.querySelectorAll('.tech-quality');
  const resolution =
    hd.rawAttributes.title.search(/720p/) >= 0 ? '720p' : '1080p';

  log.trace('[getTorrentsFromPage] Found resolution', resolution);
  const infoIndex = techQuality.findIndex((el) => el.innerText.match(/720p/));
  log.trace('[getTorrentsFromPage] index', infoIndex);
  const specs = root
    .querySelectorAll('.tech-spec-info')
    [infoIndex]?.querySelectorAll('.tech-spec-element');
  const peersMatch = specs
    ?.find((el) => el.innerHTML.match(/Peers and Seeds/))
    .innerText.trim()
    .match(/(\d+ \/ \d+)/);
  const peersAndSeeds = peersMatch ? peersMatch[1].split('/') : [0, 0];
  log.info('[getTorrentsFromPage] peersAndSeeds', peersAndSeeds);

  // download torrent
  log.info('[getTorrentsFromPage] downloading torrent', hd.rawAttributes.href);
  let res;
  try {
    res = await download(hd.rawAttributes.href);
  } catch (e) {
    log.error('[getTorrentsFromPage] error', e, res);
    res = null;
  }
  log.info('[getTorrentsFromPage] got torrent', res);
  if (!res) {
    console.log('Could not download torrent');
    return null;
  }
  // === write file to filesystem (not actually needed) ===
  // writeFileSync(title + '.torrent', res.data);
  const torrent = {
    torrentTitle: `${title} (${page.year})`,
    croppedTorrentTitle: '',
    movieTitle: title,
    year: +page.year || 0,
    torrent: {
      title: `${title} (${page.year})`,
      size: specs
        ? specs.find((el) => el.innerHTML.match(/File Size/)).innerText.trim()
        : '',
      time: '',
      magnet: '',
      desc: root
        .querySelector('#movie-sub-info')
        .querySelector('.hidden-xs')
        .innerText.trim(),
      provider: 'YTS',
      seeds: +peersAndSeeds[1],
      peers: +peersAndSeeds[0],
      torrent: res || '',
    },
  } as ITorrent;
  // const mov = imdbToIMovie(await getIMDBInfo(torrent), torrent);
  // torrent.torrent.imdb = mov.id;

  log.info('[getTorrentsFromPage]found torrent', torrent);

  // saveTorrentInDB(torrent, mov);
  // insertTorrentIntoLoadedFiles(torrent);

  return torrent;
};

export const YTSsearch = async (
  query: string
): Promise<ITranslatedMovie[] | null> => {
  try {
    const res = await axios(`https://yts.mx/ajax/search`, {
      params: {
        query,
      },
    });
    log.debug(res.data);
    if (res.data.status === 'ok') {
      const pages = res.data.data as IPage[];
      log.debug('[YTSsearch] pages', pages);
      if (!pages) return;
      const torrentsPromises = pages.map((page) => getTorrentsFromPage(page));
      log.debug('[YTSsearch] torrentsPromises', torrentsPromises);
      const torrents = await Promise.all(torrentsPromises);
      log.debug('[YTSsearch] found torrents', torrents);
      const reduced = torrents.reduce((acc: ITorrent[], torrent) => {
        if (!torrent) return acc;
        return !acc.find((movie) => movie.movieTitle === torrent.movieTitle)
          ? [...acc, torrent]
          : acc;
      }, []);
      log.debug('[YTSsearch] not null torrents', reduced);
      const movies = await loadMoviesInfo(reduced);
      log.info('[YTSsearch] found movies', movies);
      return movies;
    }
    return null;
  } catch (e) {
    log.error(e);
    return null;
  }
};
