const format = require('string-format');
const axios = require('axios');
const TorrentProvider = require('torrent-search-api/lib/TorrentProvider');
const { writeFile } = require('fs');

class _1337xx extends TorrentProvider {
  constructor() {
    super({
      name: '1337xx',
      baseUrl: 'http://www.1337x.to',
      searchUrl: '/category-search/{query}/{cat}/1/',
      categories: {
        All: 'url:/search/{query}/1/',
        Movies: 'Movies',
        TV: 'TV',
        Games: 'Games',
        Music: 'Music',
        Anime: 'Anime',
        Applications: 'Apps',
        Documentaries: 'Documentaries',
        Xxx: 'XXX',
        Other: 'Other',
        Top100: 'url:/top-100',
      },
      defaultCategory: 'All',
      resultsPerPageCount: 20,
      itemsSelector: 'tbody > tr',
      itemSelectors: {
        title: 'a:nth-child(2)',
        time: '.coll-date',
        seeds: '.seeds | int',
        peers: '.leeches | int',
        size: '.size@html | until:<sp',
        desc: 'a:nth-child(2)@href',
      },
      paginateSelector: '.pagination li:nth-last-child(2) a@href',
      torrentDetailsSelector: '.torrent-detail-page@html',
      enableCloudFareBypass: true,
    });
  }

  async downloadTorrent(torrent, path) {
    const hash = await this.xray(torrent.desc, '.infohash-box span');
    const link = await this.xray(
      format('https://torrage.info/torrent.php?h=%{0}', hash),
      '.btn-success > a@href'
    );
    const res = await axios(
      format('https://torrage.info/torrent.php?h=%{0}', hash),
      { withCredentials: true }
    );
    console.log(res);
    const torrentFile = await axios(link, { withCredentials: true });
    console.log(torrentFile);
    return path ? writeFile(path, torrentFile.data) : torrentFile.data;
  }
}

module.exports = _1337xx;
