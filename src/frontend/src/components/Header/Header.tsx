import { ButtonProps, Grid, makeStyles } from '@material-ui/core';
import { PersonOutlineRounded } from '@material-ui/icons';
import React from 'react';
import Dropdown from '../Dropdown/Dropdown';
import Login from '../Login/Login';
import Search from '../Search/Search';
import Internationalization from '../../components/Internationalization/Internationalization';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/rootReducer';
import UserInfo from '../UserInfo/UserInfo';
import Logo from '../Logo/Logo';
import { NavLink } from 'react-router-dom';
import { resetEndOfMovies } from '../../store/features/MoviesSlice';
import { useAppDispatch } from '../../store/store';
const useStyles = makeStyles({
  root: {},
});

const Header: React.FC = () => {
  const classes = useStyles();
  const dispatch = useAppDispatch();
  const { isAuth } = useSelector((state: RootState) => state.user);

  const buttonProps: ButtonProps = {
    variant: 'outlined',
  };

  return (
    <Grid container alignItems="center" className={classes.root}>
      <Grid item xs={2}>
        <NavLink onClick={() => dispatch(resetEndOfMovies())} to="/">
          <Logo />
        </NavLink>
      </Grid>
      <Grid
        container
        justify="flex-end"
        alignItems="center"
        wrap="nowrap"
        item
        xs={8}
      >
        <Internationalization />
        <Search />
      </Grid>
      <Grid container alignItems="center" justify="center" item xs={1}></Grid>
      <Grid container alignItems="center" justify="center" item xs={1}>
        <Dropdown icon={<PersonOutlineRounded />} buttonProps={buttonProps}>
          {isAuth ? <UserInfo /> : <Login />}
        </Dropdown>
      </Grid>
    </Grid>
  );
};

export default Header;
