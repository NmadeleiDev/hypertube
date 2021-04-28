import { Divider, Grid, makeStyles, Typography } from '@material-ui/core';
import React, { useEffect } from 'react';
import { Link, RouteComponentProps, useHistory } from 'react-router-dom';
import { IUser } from '../../models/MovieInfo';
import HorizontalGrid from '../HorizontalGrid/HorizontalGrid';
import CategoryHeader from '../CategoryHeader/CategoryHeader';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/rootReducer';
import Comments from '../Comments/Comments';
import { useAppDispatch } from '../../store/store';
import { loadMovie, resetError } from '../../store/features/MoviesSlice';
import { useTranslation } from 'react-i18next';
import { primaryColor } from '../../theme';
import Player from '../Player/Player';
import { useToast } from '../../hooks/useToast';
import NativePlayer from '../Player/NativePlayer';

interface TParams {
  id: string;
}

const useStyles = makeStyles({
  root: {
    fontSize: '1.3rem',
    padding: '10px 0',
  },
  Header: {
    fontSize: '2rem',
    margin: '0 10px',
    fontWeight: 700,
  },
  Divider: {
    margin: '10px 0',
  },
  Poster: {
    height: '23rem',
    marginRight: 10,
  },
  MainInfoText: {
    margin: 5,
    width: '100%',
    display: 'flex',
    alignItems: 'flex-end',
    fontSize: 'inherit',
    fontWeight: 200,
  },
  Link: {
    marginLeft: 5,
    fontSize: 'inherit',
    color: 'inherit',
    fontWeight: 'inherit',
    textTransform: 'capitalize',
    '&:hover': {
      color: primaryColor.light,
    },
  },
  Description: {
    marginTop: 10,
    fontSize: 'inherit',
  },
  Video: {
    margin: '15px 0',
  },
  AdditionalInfo: {
    margin: '10px 0',
  },
});

const MovieFullInfo = ({ match }: RouteComponentProps<TParams>) => {
  const classes = useStyles();
  const dispatch = useAppDispatch();
  const { t, i18n } = useTranslation();
  const history = useHistory();
  const { toast } = useToast();
  const { isAuth } = useSelector((state: RootState) => state.user);
  const { movies, error } = useSelector((state: RootState) => state.movies);
  const movie = movies.find((movie) => movie.en.id === match.params.id);
  const headerRef = React.useRef<HTMLHeadingElement | null>(null);

  // if no movies in redux - load some, we've landed on movie's page
  useEffect(() => {
    if (!movie) dispatch(loadMovie(match.params.id));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (error) {
      history.push('/');
      toast({ text: error[i18n.language as 'en' | 'ru'] }, 'error');
      dispatch(resetError());
    }
  }, [error, dispatch, history, t, toast, i18n.language]);

  useEffect(() => {
    if (!headerRef || !headerRef.current) return;
    headerRef.current.scrollIntoView(true);
  }, []);

  const mapItemsToLinks = (
    items: (string | IUser)[] | undefined,
    backup?: string
  ): JSX.Element[] | string => {
    if (!items || !items.length) return backup || t('No info');
    return items.map((item: string | IUser) => {
      const text = typeof item === 'string' ? item : item.name;
      return (
        <Link to={`/genres/${text}`} key={text} className={classes.Link}>
          {t(text)}
        </Link>
      );
    });
  };
  const mapItems = (items?: (string | IUser)[], backup?: string): string => {
    if (!items || !items.length) return backup || t('No info');
    const names = items.map((item: string | IUser) => {
      const text = typeof item === 'string' ? item : item.name;
      return text;
    });
    return names.join(', ');
  };

  // if (!isAuth)
  //   return (
  //     <Grid
  //       container
  //       direction="column"
  //       alignItems="center"
  //       className={classes.root}
  //     >
  //       <Typography className={classes.Description}>{t`AuthOnly`}</Typography>
  //     </Grid>
  //   );

  if (!movie) return null;
  return (
    <Grid container direction="column" className={classes.root}>
      <Typography variant="h2" className={classes.Header} ref={headerRef}>
        {movie[i18n.language as 'en' | 'ru'].title}
      </Typography>
      <Divider className={classes.Divider} />
      <Grid container wrap="nowrap">
        <img
          className={classes.Poster}
          src={movie.en.img}
          alt={`${movie.en.title} poster`}
        />
        <Grid item container direction="column">
          <Grid container className={classes.MainInfoText}>
            {t`Year`}: {movie.en.info.year}
          </Grid>
          <Grid container className={classes.MainInfoText}>
            {t`Genres`}: {mapItemsToLinks(movie.en.info.genres)}
          </Grid>
          <Grid container className={classes.MainInfoText}>
            {t`Length`}:{' '}
            {movie.en.info.length > 0
              ? `${movie.en.info.length}${t('min')}`
              : t`unknown`}
          </Grid>
          <Grid container className={classes.MainInfoText}>
            {t`Views`}: {movie.en.info.views}
          </Grid>
          <Grid container className={classes.MainInfoText}>
            {t`PG rating`}: {movie.en.info.pgRating}
          </Grid>
          <Grid container className={classes.MainInfoText}>
            {t`Directed`}:{' '}
            {mapItems(movie.en.info.directorList, movie.en.info.directors)}
          </Grid>
          <Grid container className={classes.MainInfoText}>
            {t`Actors`}: {mapItems(movie.en.info.cast, movie.en.info.stars)}
          </Grid>
        </Grid>
      </Grid>
      <Grid container className={classes.Video}>
        <NativePlayer id={movie.en.id} />
        <NativePlayer src={`/api/test/${movie.en.id}`} />

        <Player
          src={`/api/test/${movie.en.id}`}
          // id={movie.en.id}
          title={movie[i18n.language as 'en' | 'ru'].title}
        />
      </Grid>
      <Grid container direction="column" className={classes.AdditionalInfo}>
        <CategoryHeader text={t`About movie`} />
        <Typography variant="body1" className={classes.Description}>
          {movie[i18n.language as 'en' | 'ru'].info.description || t`No info`}
        </Typography>
        <Divider className={classes.Divider} />
        {movie.en.info.photos && (
          <HorizontalGrid
            sources={movie.en.info.photos}
            name={movie.en.title}
            type={'photo'}
          />
        )}
        {movie.en.info.videos && (
          <HorizontalGrid
            sources={movie.en.info.videos}
            name={movie.en.title}
            type={'video'}
          />
        )}
        <Comments movie={movie.en} />
      </Grid>
    </Grid>
  );
};

export default MovieFullInfo;
