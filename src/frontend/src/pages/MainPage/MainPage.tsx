import { Container } from '@material-ui/core';
import { useEffect, useState } from 'react';
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

  console.log(match, route);

  // load movies on component mount
  useEffect(() => {
    const location = window.location.href;
    const filter: IFilter = { limit: LIMIT };
    switch (route) {
      case 'search':
        filter.search = location.split('/').pop();
        break;
      case 'byname':
        filter.letter = location.split('/').pop();
        break;
      default:
        break;
    }
    console.log('[MainPage] useEffect. movies, filter', movies, filter);
    // load popular and search/byName
    dispatch(loadMovies({ filter: { limit: LIMIT + 10 } }));
    if (route === 'search')
      dispatch(loadMovies({ filter })).then((res) =>
        setDisplayedMovies(res || [])
      );
    else if (route === 'byname')
      dispatch(loadMovies({ filter })).then((res) => {
        console.log('[MainPage] useEffect. res, byName', res, byName);
        setDisplayedMovies(res || []);
      });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    let cards = null;
    if (match.path.includes('byname')) cards = byName;
    else if (match.path.includes('search')) cards = search;
    else cards = popular;

    console.log(cards);
    const displayedMovies: ITranslatedMovie[] = [];
    cards.forEach((id) => {
      const movie = movies.find((movie) => movie.en.id === id);
      if (movie) displayedMovies.push(movie);
    });
    console.log(displayedMovies);
    setDisplayedMovies(displayedMovies);
  }, [byName, search, popular, match.path, movies]);

  return (
    <Container>
      <AlphabetNav />
      <FilterSortPanel />
      <CardsSlider />
      <Cards movies={displayedMovies} />
    </Container>
  );
};

export default MainPage;
