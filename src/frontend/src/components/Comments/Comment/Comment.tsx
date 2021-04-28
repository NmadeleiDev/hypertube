import React from 'react';
import { Avatar, Grid, makeStyles, Typography } from '@material-ui/core';
import { IComment } from '../../../models/MovieInfo';
import { primaryColor } from '../../../theme';

const useStyles = makeStyles({
  root: {
    width: '100%',
    padding: 10,
    marginBottom: 5,
  },
  h4: {
    fontWeight: 700,
    fontSize: '1rem',
    marginBottom: 10,
  },
  Image: {
    marginRight: 10,
    background: primaryColor.main,
  },
  Text: {
    width: '100%',
    textAlign: 'justify',
  },
});

const Comment: React.FC<IComment> = ({ username, text, avatar }) => {
  const classes = useStyles();

  return (
    <Grid container wrap="nowrap" className={classes.root}>
      <Avatar
        src={avatar}
        alt={`${username}'s avatar`}
        className={classes.Image}
      >
        {username.substr(0, 2)}
      </Avatar>
      <Grid>
        <Typography variant="h4" className={classes.h4}>
          {username}
        </Typography>
        <Typography variant="body2" className={classes.Text}>
          {text}
        </Typography>
      </Grid>
    </Grid>
  );
};

export default Comment;
