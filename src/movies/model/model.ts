export interface httpResponse {
  status: boolean;
  data: any;
}

export interface IGenres {
  Action: boolean;
  Adventure: boolean;
  Animation: boolean;
  Biography: boolean;
  Comedy: boolean;
  Crime: boolean;
  Documentary: boolean;
  Drama: boolean;
  Family: boolean;
  Fantasy: boolean;
  'Film Noir': boolean;
  History: boolean;
  Horror: boolean;
  Music: boolean;
  Mystery: boolean;
  Romance: boolean;
  'Sci-Fi': boolean;
  Short: boolean;
  Sport: boolean;
  Superhero: boolean;
  Thriller: boolean;
  War: boolean;
  Western: boolean;
}
export type GenresKeys = keyof IGenres;

export interface ICountries {
  Australia: boolean;
  Canada: boolean;
  France: boolean;
  Germany: boolean;
  'Great Britain': boolean;
  India: boolean;
  Japan: boolean;
  Russia: boolean;
  Spain: boolean;
  USA: boolean;
  USSR: boolean;
  other: boolean;
}
export type CountriesKeys = keyof ICountries;

export interface IDBMovie {
  id: string;
  imdbid: string;
  title: string;
  image: string;
  year: string;
  genres: string;
  rating: string;
  votes: string;
  views: string;
  runtimemins: string;
  contentrating: string;
  countries: string;
  plot: string;
  directors: string;
  directorlist: string;
  stars: string;
  actorlist: string;
  keywordlist: string;
  images: string;
  imdbrating: string;
  maxcomments: string;
}

export interface IDBTorrent {
  movieid: string;
  torrentname: string;
  torrent: Blob;
  magnet: string;
  seeds: number;
  peers: number;
  size: number;
}

export interface ITranslatedMovie {
  en: IMovie;
  ru: IMovie;
}

export interface IKinopoiskMovie {
  kinoid: string;
  imdbid: string;
  nameen: string;
  nameru: string;
  description?: string;
  posterurl?: string;
  posterurlpreview?: string;
}

export interface IMovie {
  id: string;
  title: string;
  img: string;
  src: string;
  info: IMovieInfo;
}

export interface IMovieInfo {
  avalibility: number;
  year: number;
  genres: GenresKeys[];
  rating: number;
  imdbRating: number;
  views: number;
  length: number;
  pgRating: string;
  countries?: string[];
  comments?: IFrontComment[];
  maxComments?: number;
  description?: string;
  photos?: string[];
  videos?: string[];
  moreLikeThis?: IMovie[];
  storyline?: string;
  directors?: string;
  directorList?: IUser[];
  stars?: string;
  cast?: IUser[];
  keywords?: string[];
}

export interface IComment {
  id?: number;
  userid: string;
  movieid: string;
  text: string;
  time?: number;
}
export interface IFrontComment {
  commentid: string;
  username: string;
  avatar?: string;
  movieid: string;
  text: string;
  time: number;
}

export interface IUser {
  id: string;
  name: string;
  image?: string;
  asCharacter?: string;
  knownFor?: IMovie[];
  filmography?: IFilmography[];
  otherWorks?: string[];
  publicityListings?: string[];
  ofiicialSites?: string[];
  alternateNames?: string[];
  starSign?: string;
  born?: number;
  bornPlace?: string;
  birthName?: string;
  height?: number;
  bio?: string;
  tradeMark?: string[];
  trivia?: string[];
  personalQuotes?: string[];
}

export interface IFilmography {
  id?: string;
  job: string; //'actor', 'director'...
  movies: IMovie[];
}

interface IMDBPerson {
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
