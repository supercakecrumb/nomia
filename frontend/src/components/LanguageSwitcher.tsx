import { useTranslation } from 'react-i18next';

const LanguageSwitcher = () => {
  const { i18n, t } = useTranslation('common');

  const toggleLanguage = () => {
    const newLang = i18n.language === 'en' ? 'ru' : 'en';
    i18n.changeLanguage(newLang);
    localStorage.setItem('language', newLang);
  };

  const currentLang = i18n.language;

  return (
    <button
      onClick={toggleLanguage}
      className="flex items-center gap-2 px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors duration-200 shadow-md hover:shadow-lg"
      aria-label={t('labels.language')}
    >
      <span className="text-lg font-semibold">
        {currentLang === 'en' ? '🇺🇸 EN' : '🇷🇺 RU'}
      </span>
      <span className="text-sm opacity-90">
        {t(`languages.${currentLang}`)}
      </span>
    </button>
  );
};

export default LanguageSwitcher;