import { debounce, throttle } from 'lodash';

/**
 * find scroll position at the end of scroll event
 * fires callback1 if we're reaching bottom line of a given element
 * (some pixels reserved to load in advance)
 * isEndOfMoviesRef.current is updated in a callback after fetched result
 */
export const throttledDetectBottomLine = (
  element: HTMLElement | null,
  callback1: Function,
  callback2?: Function,
  delay: number = 500,
  reserve: number = 100
) =>
  throttle(() => {
    if (!element) return;
    const scroll = window.pageYOffset + window.innerHeight;
    const offset = element.offsetTop + element.offsetHeight;

    if (scroll > offset - reserve) {
      callback1 && callback1();
    } else {
      callback2 && callback2();
    }
  }, delay);

/**
 * find scroll position at the end of scroll event
 * fires callback1 if we're reaching bottom line of a given element
 * (some pixels reserved to load in advance)
 * isEndOfMoviesRef.current is updated in a callback after fetched result
 */
export const debouncedDetectBottomLine = (
  element: HTMLElement | null,
  callback1: Function,
  callback2?: Function,
  delay: number = 500,
  reserve: number = 100
) =>
  debounce(() => {
    if (!element) return;
    const scroll = window.pageYOffset + window.innerHeight;
    const offset = element.offsetTop + element.offsetHeight;

    if (scroll > offset - reserve) {
      callback1 && callback1();
    } else {
      callback2 && callback2();
    }
  }, delay);

interface KnownOptions {
  path?: string;
  expires?: Date | string;
}
interface IOptions extends KnownOptions {
  [key: string]: string | Date | boolean | undefined;
}

export const setCookie = (
  name: string,
  value: string,
  options: IOptions = {}
) => {
  options = {
    path: '/',
    // при необходимости добавьте другие значения по умолчанию
    ...options,
  };

  if (options.expires instanceof Date) {
    options.expires = options.expires.toUTCString();
  }

  let updatedCookie =
    encodeURIComponent(name) + '=' + encodeURIComponent(value);

  Object.keys(options).forEach((key) => {
    updatedCookie += '; ' + key;
    let optionValue = options[key];
    if (optionValue !== true) {
      updatedCookie += '=' + optionValue;
    }
  });
  document.cookie = updatedCookie;
};

export const getSearchParam = (search?: string) => {
  const searchString = search || window.location.search;
  const params = searchString ? decodeURI(searchString) : null;
  return params
    ? params
        .replaceAll('%3A', ':')
        .replaceAll('%2C', ',')
        .replaceAll('+', ' ')
        .slice(1)
        .split('&')
        .reduce((acc, cur) => {
          const keyValue = cur.split('=');
          acc[keyValue[0]] =
            keyValue[0] === 'error' ? JSON.parse(keyValue[1]) : keyValue[1];
          return acc;
        }, {} as { [key: string]: any })
    : null;
};
