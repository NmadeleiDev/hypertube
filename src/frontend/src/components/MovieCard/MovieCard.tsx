import React from 'react';
import { Card, Grid, makeStyles, Typography } from '@material-ui/core';
import cn from 'classnames';
import MovieCardMedium from '../MovieCardMedium/MovieCardMedium';
import { useTranslation } from 'react-i18next';
import { NavLink } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/rootReducer';

export interface MovieCardProps {
  display?: 'image' | 'lines' | 'grid';
  className?: string;
  id: string;
}

export const useStyles = makeStyles({
  root: {
    // cursor: 'pointer',
  },
  Grid: {
    maxWidth: 150,
    marginRight: '1rem',
    marginBottom: '1rem',
  },
  Media: {
    height: 220,
    width: 150,
  },
  Name: {
    fontSize: '1rem',
    fontWeight: 800,
  },
});

const MovieCard: React.FC<MovieCardProps> = ({
  display = 'image',
  className,
  id: card,
}: MovieCardProps) => {
  const classes = useStyles();
  const { i18n } = useTranslation();
  const { movies } = useSelector((state: RootState) => state.movies);
  const movie = movies.find((el) => el.en.id === card);

  const getCard = () => {
    if (!movie) return null;
    if (display === 'lines') return <MovieCardMedium card={movie} />;
    const bg = movie.en.img || movie.ru.img;
    return (
      <NavLink to={`/movies/${movie.en.id}`} className={classes.root}>
        <Card style={{ height: 'fit-content' }}>
          <div
            className={classes.Media}
            style={{
              background: `url(${bg})`,
              backgroundSize: 'cover',
            }}
          ></div>
        </Card>
        {display !== 'image' && (
          <Grid container direction="column" alignItems="center">
            {display === 'grid' && (
              <Typography className={classes.Name}>
                {movie[i18n.language as 'en' | 'ru'].title}
              </Typography>
            )}
          </Grid>
        )}
      </NavLink>
    );
  };

  const getClassName = () =>
    display === 'lines'
      ? cn(classes.root, className)
      : cn(classes.root, className, classes.Grid);

  return (
    <Grid
      item
      container
      wrap="wrap"
      direction={display === 'lines' ? 'column' : 'row'}
      className={getClassName()}
    >
      {getCard()}
    </Grid>
  );
};

export default MovieCard;
