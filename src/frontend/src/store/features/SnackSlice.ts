import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export type IVariant = 'default' | 'info' | 'warning' | 'success' | 'error';
export interface INotification {
  key: string;
  text: string;
  variant: IVariant;
  header?: string;
  options?: any;
  dismissed?: boolean;
}

interface ISnackState {
  notifications: INotification[];
}

const initialState = {
  notifications: [],
} as ISnackState;

const snackSlice = createSlice({
  name: 'snack',
  initialState,
  reducers: {
    addSnack(state, { payload }: PayloadAction<INotification>) {
      const hasSameMessageAndVariant = state.notifications.find(
        (msg) => msg.text === payload.text && msg.variant === payload.variant
      );
      if (!hasSameMessageAndVariant) state.notifications.push(payload);
    },
    dismissSnack(state, { payload }) {
      state.notifications = state.notifications.map((notificaton) =>
        payload.dismissAll || notificaton.key === payload.key
          ? { ...notificaton, dismissed: true }
          : notificaton
      );
    },
    removeSnack(state, { payload }) {
      state.notifications = state.notifications.filter(
        (el) => el.key !== payload.key
      );
    },
  },
});

export const { addSnack, dismissSnack, removeSnack } = {
  ...snackSlice.actions,
};

export default snackSlice.reducer;
