import { Container, Typography } from '@material-ui/core';
import { useEffect, useState } from 'react';
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
  route,
  match,
}: IMainPageProps & RouteComponentProps<TParams>) => {
  const dispatch = useAppDispatch();
  const { movies, search, byName, popular } = useSelector(
    (state: RootState) => state.movies
  );
  const [displayedMovies, setDisplayedMovies] = useState<ITranslatedMovie[]>(
    []
  );
  const { t } = useTranslation();

  console.log(match, route);

  // load movies on component mount
  useEffect(() => {
    const location = window.location.href;
    const filter: IFilter = { limit: LIMIT };
    switch (route) {
      // case 'search':
      //   filter.search = location.split('/').pop();
      //   break;
      case 'byname':
        filter.letter = location.split('/').pop();
        break;
      default:
        break;
    }
    console.log('[MainPage] useEffect. movies, filter', movies, filter);
    // load popular and search/byName
    dispatch(loadMovies({ filter: { limit: LIMIT + 10 } }));
    if (route === 'byname')
      dispatch(loadMovies({ filter })).then((res) => {
        console.log('[MainPage] useEffect. res, byName', res, byName);
        setDisplayedMovies(res || []);
      });
    // else if (route === 'search')
    //   dispatch(loadMovies({ filter })).then((res) =>
    //     setDisplayedMovies(res || [])
    //   );
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    let cards = null;
    if (match.path.includes('byname')) cards = byName;
    else if (match.path.includes('search')) cards = search;
    else cards = popular;

    console.log(cards);
    const displayedMovies: ITranslatedMovie[] = cards
      .map((movieId) => movies.find((movie) => movie.en.id === movieId))
      .filter(notEmpty);
    console.log(displayedMovies);
    setDisplayedMovies(displayedMovies);
  }, [byName, search, popular, match.path, movies]);

  return (
    <Container>
      <AlphabetNav />
      <FilterSortPanel />
      <CardsSlider />
      {displayedMovies.length ? (
        <Cards movies={displayedMovies} />
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
