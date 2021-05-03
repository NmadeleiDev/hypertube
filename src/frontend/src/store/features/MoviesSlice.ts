import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import axios from 'axios';
import { movies, search as axiosSearch } from '../../axios';
import { IComment, ITranslatedMovie } from '../../models/MovieInfo';
import { AppDispatch } from '../store';
import { CountriesKeys, GenresKeys } from './FilterSlice';

export interface MoviesItems {
  search?: ITranslatedMovie[];
  popular?: ITranslatedMovie[];
  byName?: ITranslatedMovie[];
}

export type MovieItemsKeys = keyof MoviesItems;
export interface CommentsItems {
  comments: IComment[];
  id?: string;
}
export interface ErrorItem {
  error: string;
}

export interface MoviesState {
  loading: boolean;
  error: string | null;
  movies: ITranslatedMovie[];
  search: string[];
  popular: string[];
  byName: string[];
  isEndOfMovies: boolean;
}

export interface IMoviesResponse {
  data: ITranslatedMovie[] | null;
  error?: string;
  status: boolean;
}

export interface IFilter {
  limit?: number;
  offset?: number;
  letter?: string;
  search?: string;
  genres?: GenresKeys[];
  years?: [] | number[][];
  countries?: CountriesKeys[];
  id?: string | number;
}
interface ILoaderInfo {
  activePeers: number;
  loadedPercent: number;
}

const initialState = {
  loading: false,
  error: null,
  movies: [],
  search: [],
  popular: [],
  byName: [],
  isEndOfMovies: false,
} as MoviesState;

const MoviesSlice = createSlice({
  name: 'movies',
  initialState,
  reducers: {
    addMovies(state, { payload }: PayloadAction<MoviesItems>) {
      const add = (type: MovieItemsKeys) => {
        console.log(
          '[loadMovies] payload, moviesLength',
          payload,
          state.movies.length
        );
        if (!payload[type] || !payload[type]?.length) return;
        if (state.movies.length) {
          const newMovies = payload[type]?.filter(
            (movie) => !state.movies.find((el) => el.en.id === movie.en.id)
          );
          console.log('[loadMovies] newMovies', newMovies);
          if (newMovies?.length) {
            state.movies = state.movies.concat(newMovies);
          }
          // set unique movie ids
          const set = new Set(
            state[type].concat((payload[type] || []).map((el) => el.en.id))
          );
          state[type] = Array.from(set);
        } else {
          const arr = payload[type] || [];
          state.movies = arr;
          state[type] = arr.map((el) => el.en.id);
        }
      };

      state.loading = false;
      if (payload.search) {
        add('search');
      } else if (payload.byName) {
        add('byName');
      } else if (payload.popular) {
        add('popular');
      } else {
        throw new Error('[movies:addMovies] no movies in payload');
      }
    },
    setMovies(state, { payload }: PayloadAction<MoviesItems>) {
      const set = (type: MovieItemsKeys) => {
        if (!payload[type]) return;
        // state[type] = payload[type] || [];
        const newMovies = payload[type]?.filter(
          (movie) => !state.movies.find((el) => el.en.id === movie.en.id)
        );
        if (newMovies) {
          state.movies = state.movies.concat(newMovies);
          state[type] = payload[type]?.map((el) => el.en.id) || [];
        }
      };
      state.loading = false;
      if (payload.search) {
        set('search');
      } else if (payload.byName) {
        set('byName');
      } else if (payload.popular) {
        set('popular');
      } else {
        throw new Error('[movies:setMovies] no movies in payload');
      }
    },
    updateComments(state, { payload }: PayloadAction<CommentsItems>) {
      state.loading = false;
      if (!payload.id)
        throw new Error('[movies:updateComments] no id in payload');
      if (!payload.comments)
        throw new Error('[movies:updateComments] no comments in payload');
      const movie = state.movies.find((movie) => movie.en.id === payload.id);
      if (!movie) throw new Error('[movies:updateComments] no movie found');
      const newComments = payload.comments.filter((comment) =>
        movie.en.info.comments
          ? !movie.en.info.comments.find(
              (el) => el.commentid === comment.commentid
            )
          : !!comment
      );
      movie.en.info.comments = movie.en.info.comments
        ? [...movie.en.info.comments, ...newComments]
        : [...newComments];
    },
    setEndOfMovies(state) {
      state.isEndOfMovies = true;
    },
    resetEndOfMovies(state) {
      state.isEndOfMovies = false;
    },
    setError(state, { payload }: PayloadAction<ErrorItem>) {
      state.error = payload.error;
      state.loading = false;
    },
    resetError(state) {
      state.error = null;
      state.loading = false;
    },
    startLoading(state) {
      state.loading = true;
    },
    stopLoading(state) {
      state.loading = false;
    },
  },
});

