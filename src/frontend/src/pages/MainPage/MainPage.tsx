import { Container, Typography } from '@material-ui/core';
import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { RouteComponentProps } from 'react-router-dom';
import { LIMIT } from '../..';
import AlphabetNav from '../../components/AlphabetNav/AlphabetNav';
import Cards from '../../components/Cards/Cards';
import CardsSlider from '../../components/CardsSlider/CardsSlider';
import FilterSortPanel from '../../components/FilterSortPanel/FilterSortPanel';
import { ITranslatedMovie } from '../../models/MovieInfo';
import { IFilter, loadMovies } from '../../store/features/MoviesSlice';
import { RootState } from '../../store/rootReducer';
import { useAppDispatch } from '../../store/store';
import { notEmpty } from '../../utils';

interface IMainPageProps {
  route?: string;
}

interface TParams {
  id?: string;
}

const MainPage = ({
  match,
  route,
}: IMainPageProps & RouteComponentProps<TParams>) => {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();
  const cardsRef = useRef<HTMLDivElement>(null);
  const { movies, search, byName, popular } = useSelector(
    (state: RootState) => state.movies
  );
  const [displayedMovies, setDisplayedMovies] = useState<ITranslatedMovie[]>(
    []
  );

  // load movies on component mount
  useEffect(() => {
    const m = match.url.match(/genres|byname/);
    const route = m ? m[0] : '';
    console.log('route', route);
    const filter: IFilter = { limit: LIMIT };
    switch (route) {
      case 'genres':
        filter.genre = match.params.id;
        break;
      case 'byname':
        filter.letter = match.params.id;
        break;
      default:
        break;
    }
    console.log(
      '[MainPage] useEffect. movies, filter, route',
      movies,
      filter,
      route
    );
    if (route === 'byname' || route === 'genres') {
      dispatch(loadMovies({ filter }));
      // .then((res) => {
      //   console.log('[MainPage] useEffect. res, byName', res, byName);
      //   setDisplayedMovies(res || []);
      // });
    } else {
      // load popular and search/byName
      dispatch(loadMovies({ filter }));
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [match.url]);

  useEffect(() => {
    let cards = null;
    if (match.url.includes('byname')) cards = byName;
    else if (match.url.match(/search|genres/)) cards = search;
    else cards = popular;

    console.log(cards);
    const displayedMovies: ITranslatedMovie[] = cards
      .map((movieId) => movies.find((movie) => movie.en.id === movieId))
      .filter(notEmpty);
    console.log(displayedMovies);
    setDisplayedMovies(displayedMovies);
    cardsRef.current?.scrollTo();
  }, [byName, search, popular, match.url, movies]);

  return (
    <Container>
      <AlphabetNav />
      <FilterSortPanel />
      <CardsSlider />
      {displayedMovies.length ? (
        <Cards ref={cardsRef} movies={displayedMovies} />
      ) : (
        <Typography
          variant="h5"
          style={{ width: '100%', textAlign: 'center' }}
        >{t`NoMoviesSearch`}</Typography>
      )}
    </Container>
  );
};

export default MainPage;
