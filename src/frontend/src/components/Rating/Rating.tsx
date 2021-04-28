import React, { useState } from 'react';
import { movies } from '../../axios';
import { getToken } from '../../store/features/UserSlice';
import MUIRating from '@material-ui/lab/Rating';
import StarBorderRounded from '@material-ui/icons/StarBorderRounded';
import { Grid, makeStyles, Typography } from '@material-ui/core';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/rootReducer';

interface Props {
  movieId: string;
  imdbRating: number;
  ourRating: number;
}

const useStyles = makeStyles({
  Typography: {
    padding: 5,
    marginRight: 5,
    border: '1px solid #ddd',
    borderRadius: 5,
  },
});

const Rating = ({ movieId, imdbRating, ourRating }: Props) => {
  const classes = useStyles();
  const [rating, setRating] = useState(ourRating);
  const { isAuth } = useSelector((state: RootState) => state.user);

  const handleRatingChange = async (
    e: React.ChangeEvent<{}>,
    newRating: number | null
  ) => {
    e.stopPropagation();
    if (newRating !== null) {
      console.log(newRating);
      const res = await movies.patch(
        '/rating',
        { movieId, rating: newRating },
        {
          headers: {
            accessToken: getToken(),
          },
        }
      );
      if (res.data.status) setRating(+res.data.data);
    }
  };

  return (
    <Grid container alignItems="center" justify="flex-end">
      <Typography className={classes.Typography}>IMDB</Typography>
      <Typography>{imdbRating}</Typography>
      <MUIRating
        disabled={!isAuth}
        value={rating || imdbRating / 2}
        name={`${movieId}-rating`}
        emptyIcon={<StarBorderRounded color="disabled" fontSize="inherit" />}
        onChange={handleRatingChange}
      />
    </Grid>
  );
};

export default Rating;
