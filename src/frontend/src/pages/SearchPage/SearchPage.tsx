import { Container } from '@material-ui/core';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { useHistory } from 'react-router';
import Cards from '../../components/Cards/Cards';
import FilterSortPanel from '../../components/FilterSortPanel/FilterSortPanel';
import CardLoader from '../../components/MovieCard/CardLoader/CardLoader';
import { useToast } from '../../hooks/useToast';
import { ITranslatedMovie } from '../../models/MovieInfo';
import { RootState } from '../../store/rootReducer';
import { notEmpty } from '../../utils';

const SearchPage = () => {
  const { movies, search, loading, error } = useSelector(
    (state: RootState) => state.movies
  );
  const history = useHistory();
  const { toast } = useToast();
  const { t } = useTranslation();
  if (!loading && movies.length === 0) history.push('/');
  if (error) {
    history.push('/');
    toast({ text: t(error) }, 'error');
  }

  const displayedMovies: ITranslatedMovie[] = search
    .map((movieId) => movies.find((movie) => movie.en.id === movieId))
    .filter(notEmpty);

  const cards = loading ? (
    <CardLoader display="lines" />
  ) : (
    <Cards movies={displayedMovies} />
  );

  return (
    <Container>
      <FilterSortPanel />
      {cards}
    </Container>
  );
};

export default SearchPage;
