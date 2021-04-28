import React from 'react';
import { useHistory } from 'react-router-dom';
import { getSelfInfo, saveToken } from '../../store/features/UserSlice';
import { useAppDispatch } from '../../store/store';
import { setCookie } from '../../utils';

const Auth = () => {
  const history = useHistory();
  const dispatch = useAppDispatch();

  React.useEffect(() => {
    console.log('cookie', document.cookie);
    console.log(window.location.href);
    const queryParams = decodeURI(window.location.search);
    console.log('queryParams', queryParams);
    const accessToken = queryParams
      .slice(1)
      .split('&')
      .find((el) => el.startsWith('accessToken'));
    if (accessToken) {
      const tokenStart = accessToken.search(/=(.+)/);
      const token = accessToken.slice(tokenStart + 1);
      console.log('tokenStart', token);
      saveToken(token);
      dispatch(getSelfInfo());
      setCookie('accessToken', token, { path: '/' });
      history.push('/');
    }
  }, [dispatch, history]);
  return null;
};

export default Auth;
