import React from 'react';
import { Button, Grid, makeStyles, TextField } from '@material-ui/core';
import { useTranslation } from 'react-i18next';
import { movies } from '../../../axios';
import { useSelector } from 'react-redux';
import { RootState } from '../../../store/rootReducer';
import { useAppDispatch } from '../../../store/store';
import { updateComments } from '../../../store/features/MoviesSlice';

const useStyles = makeStyles({
  root: {
    margin: '1rem 0',
    minWidth: '30rem',
  },
  Form: {
    width: '100%',
  },
  margin: {
    marginBottom: 10,
  },
});

export interface IComment {
  id?: number;
  userid: string;
  movieid: string;
  text: string;
  time?: number;
}

const SendComment = () => {
  const classes = useStyles();
  const { userId } = useSelector((state: RootState) => state.user);
  const [comment, setComment] = React.useState('');
  const { t } = useTranslation();
  const dispatch = useAppDispatch();

  const sendForm = async () => {
    if (!comment.length) return;
    console.log(comment);
    const movieid = window.location.href.split('/').pop();
    const res = await movies.post('/comment', {
      userid: `${userId}`,
      movieid,
      text: comment,
      time: new Date().getTime(),
    } as IComment);
    if (res.data.status)
      dispatch(updateComments({ id: movieid, comments: [res.data.data] }));
    setComment('');
  };

  const handleInput = (e: React.KeyboardEvent<HTMLDivElement>) => {
    if (e.ctrlKey && e.key === 'Enter') sendForm();
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    sendForm();
  };

  return (
    <Grid container direction="column" className={classes.root}>
      <form onSubmit={handleSubmit} className={classes.Form}>
        <TextField
          fullWidth
          multiline
          name="comment"
          type="textarea"
          placeholder={t`Leave a comment!`}
          className={classes.margin}
          size="small"
          onChange={(e) => setComment(e.currentTarget.value)}
          onKeyPress={handleInput}
          value={comment}
        />
        <Button type="submit" className={classes.margin} variant="contained">
          {t`Send`}
        </Button>
      </form>
    </Grid>
  );
};

export default SendComment;
