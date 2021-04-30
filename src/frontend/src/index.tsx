import React, { Suspense } from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import { BrowserRouter } from 'react-router-dom';
import { CssBaseline, ThemeProvider } from '@material-ui/core';
import { theme } from './theme';

import { Provider } from 'react-redux';
import store from './store/store';

import './i18n';

if (process.env.NODE_ENV === 'production') console.log = () => {};
export const LIMIT = +(process.env.REACT_APP_LOAD_LIMIT || 5);

// console.log = () => {};

const render = () => {
  const App = require('./App').default;

  ReactDOM.render(
    // <React.StrictMode>
    <Provider store={store}>
      <ThemeProvider theme={theme}>
        <BrowserRouter>
          <Suspense fallback={<div>Loading...</div>}>
            <CssBaseline />
            <App />
          </Suspense>
        </BrowserRouter>
      </ThemeProvider>
    </Provider>,
    // </React.StrictMode>,
    document.getElementById('root')
  );
};

render();

if (process.env.NODE_ENV === 'development' && module.hot) {
  module.hot.accept('./App', render);
}
