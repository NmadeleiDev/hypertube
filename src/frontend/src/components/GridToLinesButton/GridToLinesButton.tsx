import React from 'react';
import { Button, makeStyles } from '@material-ui/core';
import { DashboardRounded, ViewAgendaRounded } from '@material-ui/icons';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/rootReducer';
import { useAppDispatch } from '../../store/store';
import { changeView } from '../../store/features/UISlice';

const useStyles = makeStyles({
  right: {
    grow: 1,
    marginLeft: 'auto',
  },
});

const GridToLinesButton = () => {
  const classes = useStyles();
  const dispatch = useAppDispatch();
  const { view } = useSelector((state: RootState) => state.UI);

  const handleClick = () => {
    dispatch(changeView());
  };

  return (
    <Button
      className={classes.right}
      variant="text"
      size="large"
      onClick={handleClick}
    >
      {view === 'grid' ? <DashboardRounded /> : <ViewAgendaRounded />}
    </Button>
  );
};

export default GridToLinesButton;
