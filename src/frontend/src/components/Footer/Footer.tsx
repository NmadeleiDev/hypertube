import { Grid, makeStyles } from '@material-ui/core';
import { useTranslation } from 'react-i18next';

const useStyles = makeStyles({
  root: {
    width: '100%',
    padding: 10,
    position: 'fixed',
    bottom: 0,
    right: 0,
    background: 'linear-gradient(0deg, #ccc, transparent)',
  },
});

const Footer = () => {
  const styles = useStyles();
  const { t } = useTranslation();
  return (
    <Grid
      container
      alignItems="flex-end"
      justify="flex-end"
      className={styles.root}
    >
      {t`footerCreated`}
    </Grid>
  );
};

export default Footer;
