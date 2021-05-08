import { CountriesKeys, GenresKeys } from '../store/features/FilterSlice';

export interface ITranslatedMovie {
  ru: IMovie;
  en: IMovie;
}
export interface IMovie {
  id: string;
  title: string;
  img: string;
  src: string;
  isViewed: boolean;
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
  comments?: IComment[];
  maxComments: number;
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
  commentid: number;
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
