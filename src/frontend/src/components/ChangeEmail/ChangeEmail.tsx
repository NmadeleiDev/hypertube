import React, { useCallback, useMemo, useState } from 'react';
import { Grid, makeStyles } from '@material-ui/core';
import { useTranslation } from 'react-i18next';

import { email } from '../../axios';
import { useToast } from '../../hooks/useToast';

import Form, { IButtonProps } from '../Form/Form';
import { IInputProps } from '../Input/Input';

const useStyles = makeStyles({
  root: {
    padding: '2rem 1rem',
    minWidth: '30rem',
  },
});

const ChangeEmail = () => {
  const classes = useStyles();
  const [inputs, setInputs] = useState({ email: '' });
  const [valid, setValid] = useState({ email: false });
  const { t, i18n } = useTranslation();
  const formValid = valid.email;
  const { toast } = useToast();

  const handleInput = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      e.stopPropagation();
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
      console.log('[ChangeEmail] handleSubmit', inputs);

      // use cookie or repairToken to get new email
      try {
        const res = await email.post('patch', {
          body: { email: inputs.email },
        });
        console.log(res);
        if (res.status < 400) {
          toast(t`Check your email`);
        } else {
          console.log(res.data[i18n.language]);
          const message = {
            text: res.data[i18n.language] ?? t`Server error`,
          };

          toast(message, 'error');
        }
      } catch (e) {
        toast(t`Server error`, 'error');
      }
    },
    [inputs, t, i18n.language, toast]
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
          onValidate: (isValid) => {
            setValid((prev) => ({ ...prev, email: isValid }));
          },
          rules: {
            helperText: t('emailError'),
            rule: {
              minLength: 6,
              maxLength: 40,
              regex: /^([\w%+-.]+)@([\w-]+\.)+([\w]{2,})$/i,
            },
          },
        },
      ] as IInputProps[],
      buttons: [
        {
          type: 'submit',
          variant: 'contained',
          onClick: handleSubmit,
          text: t`Save new email`,
          disabled: !formValid,
        },
      ] as IButtonProps[],
    }),
    [inputs.email, formValid, handleInput, handleSubmit, t]
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

export default ChangeEmail;
