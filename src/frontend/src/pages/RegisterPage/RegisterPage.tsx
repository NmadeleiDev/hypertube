import { Container, makeStyles } from '@material-ui/core';
import React from 'react';
import Register from '../../components/Register/Register';

const useStyles = makeStyles({
  root: {},
});

const RegisterPage = () => {
  const classes = useStyles();
  return (
    <Container className={classes.root}>
      <Register />
    </Container>
  );
};

export default RegisterPage;
