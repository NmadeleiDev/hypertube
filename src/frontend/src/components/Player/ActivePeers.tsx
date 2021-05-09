import {
  Box,
  Grid,
  LinearProgress,
  makeStyles,
  Typography,
} from '@material-ui/core';
import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { inherits } from 'util';
import { useToast } from '../../hooks/useToast';
import { loadActivePeers } from '../../store/features/MoviesSlice';

const useStyles = makeStyles({
  root: {
    width: '100%',
  },
  Text: {
    marginRight: 5,
    display: 'inline-block',
  },
  Progress: {
    marginBottom: 2,
  },
});
interface Props {
  movieId?: string | number;
}

const WARNING_MESSAGE_TIMEOUT = 30000;
const STATS_REQUEST_INTERVAL = 10000;

const initialState = { activePeers: 0, loadedPercent: 0 };
const ActivePeers = ({ movieId }: Props) => {
  const classes = useStyles();
  const [state, setState] = useState(initialState);
  const { toast } = useToast();
  const { t } = useTranslation();

  useEffect(() => {
    if (!movieId) return;
    const timeout = setTimeout(() => {
      if (!state.activePeers || state.loadedPercent < 100) {
        toast({ text: t`loadingTimeoutWarning` }, 'warning', {
          autoHideDuration: 7000,
        });
      }
    }, WARNING_MESSAGE_TIMEOUT);
    const interval = setInterval(async () => {
      const data = await loadActivePeers(`${movieId}`);
      if (data) setState(data);
      if (data?.activePeers || data?.loadedPercent === 100)
        clearTimeout(timeout);
    }, STATS_REQUEST_INTERVAL);
    return () => {
      clearInterval(interval);
      clearTimeout(timeout);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <Grid container direction="column" className={classes.root}>
      <Box>
        <Typography component="span" className={classes.Text}>
          Active peers:
        </Typography>
        <Typography
          className={classes.Text}
          style={{ color: state.activePeers < 5 ? 'darkred' : 'darkgreen' }}
        >
          {state.activePeers}
        </Typography>
      </Box>
      <Box>
        <Typography className={classes.Text}>Loaded:</Typography>
        <LinearProgress
          className={classes.Progress}
          variant="determinate"
          value={state.loadedPercent}
        />
      </Box>
    </Grid>
  );
};

export default ActivePeers;
