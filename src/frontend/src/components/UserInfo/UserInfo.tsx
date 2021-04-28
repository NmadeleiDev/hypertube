import {
  Avatar,
  Button,
  Grid,
  makeStyles,
  Typography,
} from '@material-ui/core';
import React from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { useHistory } from 'react-router';
import { authFail, removeToken } from '../../store/features/UserSlice';
import { RootState } from '../../store/rootReducer';
import { useAppDispatch } from '../../store/store';
import { handlersContext } from '../Dropdown/Dropdown';

const useStyles = makeStyles({
  Type: {
    fontWeight: 100,
    marginRight: 10,
  },
  Text: {
    fontWeight: 300,
  },
  marginAuto: { marginLeft: 'auto', marginRight: 10 },
});

const UserInfo = () => {
  const classes = useStyles();
  const history = useHistory();
  const dispatch = useAppDispatch();
  const { t } = useTranslation();
  const { onClose } = React.useContext(handlersContext);
  const { username, email, imageBody } = useSelector(
    (state: RootState) => state.user
  );
  const avatar = imageBody || undefined;

  const handleLogout = () => {
    removeToken();
    dispatch(authFail());
    onClose();
  };

  const handleUpdateInfo = () => {
    onClose();
    history.push('/settings');
  };

  return (
    <Grid container>
      <Grid item container wrap="nowrap" alignItems="center">
        <Avatar src={avatar} alt={username} />
        <Button
          onClick={handleUpdateInfo}
          className={classes.marginAuto}
        >{t`Update info`}</Button>
        <Button onClick={handleLogout}>{t`Logout`}</Button>
      </Grid>
      <Grid item container alignItems="center">
        <Typography className={classes.Type}>{t`Username`}</Typography>
        <Typography className={classes.Text}>{username}</Typography>
      </Grid>
      <Grid item container alignItems="center">
        <Typography className={classes.Type}>{t`Email`}</Typography>
        <Typography className={classes.Text}>{email}</Typography>
      </Grid>
    </Grid>
  );
};

export default UserInfo;
