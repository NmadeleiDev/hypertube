import { Grid, Typography } from '@material-ui/core';
import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useToast } from '../../hooks/useToast';
import { loadActivePeers } from '../../store/features/MoviesSlice';

interface Props {
  movieId?: string | number;
}

const WARNING_MESSAGE_TIMEOUT = 30000;
const STATS_REQUEST_INTERVAL = 10000;

const initialState = { activePeers: 0, loadedPercent: 0 };
const ActivePeers = ({ movieId }: Props) => {
  const [state, setState] = useState(initialState);
  const { toast } = useToast();
  const { t } = useTranslation();

  useEffect(() => {
    if (!movieId) return;
    const timeout = setTimeout(() => {
      if (!state.activePeers) {
        toast({ text: t`loadingTimeoutWarning` }, 'warning', {
          autoHideDuration: 7000,
        });
      }
    }, WARNING_MESSAGE_TIMEOUT);
    const interval = setInterval(async () => {
      const data = await loadActivePeers(`${movieId}`);
      if (data) setState(data);
      if (data?.activePeers) clearTimeout(timeout);
    }, STATS_REQUEST_INTERVAL);
    return () => clearInterval(interval);
  }, []);

  return (
    <Grid container>
      <Typography style={{ marginRight: 5 }}>
        Active peers: {state.activePeers}
      </Typography>
      <Typography>Percent loaded: {state.loadedPercent}%</Typography>
    </Grid>
  );
};

export default ActivePeers;
