import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { useSnackbar } from 'notistack';
import Typography from '@material-ui/core/Typography';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';
import { useSelector } from 'react-redux';
import classnames from 'classnames';
import { RootState } from '../../../store/rootReducer';
import { Paper } from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  root: {
    [theme.breakpoints.up('sm')]: {
      minWidth: '344px !important',
    },
  },
  card: {
    display: 'flex',
    width: '100%',
  },
  default: {
    backgroundColor: theme.palette.grey[800],
  },
  prompt: {
    backgroundColor: theme.palette.grey[800],
    width: 550,
  },
  info: {
    backgroundColor: theme.palette.info.main,
  },
  success: {
    backgroundColor: theme.palette.success.main,
  },
  warning: {
    backgroundColor: theme.palette.warning.dark,
  },
  error: {
    backgroundColor: theme.palette.error.main,
  },
  header: {
    display: 'block',
    fontWeight: 'bold',
    color: theme.palette.common.white,
  },
  actionRoot: {
    padding: '8px 8px 8px 26px',
  },
  close: {
    marginLeft: 'auto !important',
    color: theme.palette.common.white,
  },
  text: {
    fontSize: '1rem',
    color: theme.palette.common.white,
  },
}));

interface Props {
  id: string;
}

const SnackMessage = React.forwardRef<HTMLDivElement, Props>(({ id }, ref) => {
  const classes = useStyles();
  const { closeSnackbar } = useSnackbar();
  const notifications = useSelector(
    (state: RootState) => state.snack.notifications
  );
  const selected = notifications.find((el) => el.key === id);
  const { text, variant, header } = { ...selected };

  const handleDismiss = () => {
    closeSnackbar(id);
  };

  return (
    <Paper ref={ref} className={classes.root}>
      <Card className={classnames(classes.card, classes[variant])}>
        <CardActions classes={{ root: classes.actionRoot }}>
          {header && (
            <Typography variant="subtitle2" className={classes.header}>
              {header}:
            </Typography>
          )}
          {text && (
            <Typography variant="body1" className={classes.text}>
              {text}
            </Typography>
          )}
        </CardActions>

        <IconButton className={classes.close} onClick={handleDismiss}>
          <CloseIcon color="inherit" />
        </IconButton>
      </Card>
    </Paper>
  );
});

export default SnackMessage;
