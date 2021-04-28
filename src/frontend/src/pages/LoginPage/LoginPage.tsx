import { Container, makeStyles } from '@material-ui/core';
import React from 'react';
import Login from '../../components/Login/Login';

const useStyles = makeStyles({
  root: {},
});

const LoginPage = () => {
  const classes = useStyles();
  return (
    <Container className={classes.root}>
      <Login />
    </Container>
  );
};

export default LoginPage;
