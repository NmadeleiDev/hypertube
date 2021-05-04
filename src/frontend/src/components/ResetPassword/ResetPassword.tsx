import React, { useMemo, useState } from 'react';
import { Grid, makeStyles } from '@material-ui/core';
import { useTranslation } from 'react-i18next';

import { passwd } from '../../axios';
import { useToast } from '../../hooks/useToast';

import Form, { IButtonProps } from '../Form/Form';
import { IInputProps } from '../Input/Input';
import { getSearchParam } from '../../utils';
import { useHistory } from 'react-router';

const useStyles = makeStyles({
  root: {
    padding: '2rem 1rem',
    minWidth: '30rem',
  },
});

const ResetPassword = () => {
  const classes = useStyles();
  const [inputs, setInputs] = useState({ password: '', confirm: '' });
  const [valid, setValid] = useState({ password: false, confirm: false });
  const { t, i18n } = useTranslation();
  const formValid = valid.password && valid.confirm;
  const { toast } = useToast();
  const history = useHistory();

  const formData = useMemo(() => {
    const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
      e.stopPropagation();
      setInputs({
        ...inputs,
        [e.target.name]: e.target.value,
      });
    };

    const handleSubmit = async (e: React.FormEvent) => {
      e.preventDefault();
      console.log('[ForgotPassword] handleSubmit', inputs);

      const searchParams = getSearchParam();
      if (!searchParams?.repairToken) {
        toast({ text: t`ServerErrorTry` }, 'error');
        history.push('/forgot_password');
        return;
      }
      try {
        const repairToken = searchParams.repairToken.replaceAll('%3D', '=');
        const res = await passwd.patch(
          '/repair/patch',
          { passwd: inputs.password },
          { headers: { repairToken } }
        );
        console.log(res);
        if (res.status < 400) {
          toast({ text: t`Password changed` });
          history.push('/login');
        } else {
          toast({ text: res.data[i18n.language] ?? t`Server error` }, 'error');
        }
      } catch (e) {
        toast(t`Server error`, 'error');
      }
    };

    return {
      inputs: [
        {
          name: 'password',
          type: 'password',
          label: t`New password`,
          placeholder: t`Enter password`,
          value: inputs.password,
          onChange: handleInput,
          size: 'small',
          fullWidth: true,
          required: true,
          onValidate: (isValid) =>
            setValid((prev) => ({ ...prev, password: isValid })),
          rules: {
            helperText: t`passwordError`,
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
          label: t`Confirm`,
          placeholder: t`Confirm password`,
          value: inputs.confirm,
          onChange: handleInput,
          size: 'small',
          fullWidth: true,
          required: true,
          onValidate: (isValid) =>
            setValid((prev) => ({ ...prev, confirm: isValid })),
          rules: {
            helperText: t`confirmError`,
            rule: {
              minLength: 6,
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
          text: t`Save new password`,
          disabled: !formValid,
        },
      ] as IButtonProps[],
    };
  }, [inputs, formValid, t, i18n.language, toast, history]);

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

export default ResetPassword;
