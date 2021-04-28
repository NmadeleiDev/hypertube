import { Grid, makeStyles, Paper, Card } from '@material-ui/core';
import { Skeleton } from '@material-ui/lab';

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
  Text: {
    marginTop: 15,
    height: 10,
    width: `${Math.floor(Math.random() * 100)}%`,
  },
  Info: {
    height: '15rem',
    flexWrap: 'nowrap',
  },
});

const CardLoader = ({ display }: Props) => {
  const classes = useStyles();

  return display === 'lines' ? (
    <Paper className={classes.Paper}>
      <Skeleton animation="wave" variant="rect" className={classes.Img} />
      <Grid container direction="column" className={classes.Info}>
        <Skeleton animation="wave" className={classes.Header} />
        <Skeleton animation="wave" width="50%" className={classes.Text} />
        <Skeleton animation="wave" width="60%" className={classes.Text} />
        <Skeleton animation="wave" width="53%" className={classes.Text} />
        <Skeleton animation="wave" width="30%" className={classes.Text} />
        <Skeleton animation="wave" width="40%" className={classes.Text} />
        <Skeleton animation="wave" width="70%" className={classes.Text} />
        <Skeleton animation="wave" width="47%" className={classes.Text} />
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
