import { Grid, Typography } from '@material-ui/core';
import React, { useEffect, useState } from 'react';
import { loadActivePeers } from '../../store/features/MoviesSlice';

interface Props {
  movieId?: string | number;
}
const initialState = { activePeers: 0, loadedPercent: 0 };
const ActivePeers = ({ movieId }: Props) => {
  const [state, setState] = useState(initialState);

  useEffect(() => {
    if (!movieId) return;
    const interval = setInterval(async () => {
      const data = await loadActivePeers(`${movieId}`);
      if (data) setState(data);
    }, 10000);
    return () => clearInterval(interval);
  }, []);
  return (
    <Grid container>
      <Typography style={{ marginRight: 5 }}>
        activePeers: {state.activePeers}
      </Typography>
      <Typography>loadedPercent: {state.loadedPercent}</Typography>
    </Grid>
  );
};

export default ActivePeers;
