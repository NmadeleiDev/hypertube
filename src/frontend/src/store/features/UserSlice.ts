import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { profile } from '../../axios';
import { AppDispatch } from '../store';

const tokenKey = 'hypertubeAuthToken';

interface IUserPayload {
  user: IUser;
}

export interface IUser {
  userId: number;
  email: string;
  username: string;
  firstName: string;
  lastName: string;
  imageBody: string | null;
}

export type IUserState = IUser & {
  isAuth: boolean;
  isLoading: boolean;
};

const initialState: IUserState = {
  isAuth: false,
  isLoading: false,
  userId: 0,
  email: '',
  username: '',
  firstName: '',
  lastName: '',
  imageBody: null,
};

const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    saveUserState(state, { payload }: PayloadAction<IUserPayload>) {
      if (!payload.user)
        throw new Error('[userSlice:saveState] no user in payload');
      //   Object.keys(payload.user).map((k) => {
      //     const key = k as IUserKeys;
      //     state[key] = payload.user[key];
      //   });
      return {
        isAuth: state.isAuth,
        isLoading: state.isLoading,
        userId: payload.user.userId || state.userId,
        email: payload.user.email || state.email,
        username: payload.user.username || state.username,
        firstName: payload.user.firstName || state.firstName,
        lastName: payload.user.lastName || state.lastName,
        imageBody: payload.user.imageBody || state.imageBody,
      };
    },
    authSuccess(state) {
      state.isAuth = true;
    },
    authFail(state) {
      state.isAuth = false;
    },
    startLoading(state) {
      state.isLoading = true;
    },
    stopLoading(state) {
      state.isLoading = false;
    },
  },
});

export const getSelfInfo = () => async (dispatch: AppDispatch) => {
  dispatch(startLoading());
  try {
    const res = await profile('get', {
      headers: {
        accessToken: getToken(),
      },
    });
    console.log(res.data);
    if (res.data.en) {
      dispatch(authFail());
      removeToken();
    } else {
      dispatch(
        saveUserState({
          user: res.data as IUser,
        })
      );
      dispatch(authSuccess());
    }
  } catch (e) {
    console.log(e);
  }
  dispatch(stopLoading());
};

export const saveToken = (token: string) => {
  localStorage.setItem(tokenKey, token);
};

export const removeToken = (key: string = tokenKey) => {
  localStorage.removeItem(key);
};

export const getToken = () => localStorage.getItem(tokenKey);

export const {
  saveUserState,
  authSuccess,
  authFail,
  startLoading,
  stopLoading,
} = userSlice.actions;

export default userSlice.reducer;
