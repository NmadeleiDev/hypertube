import { useSnackbar } from 'notistack';
import { useEffect } from 'react';
import { useSelector } from 'react-redux';
import { removeSnack } from '../../store/features/SnackSlice';
import { RootState } from '../../store/rootReducer';
import { useAppDispatch } from '../../store/store';

let displayed: string[] = [];

const Notifier = () => {
  const dispatch = useAppDispatch();
  const notifications = useSelector(
    (state: RootState) => state.snack.notifications
  );
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  const storeDisplayed = (id: string) => {
    displayed = [...displayed, id];
  };
  const removeDisplayed = (id: string) => {
    displayed = displayed.filter((key) => id !== key);
  };

  useEffect(() => {
    notifications.forEach(({ key, text, options = {}, dismissed = false }) => {
      if (dismissed) {
        // dismiss snackbar using notistack
        closeSnackbar(key);
        return;
      }

      // do nothing if snackbar is already displayed
      if (displayed.includes(key)) return;

      // display snackbar using notistack
      enqueueSnackbar(text, {
        key,
        ...options,
        onClose: (event, reason, key) => {
          if (options.onClose) {
            options.onClose(event, reason, key);
          }
        },
        onExited: (_event, key: string) => {
          // remove this snackbar from redux store
          dispatch(removeSnack({ key }));
          removeDisplayed(key);
        },
      });

      // keep track of snackbars that we've displayed
      storeDisplayed(key);
    });
  }, [notifications, enqueueSnackbar, dispatch, closeSnackbar]);

  return null;
};

export default Notifier;
