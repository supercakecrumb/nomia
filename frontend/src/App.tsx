import { useTranslation } from 'react-i18next';
import LanguageSwitcher from './components/LanguageSwitcher';

function App() {
  const { t } = useTranslation(['common', 'pages']);

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow-md">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-gray-800">
            {t('common:appName')}
          </h1>
          <LanguageSwitcher />
        </div>
      </header>
      
      <main className="container mx-auto px-4 py-8">
        <div className="bg-white p-8 rounded-lg shadow-md max-w-2xl mx-auto">
          <h2 className="text-3xl font-bold text-gray-800 mb-4">
            {t('pages:main.title')}
          </h2>
          <p className="text-xl text-gray-600 mb-4">
            {t('pages:main.subtitle')}
          </p>
          <p className="text-gray-600 mb-6">
            {t('pages:main.description')}
          </p>
          
          <div className="space-y-2 text-sm text-gray-500 mt-8">
            <div className="flex items-center gap-2">
              <span className="font-semibold">✓</span>
              <span>React 19.2.0</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="font-semibold">✓</span>
              <span>React Router 7.9.6</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="font-semibold">✓</span>
              <span>TanStack Query 5.90.9</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="font-semibold">✓</span>
              <span>Tailwind CSS 4.1.17</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="font-semibold">✓</span>
              <span>i18next 23.16.8 + react-i18next 16.3.3</span>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

export default App
