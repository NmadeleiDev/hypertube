import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { RootState } from '../rootReducer';
import { AppDispatch } from '../store';
import { IFilter, loadMovies } from './MoviesSlice';

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
export interface IYears {
  'before 1980': boolean;
  '1980-1990': boolean;
  '1990-2000': boolean;
  '2000-2010': boolean;
  '2010-2020': boolean;
}
export type YearsKeys = keyof IYears;
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

export const items = {
  genres: [
    'Action',
    'Adventure',
    'Animation',
    'Biography',
    'Comedy',
    'Crime',
    'Documentary',
    'Drama',
    'Family',
    'Fantasy',
    'Film Noir',
    'History',
    'Horror',
    'Music',
    'Mystery',
    'Romance',
    'Sci-Fi',
    'Short',
    'Sport',
    'Superhero',
    'Thriller',
    'War',
    'Western',
  ] as GenresKeys[],
  years: [
    'before 1980',
    '1980-1990',
    '1990-2000',
    '2000-2010',
    '2010-2020',
  ] as YearsKeys[],
  countries: [
    'Australia',
    'Canada',
    'France',
    'Germany',
    'Great Britain',
    'India',
    'Japan',
    'Russia',
    'Spain',
    'USA',
    'USSR',
    'other',
  ] as CountriesKeys[],
  // additional: ['subtitles', 'HD', 'Multi-audio'],
};

export interface GenresItem {
  category: 'genres';
  name: GenresKeys;
  value: boolean;
}

export interface YearsItem {
  category: 'years';
  name: YearsKeys;
  value: boolean;
}

export interface CountriesItem {
  category: 'countries';
  name: CountriesKeys;
  value: boolean;
}

export interface SetFilterPayload {
  filter: GenresItem | YearsItem | CountriesItem;
}

export interface FilterState {
  genres: IGenres;
  years: IYears;
  countries: ICountries;
}
export type FilterStateKeys = keyof FilterState;

const initialState = {
  genres: {
    Action: false,
    Adventure: false,
    Animation: false,
    Biography: false,
    Comedy: false,
    Crime: false,
    Documentary: false,
    Drama: false,
    Family: false,
    Fantasy: false,
    'Film Noir': false,
    History: false,
    Horror: false,
    Music: false,
    Mystery: false,
    Romance: false,
    'Sci-Fi': false,
    Short: false,
    Sport: false,
    Superhero: false,
    Thriller: false,
    War: false,
    Western: false,
  },
  years: {
    'before 1980': false,
    '1980-1990': false,
    '1990-2000': false,
    '2000-2010': false,
    '2010-2020': false,
  },
  countries: {
    Australia: false,
    Canada: false,
    France: false,
    Germany: false,
    'Great Britain': false,
    India: false,
    Japan: false,
    Russia: false,
    Spain: false,
    USA: false,
    USSR: false,
    other: false,
  },
  // additional: { subtitles: false, HD: false, 'Multi-audio': false },
};

const filterSlice = createSlice({
  name: 'filter',
  initialState,
  reducers: {
    setFilterState(state, { payload }: PayloadAction<SetFilterPayload>) {
      const { category, name, value } = payload.filter;
      return {
        ...state,
        [category]: {
          ...state[category],
          [name]: value,
        },
      };
    },
    resetFilterState() {
      return initialState;
    },
  },
});
// category: 'genres' | 'years' | 'countries',

const isCountriesKey = (value: string): value is CountriesKeys => {
  const allowedKeys: string[] = items.countries.filter((el) => el !== 'other');
  return allowedKeys.indexOf(value) !== -1;
};
const isGenresKey = (value: string): value is GenresKeys => {
  const allowedKeys: string[] = items.genres;
  return allowedKeys.indexOf(value) !== -1;
};
const isYearsKey = (value: string): value is YearsKeys => {
  const allowedKeys: string[] = items.years;
  return allowedKeys.indexOf(value) !== -1;
};

export const isCountiesStateEmpty = (state: ICountries) => {
  for (let country in state) {
    if (isCountriesKey(country) && state[country]) return false;
  }
  return true;
};
export const isYearsStateEmpty = (state: IYears) => {
  for (let year in state) {
    if (isYearsKey(year) && state[year]) return false;
  }
  return true;
};
export const isGenresStateEmpty = (state: IGenres) => {
  for (let genre in state) {
    if (isGenresKey(genre) && state[genre]) return false;
  }
  return true;
};

const getYearsFilter = (years: IYears): number[][] | [] => {
  const yearsArray = items.years.filter(
    (value) => years[value] && [value, years[value]]
  );

  const result: number[][] = [];
  yearsArray.forEach((year) => {
    if (year.startsWith('before')) {
      result.push([0, 1980]);
    } else if (
      result.length > 0 &&
      year.match(`${result[result.length - 1][1]}`)
    ) {
      // if last value in result contains one of current years
      // (actually, if it's true - should find the first one)
      // then we replace that year with the second one
      const arr = year.split('-');
      result[result.length - 1][1] = +arr[1];
    } else {
      const arr = year.split('-');
      result.push([+arr[0], +arr[1]]);
    }
  });
  console.log(result);
  return result;
};

/**
 * Requests movies, based on filterState.
 * Empty or object with fields required.
 * @param param0 object with _limit and _start fields
 */
export const reqeuestMoviesWithFilters = (
  filter: IFilter = {
    limit: 20,
    offset: 0,
  }
) => async (dispatch: AppDispatch, getState: () => RootState) => {
  const filterState = getState().filter;
  const genres = items.genres.filter(
    (value) => filterState.genres[value] && [value, filterState.genres[value]]
  );
  const years = getYearsFilter(filterState.years);
  const countries = items.countries.filter(
    (value) =>
      filterState.countries[value] && [value, filterState.countries[value]]
  );
  const filterParams = { ...filter, genres, years, countries };

  console.log(filterParams);
  dispatch(loadMovies({ filter: filterParams }));
};

export const yearInRange = (year: number, years: IYears) => {
  for (let entry of Object.entries(years)) {
    if (!entry[1]) continue;
    const rangeNumbers = entry[0].startsWith('before')
      ? [0, +entry[0].split(' ')[1]]
      : entry[0].split('-').map((x) => +x);
    if (year >= rangeNumbers[0] && year < rangeNumbers[1]) return true;
  }
  return false;
};
export const genreInGenres = (
  movieGenres: GenresKeys[],
  filterGenres: IGenres
) => {
  for (let genre of movieGenres) {
    if (filterGenres[genre]) return true;
  }
  return false;
};
export const countryInCountries = (
  movieCountries: string[],
  filterCountries: ICountries
) => {
  for (let country of movieCountries) {
    return isCountriesKey(country)
      ? filterCountries[country]
      : filterCountries.other;
  }
  return false;
};

export const { setFilterState, resetFilterState } = filterSlice.actions;

export default filterSlice.reducer;
