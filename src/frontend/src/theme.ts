import { createMuiTheme } from '@material-ui/core';

export const primaryColor = {
  light: '#05829E',
  main: '#035263',
  dark: '#02313B',
  contrastText: '#fff',
};
export const secondaryColor = {
  light: '#F19941',
  main: '#EE851B',
  dark: '#D17010',
  contrastText: '#fff',
  contrastTextLighter: '#333',
};
export const backgroundColor = {
  light: '#FFFFFF',
  default: '#FAFAFA',
  dark: '#DDDDDD',
  contrastText: '#121212',
};

export const theme = createMuiTheme({
  palette: {
    primary: primaryColor,
    secondary: secondaryColor,
    background: backgroundColor,
  },
  typography: {
    fontFamily: [
      '"Quicksand"',
      '"Helvetica Neue"',
      'Roboto',
      'Arial',
      '-apple-system',
      'BlinkMacSystemFont',
      '"Segoe UI"',
      'sans-serif',
      '"Apple Color Emoji"',
      '"Segoe UI Emoji"',
      '"Segoe UI Symbol"',
    ].join(','),
  },
  props: {
    MuiButton: {
      variant: 'outlined',
      color: 'primary',
    },
    MuiTextField: {
      variant: 'outlined',
      color: 'primary',
    },
  },
});
