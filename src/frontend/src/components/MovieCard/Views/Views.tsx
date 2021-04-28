import { Grid, makeStyles, Typography } from '@material-ui/core';
import { VisibilityRounded } from '@material-ui/icons';
import React from 'react';

interface Props {
  views: number;
}

export const useStyles = makeStyles((theme) => ({
  SecondaryText: {
    fontSize: '0.8rem',
    fontWeight: 400,
    color: theme.palette.grey[700],
  },
  Icon: {
    paddingRight: '0.3rem',
  },
}));

const Views: React.FC<Props> = ({ views }: Props) => {
  const classes = useStyles();
  return (
    <Grid container alignItems="center">
      <VisibilityRounded color="secondary" className={classes.Icon} />
      <Typography className={classes.SecondaryText}>{views}</Typography>
    </Grid>
  );
};

export default Views;
