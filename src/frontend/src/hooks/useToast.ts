import { addSnack, IVariant } from '../store/features/SnackSlice';
import { v4 as uuid } from 'uuid';
import { useAppDispatch } from '../store/store';
import { useCallback } from 'react';

interface IMessage {
  text: string;
  header?: string;
}

export const useToast = () => {
  const dispatch = useAppDispatch();

  const toast = useCallback(
    (message: IMessage, variant: IVariant = 'success', options = {}) => {
      dispatch(
        addSnack({
          key: uuid(),
          text: message.text,
          header: message.header,
          variant: variant || 'default',
          options: {
            autoHideDuration: 3000,
            ...options,
          },
        })
      );
    },
    [dispatch]
  );

  return { toast };
};
