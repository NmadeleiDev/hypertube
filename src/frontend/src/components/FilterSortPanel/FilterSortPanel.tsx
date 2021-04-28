import React from 'react';
import { Grid, makeStyles } from '@material-ui/core';
import Filter from '../Filter/Filter';
import Sort from '../Sort/Sort';
import GridToLinesButton from '../GridToLinesButton/GridToLinesButton';

const useStyles = makeStyles({
  root: {
    padding: '10px 30px',
  },
});

const FilterSortPanel = () => {
  const classes = useStyles();

  return (
    <Grid container className={classes.root} alignItems="center">
      <Filter />
      <Sort />
      <GridToLinesButton />
    </Grid>
  );
};

export default FilterSortPanel;
