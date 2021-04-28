import { Grid, makeStyles } from '@material-ui/core';
import React, { useCallback, useState, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useHistory } from 'react-router';
import { profile } from '../../axios';
import { useToast } from '../../hooks/useToast';
import { saveUserState } from '../../store/features/UserSlice';
import { useAppDispatch } from '../../store/store';
import Form, { IButtonProps } from '../Form/Form';
import { IInputProps } from '../Input/Input';

const useStyles = makeStyles({
  root: {
    padding: '2rem 1rem',
    minWidth: '30rem',
  },
});

const Register = () => {
  const [inputs, setInputs] = useState({
    email: '',
    username: '',
    password: '',
    confirm: '',
  });
  const [valid, setValid] = useState({
    email: false,
    password: false,
    confirm: false,
  });
  const formValid = valid.email && valid.password;
  const classes = useStyles();
  const { t, i18n } = useTranslation();
  const history = useHistory();
  const dispatch = useAppDispatch();
  const { toast } = useToast();

  const handleInput = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      e.preventDefault();
      setInputs({
        ...inputs,
        [e.target.name]: e.target.value,
      });
    },
    [inputs]
  );

  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      console.log('[Register] handleSubmit', inputs);

      const body = {
        email: inputs.email,
        username: inputs.username,
        passwd: inputs.password,
      };
      const res = await profile.put('create', body);
      console.log(res);
      if (res.data.email) {
        dispatch(saveUserState({ user: res.data }));
        history.push('/');
        toast({ text: t`Check your email` });
      } else {
        toast({ text: res.data[i18n.language] }, 'error');
      }
    },
    [inputs, toast, t, i18n.language, dispatch, history]
  );

  const formData = useMemo(
    () => ({
      inputs: [
        {
          name: 'email',
          type: 'email',
          label: t('Email'),
          placeholder: t('Enter email'),
          value: inputs.email,
          onChange: handleInput,
          size: 'small',
          fullWidth: true,
          required: true,
          onValidate: (isValid) =>
            setValid((prev) => ({ ...prev, email: isValid })),
          rules: {
            helperText: t('emailError'),
            rule: {
              minLength: 6,
              maxLength: 40,
              regex: /^([\w%+-.]+)@([\w-]+\.)+([\w]{2,})$/i,
            },
          },
        },
        {
          name: 'username',
          type: 'text',
          label: t('Username'),
          placeholder: t('Enter username'),
          value: inputs.username,
          onChange: handleInput,
          size: 'small',
          fullWidth: true,
          required: true,
          onValidate: (isValid) =>
            setValid((prev) => ({ ...prev, username: isValid })),
          rules: {
            helperText: t('usernameError'),
            rule: {
              minLength: 3,
              maxLength: 20,
              regex: /^[\w-+.]+$/,
            },
          },
        },
        {
          name: 'password',
          type: 'password',
          label: t('Password'),
          placeholder: t('Enter password'),
          value: inputs.password,
          onChange: handleInput,
          size: 'small',
          fullWidth: true,
          required: true,
          onValidate: (isValid) =>
            setValid((prev) => ({ ...prev, password: isValid })),
          rules: {
            helperText: t('passwordError'),
            rule: {
              minLength: 6,
              maxLength: 25,
              regex: /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[-+_!@#$%^&*.,?]).{6,}$/,
            },
          },
        },
        {
          name: 'confirm',
          type: 'password',
          label: t('Confirm'),
          placeholder: t('Confirm password'),
          value: inputs.confirm,
          onChange: handleInput,
          size: 'small',
          fullWidth: true,
          required: true,
          onValidate: (isValid) =>
            setValid((prev) => ({ ...prev, confirm: isValid })),
          rules: {
            helperText: t('confirmError'),
            rule: {
              minLength: 3,
              maxLength: 25,
              regex: new RegExp(`^${inputs.password}$`),
            },
          },
        },
      ] as IInputProps[],
      buttons: [
        {
          type: 'submit',
          variant: 'contained',
          onClick: handleSubmit,
          text: t('Register'),
          disabled: !formValid,
        },
        {
          to: '/login',
          variant: 'outlined',
          text: t('Login'),
        },
      ] as IButtonProps[],
    }),
    [inputs, handleInput, handleSubmit, formValid, t]
  );

  return (
    <Grid
      container
      direction="column"
      alignItems="center"
      className={classes.root}
    >
      <Form {...formData} />
    </Grid>
  );
};

export default Register;
