import React, { useCallback, useMemo, useState } from 'react';
import { Grid, makeStyles } from '@material-ui/core';
import { useTranslation } from 'react-i18next';
import { useHistory } from 'react-router';

import { auth } from '../../axios';
import { useToast } from '../../hooks/useToast';
import { useAppDispatch } from '../../store/store';

import { getSelfInfo, saveToken } from '../../store/features/UserSlice';
import { handlersContext } from '../Dropdown/Dropdown';
import Form, { IButtonProps } from '../Form/Form';
import { IInputProps } from '../Input/Input';
import OAuthButton from '../OAuthButton/OAuthButton';

const useStyles = makeStyles({
  root: {
    padding: '2rem 1rem',
    minWidth: '30rem',
  },
  Form: {
    width: '100%',
  },
  margin: {
    marginBottom: 10,
  },
  padding: {
    padding: '5px 10px',
  },
  Link: {
    textDecoration: 'none',
  },
});

const Login = () => {
  const classes = useStyles();
  const [inputs, setInputs] = useState({ email: '', password: '' });
  const [valid, setValid] = useState({ email: false, password: false });
  const formValid = valid.email && valid.password;
  const { onClose } = React.useContext(handlersContext);
  const { t, i18n } = useTranslation();
  const { toast } = useToast();
  const history = useHistory();
  const dispatch = useAppDispatch();

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
      console.log('[Login] handleSubmit', inputs);

      const authHeader =
        'Basic ' +
        btoa(encodeURI(inputs.email) + ':' + encodeURI(inputs.password));
      try {
        const res = await auth('basic', {
          headers: {
            Authorization: authHeader,
          },
        });
        if (res.data.accessToken) {
          saveToken(res.data.accessToken);
          dispatch(getSelfInfo());
          history.push('/');
        } else {
          toast({ text: res.data[i18n.language] }, 'error');
        }
      } catch (e) {
        console.log(e);
        toast({ text: t`Server error` }, 'error');
      }
      onClose();
    },
    [inputs, toast, i18n.language, dispatch, history, onClose, t]
  );

  const handleForgot = useCallback(() => {
    console.log('[Login] handleForgot');
    history.push('forgot_password');
    onClose();
  }, [onClose, history]);

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
          onValidate: (isValid: boolean) =>
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
      ] as IInputProps[],
      buttons: [
        {
          type: 'submit',
          variant: 'contained',
          onClick: handleSubmit,
          text: t('Login'),
          disabled: !formValid,
        },
        {
          variant: 'outlined',
          onClick: handleForgot,
          text: t('Forgot password?'),
        },
        {
          to: '/register',
          variant: 'outlined',
          text: t('Register'),
          onClick: onClose,
        },
      ] as IButtonProps[],
    }),
    [inputs, handleInput, handleSubmit, handleForgot, onClose, formValid, t]
  );

  return (
    <Grid
      container
      direction="column"
      alignItems="center"
      className={classes.root}
    >
      <Form {...formData}>
        <Grid container alignItems="center" justify="space-evenly">
          <OAuthButton variant="42" />
          <OAuthButton variant="facebook" />
          <OAuthButton variant="vk" />
        </Grid>
      </Form>
    </Grid>
  );
};

export default Login;
