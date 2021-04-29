import axios from "axios";
import parse from "node-html-parser";
import log from "../logger/logger";
import { IMovie } from "../model/model";
import { IMDBMovie, IMDBPerson } from "./imdb";

export interface IQueryImage {
  height: number; // 4000
  imageUrl: string;
  width: number; // 2634
}

export interface IQueryVideo {
  i: IQueryImage;
  id: string; // "vi3503601689"
  l: string; // "All At Once"
  s: string; // "2:26"
}

export interface IQuery {
  i: IQueryImage;
  id: string; //"tt5073756"
  l?: string; //"All At Once"
  q?: "feature";
  rank: 49451;
  s: string; //"Jon Abrahams, Martin Abrahams"
  v: IQueryVideo[];
  vt: number; // 1
  y: number; //2016
}

const getReleaseYear = (datePublished: string): string => {
  try {
    log.debug("[getReleaseYear]", datePublished);
    return datePublished.split("-")[0];
  } catch (e) {
    log.error(e);
    return "";
  }
};

const getPlot = (description: string): string => {
  try {
    log.debug("[getPlot]", description);
    const sentenses = description.split(". ");
    return sentenses.slice(1).join(". ");
  } catch (e) {
    log.error(e);
    return null;
  }
};

const getDuration = (duration: string): string => {
  try {
    log.debug("[getDuration]", duration);
    if (!duration) return undefined;
    const regex = /PT(\d+)H(\d+)M/;
    const match = duration.match(regex);
    log.debug("[getDuration] match", match);
    if (match) {
      return +match[1] * 60 + +match[2] + "";
    }
  } catch (e) {
    log.error(e);
  }
  return null;
};

function isArray<T>(value: T[] | unknown): value is T[] {
  return !!(value as IMDBPerson[]).length;
}

const getPersons = (persons: IMDBPerson[] | IMDBPerson): IMDBPerson[] => {
  try {
    log.debug("[getPersons]", persons);
    return isArray(persons) ? persons : [persons];
  } catch (e) {
    log.error(e);
    return null;
  }
};

const getPersonsString = (persons: IMDBPerson[]): string => {
  return persons.map((dir) => dir.name).join(", ");
};

export const searchMovieID = async (query: string) => {
  try {
    const q = `${
      query.toLocaleLowerCase()[0]
    }/${query.toLocaleLowerCase().replace(/\s/gi, "_")}`;
    const res = await axios(
      `https://v2.sg.media-imdb.com/suggestion/${q}.json`,
      {
        headers: {
          Host: "v2.sg.media-imdb.com",
          Origin: "https://www.imdb.com",
          Referer: "https://www.imdb.com/interfaces/",
        },
      }
    );
    log.debug("[searchMovieID]", res.data);
    if (res.data && res.data.d) {
      return res.data.d[0]?.id;
    }
  } catch (e) {
    log.error(e);
  }
};

export const getMovieData = async (imdbId: string): Promise<IMDBMovie> => {
  try {
    const html = await axios(`https://www.imdb.com/title/${imdbId}`);
    if (html.status >= 400) throw new Error("[getMovieData] Request error");
    //<script type="application/ld+json">
    const root = parse(html.data);
    log.debug("root", root);
    const script = root.querySelector('script[type="application/ld+json"]');
    const json = JSON.parse(script.innerHTML);
    log.debug(json);
    const movie = {
      id: imdbId,
      title: json.name,
      type: json["@type"],
      year: getReleaseYear(json.datePublished),
      image: json.image,
      runtimeMins: getDuration(json.duration),
      plot: getPlot(json.description),
      genres: json.genre.join ? json.genre.join(", ") : json.genre,
      contentRating: json.contentRating,
      imDbRating: json.aggregateRating?.ratingValue,
      directors: getPersonsString(getPersons(json.director)),
      directorList: getPersons(json.director),
      stars: getPersonsString(getPersons(json.actor)),
      actorList: getPersons(json.actor),
      keywordList: json.keywords.split(","),
    };
    log.debug("[getMovieData] movie", movie);
    return movie;
  } catch (e) {
    log.error(e);
    return null;
  }
};

export const serachIMDBMovie = async (query: string): Promise<IMDBMovie> => {
  const id = await searchMovieID(query);
  const movie = getMovieData(id);
  return movie;
};
