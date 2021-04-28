import {
  Divider,
  Grid,
  makeStyles,
  Paper,
  Typography,
} from '@material-ui/core';
import PlayArrow from '@material-ui/icons/PlayArrow';
import { Link, NavLink } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { ITranslatedMovie, IUser } from '../../models/MovieInfo';
import { primaryColor } from '../../theme';
import Rating from '../Rating/Rating';

interface MovieCardMediumProps {
  card: ITranslatedMovie;
}

const useStyles = makeStyles({
  Paper: {
    display: 'flex',
    marginBottom: 20,
    padding: 10,
  },
  Img: {
    position: 'relative',
    cursor: 'pointer',
    '& img': {
      height: '15rem',
      width: '10rem',
      maxWidth: '15rem',
      objectFit: 'cover',
      marginRight: 20,
      borderRadius: 5,
    },
    '&:hover $PlayIconWrapper': {
      opacity: 1,
    },
  },
  PlayIconWrapper: {
    position: 'absolute',
    left: 'calc(50% - 3.5rem)',
    top: 'calc(50% - 3rem)',
    fontSize: '1rem',
    padding: '1rem',
    borderRadius: '50%',
    backgroundColor: '#000000aa',
    color: '#fff',
    opacity: 0,
    transition: '0.3s',
  },
  PlayIcon: {
    width: '4rem',
    height: '4rem',
  },
  Info: {
    height: '15rem',
    flexWrap: 'nowrap',
  },
  Caption: {
    fontSize: '1.5rem',
    textDecoration: 'none',
    color: 'inherit',
    marginBottom: 10,
  },
  Text: {
    fontSize: '1rem',
    fontWeight: 200,
  },
  Description: {
    maxHeight: '8rem',
    paddingTop: 10,
    fontSize: '0.9rem',
    position: 'relative',
    flexGrow: 2,
    overflow: 'hidden',
  },
  Shadow: {
    height: 40,
    width: '100%',
    position: 'absolute',
    bottom: 0,
    background: 'linear-gradient(0deg, #ffffff, #ffffff33)',
  },
  Link: {
    marginLeft: 5,
    fontSize: 'inherit',
    fontWeight: 'inherit',
    color: 'inherit',
    textTransform: 'capitalize',
    '&:hover': {
      color: primaryColor.light,
    },
  },
});

const MovieCardMedium = ({ card }: MovieCardMediumProps) => {
  const classes = useStyles();
  const { t, i18n } = useTranslation();

  const { id, img, info } = card.en;
  const title = card[i18n.language as 'en' | 'ru'].title;
  const description = card[i18n.language as 'en' | 'ru'].info.description;

  const mapItemsToLinks = (
    items: (string | IUser)[] | undefined
  ): JSX.Element[] | string => {
    if (!items || !items.length) return t('No info');
    return items.map((item: string | IUser) => {
      const text = typeof item === 'string' ? item : item.name;
      return (
        <Link to={`/genres/${text}`} key={text} className={classes.Link}>
          {t(text)}
        </Link>
      );
    });
  };

  return (
    <Paper className={classes.Paper}>
      <NavLink to={`/movies/${id}`} className={classes.Img}>
        <img src={img} alt={`${title} poster`} />
        <div className={classes.PlayIconWrapper}>
          <PlayArrow className={classes.PlayIcon} fontSize="large" />
        </div>
      </NavLink>
      <Grid container direction="column" className={classes.Info}>
        <Grid container alignItems="baseline" justify="space-between">
          <NavLink to={`/movies/${id}`} className={classes.Caption}>
            <Typography variant="caption" className={classes.Caption}>
              {title}
            </Typography>
          </NavLink>
          <Rating
            movieId={id}
            imdbRating={info.imdbRating}
            ourRating={info.rating}
          />
        </Grid>
        <Grid container className={classes.Text}>
          {t('Year')}: {info.year}
        </Grid>
        <Grid container className={classes.Text}>
          {t('Genres')}: {mapItemsToLinks(info.genres)}
        </Grid>
        <Grid container className={classes.Text}>
          {t('Length')}: {info.length ? `${info.length}${t('min')}` : 'unknown'}
        </Grid>
        <Divider />
        <Grid container className={classes.Description}>
          {description}
          <div className={classes.Shadow} />
        </Grid>
      </Grid>
    </Paper>
  );
};

export default MovieCardMedium;
