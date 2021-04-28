import React, { useState } from 'react';
import { Button, makeStyles } from '@material-ui/core';
import { useTranslation } from 'react-i18next';
import { theme } from '../../theme';

const languages = ['en', 'ru'];

const useStyles = makeStyles({
  Button: {
    [theme.breakpoints.down('sm')]: {
      fontSize: '0.7rem',
    },
  },
});

const Internationalization = () => {
  const classes = useStyles();
  const [lang, setLang] = useState(0);
  const nextLang = (lang + 1) % languages.length;
  const { t, i18n } = useTranslation();

  const handleClick = () => {
    i18n.changeLanguage(languages[nextLang]);
    setLang(nextLang);
  };
  return (
    <div>
      <Button onClick={handleClick} className={classes.Button}>{`${t(
        'translate'
      )} ${t(languages[nextLang])}`}</Button>
    </div>
  );
};

export default Internationalization;
