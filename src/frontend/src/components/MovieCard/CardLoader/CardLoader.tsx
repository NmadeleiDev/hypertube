import { useEffect, useState } from 'react';
import { Grid, makeStyles, Paper, Card, Typography } from '@material-ui/core';
import { Skeleton } from '@material-ui/lab';
import { useTranslation } from 'react-i18next';

interface Props {
  display?: 'lines' | 'grid';
}

const useStyles = makeStyles({
  Paper: {
    display: 'flex',
    marginBottom: 20,
    padding: 10,
    width: '100%',
  },
  Img: {
    borderRadius: 5,
    height: '15rem',
    width: '13rem',
    marginRight: '1rem',
  },
  Header: {
    width: '90%',
    marginBottom: 10,
    fontSize: 30,
  },
  Line: {
    marginBottom: 15,
    height: 10,
    width: `${Math.floor(Math.random() * 100)}%`,
  },
  Info: {
    height: '15rem',
    flexWrap: 'nowrap',
  },
  Text: {
    fontSize: '1rem',
    height: 10,
    marginBottom: 15,
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    textTransform: 'uppercase',
  },
});

const lines = ['line0', 'line1', 'line2', 'line3', 'line4', 'line5', 'line6'];

const CardLoader = ({ display }: Props) => {
  const classes = useStyles();
  const [index, setIndex] = useState<number>(0);
  const { t } = useTranslation();

  useEffect(() => {
    const interval = setInterval(() => {
      setIndex((prev) => (prev + 1) % lines.length);
    }, 3000);
    return () => clearInterval(interval);
  });

  return display === 'lines' ? (
    <Paper className={classes.Paper}>
      <Skeleton animation="wave" variant="rect" className={classes.Img} />
      <Grid container direction="column" className={classes.Info}>
        <Skeleton animation="wave" className={classes.Header} />
        <Skeleton animation="wave" width="50%" className={classes.Line} />
        <Skeleton animation="wave" width="60%" className={classes.Line} />
        <Skeleton animation="wave" width="53%" className={classes.Line} />
        {/* <Skeleton animation="wave" width="30%" className={classes.Line} /> */}
        <Typography className={classes.Text}>{t(lines[index])}</Typography>
        <Skeleton animation="wave" width="40%" className={classes.Line} />
        <Skeleton animation="wave" width="70%" className={classes.Line} />
        <Skeleton animation="wave" width="47%" className={classes.Line} />
      </Grid>
    </Paper>
  ) : (
    <Card style={{ height: 'fit-content' }}>
      <div
        className={classes.Img}
        style={{
          backgroundSize: 'cover',
        }}
      >
        <Skeleton animation="wave" variant="rect" className={classes.Img} />
      </div>
    </Card>
  );
};

export default CardLoader;
