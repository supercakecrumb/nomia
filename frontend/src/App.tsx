import { useTranslation } from 'react-i18next';
import LanguageSwitcher from './components/LanguageSwitcher';
import { useMetaYears } from './hooks/useMetaYears';
import { getApiMode } from './api/client';

function App() {
  const { t } = useTranslation(['common', 'pages']);
  const { data: yearData, isLoading, error } = useMetaYears();

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

          {/* API Client Verification */}
          <div className="mt-8 p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <h3 className="text-lg font-semibold text-blue-900 mb-2">
              API Client Status
            </h3>
            <div className="space-y-2 text-sm">
              <div className="flex items-center gap-2">
                <span className="font-semibold text-blue-700">Mode:</span>
                <span className="text-blue-900 font-mono">{getApiMode()}</span>
              </div>
              
              {isLoading && (
                <div className="flex items-center gap-2 text-blue-700">
                  <span className="animate-spin">⟳</span>
                  <span>Loading year range...</span>
                </div>
              )}
              
              {error && (
                <div className="flex items-center gap-2 text-red-700">
                  <span className="font-semibold">✗</span>
                  <span>Error: {error.message}</span>
                </div>
              )}
              
              {yearData && (
                <div className="flex items-center gap-2 text-green-700">
                  <span className="font-semibold">✓</span>
                  <span>
                    Data available: {yearData.min_year} - {yearData.max_year}
                  </span>
                </div>
              )}
            </div>
          </div>
          
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
