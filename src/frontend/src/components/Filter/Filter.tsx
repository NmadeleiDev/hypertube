import React, { useState } from 'react';
import {
  Button,
  ButtonGroup,
  Checkbox,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  Grid,
  makeStyles,
} from '@material-ui/core';
import { TuneRounded } from '@material-ui/icons';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';

import {
  setFilterState,
  resetFilterState,
  GenresKeys,
  YearsKeys,
  CountriesKeys,
  FilterStateKeys,
  items,
  reqeuestMoviesWithFilters,
} from '../../store/features/FilterSlice';
import { RootState } from '../../store/rootReducer';
import { useAppDispatch } from '../../store/store';

const useStyles = makeStyles({
  '@media (orientation: portrait)': {
    root: {
      width: '80vw',
      height: '60vh',
    },
  },
  root: {
    minWidth: 300,
    minHeight: 300,
    width: '80vw',
    height: '60vh',
    '@media (max-width: 600px)': {
      width: '90vw',
      height: '90vh',
    },
  },
  Actions: {
    justifyContent: 'center',
  },
});

const buttons: ['genres', 'years', 'countries'] = [
  'genres',
  'years',
  'countries',
];

const Filter = () => {
  const classes = useStyles();
  const dispatch = useAppDispatch();
  const [open, setOpen] = useState(false);
  const [category, setCategory] = useState<FilterStateKeys>('genres');
  const filterState = useSelector((state: RootState) => state.filter);
  const { t } = useTranslation();

  const handleReset = () => {
    dispatch(resetFilterState());
  };
  const handleShow = async () => {
    console.log('[Filter] handleShow', filterState);
    dispatch(reqeuestMoviesWithFilters());
    setOpen(false);
  };

  const handleFilterChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    switch (category) {
      case 'genres':
        return dispatch(
          setFilterState({
            filter: {
              category,
              name: e.currentTarget.name as GenresKeys,
              value: e.currentTarget.checked,
            },
          })
        );
      case 'years':
        return dispatch(
          setFilterState({
            filter: {
              category,
              name: e.currentTarget.name as YearsKeys,
              value: e.currentTarget.checked,
            },
          })
        );
      case 'countries':
        return dispatch(
          setFilterState({
            filter: {
              category,
              name: e.currentTarget.name as CountriesKeys,
              value: e.currentTarget.checked,
            },
          })
        );
      default:
        return;
    }
  };

  const getCheckboxes = (category: string): JSX.Element[] | null => {
    switch (category) {
      case 'genres':
        return items.genres.map((item: GenresKeys) => (
          <Grid key={item} item xs={6} sm={4}>
            <FormControlLabel
              control={
                <Checkbox
                  style={{ marginLeft: 20 }}
                  checked={filterState['genres'][item]}
                  onChange={handleFilterChange}
                  name={item}
                  color="primary"
                />
              }
              label={t(item)}
            />
          </Grid>
        ));
      case 'years':
        return items.years.map((item: YearsKeys) => (
          <Grid key={item} item xs={6} sm={4}>
            <FormControlLabel
              control={
                <Checkbox
                  style={{ marginLeft: 20 }}
                  checked={filterState['years'][item]}
                  onChange={handleFilterChange}
                  name={item}
                  color="primary"
                />
              }
              label={t(item)}
            />
          </Grid>
        ));
      case 'countries':
        return items.countries.map((item: CountriesKeys) => (
          <Grid key={item} item xs={6} sm={4}>
            <FormControlLabel
              control={
                <Checkbox
                  style={{ marginLeft: 20 }}
                  checked={filterState['countries'][item]}
                  onChange={handleFilterChange}
                  name={item}
                  color="primary"
                />
              }
              label={t(item)}
            />
          </Grid>
        ));
      default:
        return null;
    }
  };

  return (
    <div>
      <Button
        variant="text"
        onClick={() => setOpen(true)}
        startIcon={<TuneRounded />}
      >
        {t('filterRoot')}
      </Button>
      <Dialog
        open={open}
        onClose={() => setOpen(false)}
        aria-labelledby="filter-dialog"
        classes={{ paper: classes.root }}
      >
        <DialogTitle id="filter-dialog-title">{t('filterRoot')}</DialogTitle>
        <DialogContent>
          <Grid container direction="column" alignItems="center">
            <ButtonGroup style={{ marginBottom: 20 }}>
              {buttons.map((text) => (
                <Button
                  color="primary"
                  variant={text === category ? 'contained' : 'outlined'}
                  key={text}
                  onClick={() => setCategory(text)}
                >
                  {t(text)}
                </Button>
              ))}
            </ButtonGroup>
          </Grid>
          <Grid container wrap="wrap">
            {getCheckboxes(category)}
          </Grid>
        </DialogContent>
        <DialogActions className={classes.Actions}>
          <Button onClick={handleReset} color="primary">
            Reset filters
          </Button>
          <Button onClick={handleShow} color="secondary">
            Show results
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default Filter;