const loadMoviesAsync = async (
  params?: IFilter
): Promise<IMoviesResponse | null> => {
  try {
    console.log('[loadMoviesAsync] params', params);
    let res = null;
    if (params?.search)
      res = await axiosSearch('find', {
        params: {
          category: 'All',
          search: params?.search,
        },
      });
    else if (params?.letter)
      res = await movies(`byname`, {
        params,
      });
    else
      res = await movies('movies', {
        params,
      });

    // console.log(res.data);
    return res.data;
  } catch (e) {
    console.log(e);
    return null;
  }
};

export const loadMovies = ({
  filter = {
    limit: 20,
    offset: 0,
  },
  callback,
}: {
  filter?: IFilter;
  callback?: (arg0: boolean) => void;
}) => async (dispatch: AppDispatch) => {
  dispatch(startLoading());
  const limit = filter.limit || 20;
  console.log('[loadMovies] filter', filter);
  const movies = await loadMoviesAsync(filter);
  console.log('[loadMovies] res', movies);

  if (movies?.status && movies.data) {
    if (movies.data.length < limit) dispatch(setEndOfMovies());
    const removedNulls = movies.data?.map((el) =>
      !el.ru ? { en: el.en, ru: el.en } : el
    );
    console.log('[loadMovies] removedNulls:', removedNulls);
    if (filter.letter) dispatch(addMovies({ byName: removedNulls }));
    else if (filter.search)
      filter.offset && filter.offset > 0
        ? dispatch(addMovies({ search: removedNulls }))
        : dispatch(setMovies({ search: removedNulls }));
    else dispatch(addMovies({ popular: removedNulls }));
    if (callback)
      callback(movies && movies.data ? movies.data.length < limit : true);
    return removedNulls;
  } else {
    dispatch(stopLoading());
    const error = {
      en: 'Loading error',
      ru: 'Не удалось загрузить фильмы',
    };
    dispatch(setError({ error: movies?.error || 'Loading error' }));
    return null;
  }
};

export const loadActivePeers = async (
  id: string
): Promise<ILoaderInfo | null> => {
  let data = null;
  try {
    const res = await axios(`/api/loader/stats/${id}`);
    if (res.status < 400 && res.data.status) {
      data = res.data.data as ILoaderInfo;
    }
  } catch (e) {
    console.log(e);
  }
  return data;
};

export const loadMovie = (id: string) => async (dispatch: AppDispatch) => {
  dispatch(startLoading());
  const movie = await loadMoviesAsync({ id });
  console.log('loadMovie', movie);
  if (movie?.status && movie.data) dispatch(addMovies({ search: movie.data }));
  else {
    dispatch(setError({ error: 'Loading error' }));
  }
};

/**
 * loads comments
 * @param params object with values:
 *    @param movieId movieId to find it in redux
 *    @param limit number of comments to fetch, default 5
 *    @param offset comments offset, default 0
 * @param callback fires when end of comments is reached
 */
export const loadComments = (
  {
    movieId,
    limit = 5,
    offset = 0,
  }: { movieId: string; limit?: number; offset?: number },
  callback?: (length: number) => void
) => async (dispatch: AppDispatch) => {
  dispatch(startLoading());
  const res = await movies('comments', {
    params: {
      movieId,
      limit,
      offset,
    },
  });
  console.log('loadComments', res.data);
  if (res.data.status) {
    const comments = res.data.data as IComment[];
    comments.length && dispatch(updateComments({ comments, id: movieId }));
    if (callback) callback(comments.length);
  }
};

export const {
  addMovies,
  setMovies,
  updateComments,
  setEndOfMovies,
  resetEndOfMovies,
  setError,
  resetError,
  startLoading,
  stopLoading,
} = MoviesSlice.actions;

export default MoviesSlice.reducer;
