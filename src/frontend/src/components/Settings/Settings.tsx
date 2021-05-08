import { Avatar, Grid, makeStyles, Typography } from '@material-ui/core';
import React, { useState, useMemo, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { email, passwd, profile } from '../../axios';
import Form, { IButtonProps } from '../Form/Form';
import { IInputProps } from '../Input/Input';
import defaultAvatar from '../../images/defaultAvatar.png';
import { BackupRounded } from '@material-ui/icons';
import { RootState } from '../../store/rootReducer';
import { useSelector } from 'react-redux';
import { useToast } from '../../hooks/useToast';
import { getToken, saveUserState } from '../../store/features/UserSlice';
import { useAppDispatch } from '../../store/store';
import { useHistory } from 'react-router';

const useStyles = makeStyles({
  root: {
    padding: '2rem 1rem',
    minWidth: '30rem',
  },
  avatarWrapper: {
    position: 'relative',
    '&:hover': {
      cursor: 'pointer',
      '& $uploadIcon': {
        opacity: 0.7,
      },
      '& $avatarImage': {
        opacity: 0.5,
      },
    },
  },
  avatarImage: {
    width: '100px',
    height: '100px',
    border: '2px solid rgba(3, 82, 99, 0.5)',
    opacity: 1,
    transition: '0.2s',
  },
  inputNone: {
    display: 'none',
  },
  uploadIcon: {
    position: 'absolute',
    transition: '0.2s',
    top: 20,
    left: 20,
    width: 60,
    height: 60,
    opacity: 0,
    zIndex: 2,
  },
});

const Settings = () => {
  const user = useSelector((state: RootState) => state.user);
  const [inputs, setInputs] = useState({
    email: user.email,
    username: user.username,
    firstName: user.firstName || '',
    lastName: user.lastName || '',
    newPassword: '',
    currentPassword: '',
    confirm: '',
    avatarImage: user.imageBody,
  });
  const [valid, setValid] = useState({
    email: true,
    newPassword: true,
    currentPassword: true,
    confirm: true,
    username: true,
    firstName: true,
    lastName: true,
  });
  const classes = useStyles();
  const { t, i18n } = useTranslation();
  const inputRef = useRef<HTMLInputElement>(null);
  const { toast } = useToast();
  const dispatch = useAppDispatch();
  const history = useHistory();

  const validateForm = () => {
    const {
      email,
      username,
      currentPassword,
      confirm,
      firstName,
      lastName,
    } = valid;
    return inputs.newPassword
      ? currentPassword &&
          confirm &&
          (email || username || firstName || lastName)
      : email || username || firstName || lastName;
  };
  const formValid = validateForm();

  const validateUpload = (file: File) => {
    return (
      file.type === 'image/png' ||
      file.type === 'image/jpg' ||
      file.type === 'image/jpeg'
    );
  };

  const uploadAvatar = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files) return;
    let file = e.target.files[0];
    if (!validateUpload(file)) {
      e.target.value = '';
      return toast({ text: 'Please add image only!' });
    }
    if (file.size > 2000000) {
      e.target.value = '';
      return toast({ text: 'Please upload picture size less then 2 mb!' });
    }
    let reader = new FileReader();
    reader.onload = async function (e) {
      if (!e.target || !e.target.result) return;
      const image = e.target.result.toString();
      setInputs({
        ...inputs,
        avatarImage: image,
      });
      const body = {
        imageBody: image,
      };
      const res = await profile.patch('patch', body, {
        headers: {
          accessToken: getToken(),
        },
      });
      // console.log(res.data);
      if (res.status >= 400) toast({ text: t`Avatar upload error` }, 'error');
      else {
        dispatch(
          saveUserState({
            user: {
              ...user,
              imageBody: image,
            },
          })
        );
        toast({ text: t`Avatar uploaded` });
      }
    };
    reader.readAsDataURL(file);
  };

  const formData = useMemo(() => {
    const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
      e.preventDefault();
      setInputs((prevInputs) => ({
        ...prevInputs,
        [e.target.name]: e.target.value,
      }));
    };

    const updateEmail = async () => {
      const body = {
        email: inputs.email,
      };
      const res = await email.patch('patch', body, {
        headers: {
          accessToken: getToken(),
        },
      });
      console.log(res.data);
      if (res.status >= 400)
        toast({ text: res.data[i18n.language] || t`Update error` }, 'error');
      else {
        dispatch(
          saveUserState({
            user: {
              ...user,
              email: inputs.email,
            },
          })
        );
        toast({ text: t`Check your email` });
      }
    };

    const updateInfo = async () => {
      const body = {
        username: inputs.username,
        firstName: inputs.firstName,
        lastName: inputs.lastName,
      };
      const res = await profile.patch('patch', body, {
        headers: {
          accessToken: getToken(),
        },
      });
      console.log(res.data);
      if (res.status >= 400)
        toast({ text: res.data[i18n.language] || t`Update error` }, 'error');
      else {
        dispatch(
          saveUserState({
            user: {
              ...user,
              username: inputs.username,
              firstName: inputs.firstName,
              lastName: inputs.lastName,
            },
          })
        );
        toast({ text: t`User info updated` });
      }
    };

    const updatePassword = async () => {
      const body = {
        passwd: inputs.currentPassword,
        newPasswd: inputs.newPassword,
      };
      const res = await passwd.patch('patch', body, {
        headers: {
          accessToken: getToken(),
        },
      });
      console.log(res.data);
      if (res.status >= 400)
        toast({ text: res.data[i18n.language] || t`Update error` }, 'error');
      else {
        toast({ text: t`User info updated` });
        setInputs({
          ...inputs,
          currentPassword: '',
          confirm: '',
          newPassword: '',
        });
      }
    };

    const handleSubmit = async (e: React.FormEvent) => {
      e.preventDefault();
      if (inputs.email !== user.email) updateEmail();
      if (
        inputs.username !== user.username ||
        inputs.firstName !== user.firstName ||
        inputs.lastName !== user.lastName
      )
        updateInfo();
      if (inputs.newPassword) updatePassword();
      console.log('[Settings] handleSubmit', inputs);
    };

    return {
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
          onValidate: (isValid) =>
            setValid((prev) => ({
              ...prev,
              currentPassword: isValid,
            })),
          rules: {
            helperText: t('emailError'),
            rule: {
              minLength: 0,
              maxLength: 40,
              regex: /^(([\w%+-.]+)@([\w-]+\.)+([\w]{2,}))?$/i,
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
          onValidate: (isValid) =>
            setValid((prev) => ({
              ...prev,
              currentPassword: isValid,
            })),
          rules: {
            helperText: t('usernameError'),
            rule: {
              minLength: 0,
              maxLength: 20,
              regex: /^([\w.+-]+)?$/,
            },
          },
        },
        {
          name: 'firstName',
          type: 'text',
          label: t('First Name'),
          placeholder: t('Enter first name'),
          value: inputs.firstName,
          onChange: handleInput,
          size: 'small',
          fullWidth: true,
          onValidate: (isValid) =>
            setValid((prev) => ({
              ...prev,
              currentPassword: isValid,
            })),
          rules: {
            helperText: t('firstNameError'),
            rule: {
              minLength: 0,
              maxLength: 20,
              regex: /^([\w.+-]+)?$/,
            },
          },
        },
        {
          name: 'lastName',
          type: 'text',
          label: t('Last Name'),
          placeholder: t('Enter last name'),
          value: inputs.lastName,
          onChange: handleInput,
          size: 'small',
          fullWidth: true,
          onValidate: (isValid) =>
            setValid((prev) => ({
              ...prev,
              currentPassword: isValid,
            })),
          rules: {
            helperText: t('lastNameError'),
            rule: {
              minLength: 0,
              maxLength: 20,
              regex: /^([\w.+-]+)?$/,
            },
          },
        },
        {
          name: 'newPassword',
          type: 'password',
          label: t('New password'),
          placeholder: t('Enter new password'),
          value: inputs.newPassword,
          onChange: handleInput,
          size: 'small',
          fullWidth: true,
          onValidate: (isValid) =>
            setValid((prev) => ({
              ...prev,
              currentPassword: isValid,
            })),
          rules: {
            helperText: t('passwordError'),
            rule: {
              minLength: 0,
              maxLength: 25,
              regex: /^((?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[-+_!@#$%^&*.,?]).{6,})?$/,
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
          onValidate: (isValid) =>
            setValid((prev) => ({ ...prev, confirm: isValid })),
          rules: {
            helperText: t('confirmError'),
            rule: {
              minLength: 0,
              maxLength: 25,
              regex: new RegExp(`^${inputs.newPassword}$`),
            },
          },
        },
        {
          name: 'currentPassword',
          type: 'password',
          label: t('Current password'),
          placeholder: t('Enter current password'),
          value: inputs.currentPassword,
          onChange: handleInput,
          size: 'small',
          fullWidth: true,
          onValidate: (isValid) =>
            setValid((prev) => ({
              ...prev,
              currentPassword: isValid,
            })),
          rules: {
            helperText: t('currentPasswordError'),
            rule: {
              minLength: inputs.newPassword.length ? 1 : 0,
              maxLength: 25,
            },
          },
        },
      ] as IInputProps[],
      buttons: [
        {
          type: 'submit',
          variant: 'contained',
          onClick: handleSubmit,
          text: t('Save'),
          disabled: !formValid,
        },
        {
          variant: 'outlined',
          onClick: () => history.goBack(),
          text: t`Back`,
        },
      ] as IButtonProps[],
    };
  }, [inputs, formValid, t, dispatch, i18n.language, toast, user, history]);

  return (
    <Grid
      container
      direction="column"
      alignItems="center"
      className={classes.root}
    >
      <Typography>{t('Settings').toUpperCase()}</Typography>
      <div
        className={classes.avatarWrapper}
        onClick={() => {
          if (inputRef && inputRef.current) {
            inputRef.current.click();
          }
        }}
      >
        <Avatar
          className={classes.avatarImage}
          alt="profile image"
          src={inputs.avatarImage || defaultAvatar}
        />
        <BackupRounded className={classes.uploadIcon} />
        <input
          className={classes.inputNone}
          type="file"
          accept="image/png, image/jpeg, image/jpg"
          ref={inputRef}
          onChange={uploadAvatar}
        />
      </div>
      <Form {...formData} />
    </Grid>
  );
};

export default Settings;
