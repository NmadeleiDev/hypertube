import i18n from 'i18next';
import Backend from 'i18next-http-backend';
import languageDetector from 'i18next-browser-languagedetector';
import { initReactI18next } from 'react-i18next';

const host = window.location.host;
i18n
  .use(Backend)
  .use(languageDetector)
  .use(initReactI18next)
  .init({
    lng: 'en',
    fallbackLng: 'en',
    backend: {
      loadPath: `/locales/{{lng}}/{{ns}}.json`,
    },
    detection: ['queryString', 'cookie'],
    cache: ['cookie'],
    debug: false,
    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
