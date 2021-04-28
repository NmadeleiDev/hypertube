import React from 'react';
import { makeStyles, Typography } from '@material-ui/core';

interface Props {
  text: string;
}

const useStyles = makeStyles({
  h3: {
    fontSize: '1.6rem',
    textTransform: 'uppercase',
    margin: '10px 0',
    fontWeight: 900,
  },
});

const CategoryHeader = ({ text }: Props) => {
  const classes = useStyles();

  return (
    <Typography variant="h3" className={classes.h3}>
      {text}
    </Typography>
  );
};

export default CategoryHeader;
