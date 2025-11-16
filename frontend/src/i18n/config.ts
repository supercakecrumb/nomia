import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import HttpBackend from 'i18next-http-backend';

// Initialize i18next with configuration
i18n
  .use(HttpBackend)
  .use(initReactI18next)
  .init({
    // Set default language
    lng: localStorage.getItem('language') || 'en',
    fallbackLng: 'en',
    
    // Enable debug mode in development
    debug: import.meta.env.DEV,
    
    // Configure namespaces
    ns: ['common', 'filters', 'pages'],
    defaultNS: 'common',
    
    // Interpolation settings
    interpolation: {
      escapeValue: false, // React already escapes values
    },
    
    // Backend configuration for loading translations
    backend: {
      loadPath: '/locales/{{lng}}/{{ns}}.json',
    },
    
    // React specific options
    react: {
      useSuspense: false,
    },
  });

export default i18n;