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
  avalibility: string;
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
  avalibility?: number;
  year: number;
  genres: GenresKeys[];
  rating: number;
  imdbRating: number;
  views: number;
  length: number;
  pgRating: string;
  countries?: string[];
  comments?: IFrontComment[];
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
  id: string;
  username: string;
  avatar?: string;
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
