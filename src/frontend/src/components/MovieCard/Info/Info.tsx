import React from 'react';
import { Grid, makeStyles, Typography } from '@material-ui/core';
import Views from '../Views/Views';
import Rating from '@material-ui/lab/Rating';
import { StarOutlineRounded } from '@material-ui/icons';
import { IMovieInfo } from '../../../models/MovieInfo';

export interface InfoProps {
  info: IMovieInfo;
  className?: string;
}

export const useStyles = makeStyles((theme) => ({
  Info: {
    maxWidth: 200,
    minWidth: 200,
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'flex-start',
  },
  InfoHorizontal: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
  },
  SecondaryInfo: {
    alignItems: 'center',
  },
  SecondaryText: {
    fontSize: '0.8rem',
    fontWeight: 400,
    color: theme.palette.grey[700],
  },
  Genre: {
    marginLeft: '0.3rem',
    display: 'flex',
    '& $SecondaryText': {
      textTransform: 'capitalize',
    },
    '& $SecondaryText:nth-child(n+2)::before': {
      content: '"·"',
      margin: '0 0.1rem',
    },
  },
}));

const Info = ({
  info: { year, genres, rating, views, length, pgRating },
}: InfoProps) => {
  const classes = useStyles();
  return (
    <Grid container className={classes.InfoHorizontal}>
      <Grid item>
        <Grid container className={classes.SecondaryInfo}>
          <Typography className={classes.SecondaryText}>{year}</Typography>
          <Grid className={classes.Genre}>
            {genres.map((genre) => (
              <span key={genre} className={classes.SecondaryText}>
                {genre}
              </span>
            ))}
          </Grid>
        </Grid>
        <>
          <Rating
            name="customized-empty"
            defaultValue={rating}
            precision={0.5}
            emptyIcon={
              <StarOutlineRounded color="secondary" fontSize="inherit" />
            }
          />
          <Views views={views} />
        </>
      </Grid>
      {orientation === 'horizontal' && (
        <Grid item className={classes.InfoHorizontal}>
          <Typography className={classes.SecondaryText}>{length}min</Typography>
          <Typography className={classes.SecondaryText}>{pgRating}</Typography>
        </Grid>
      )}
    </Grid>
  );
};

export default Info;
