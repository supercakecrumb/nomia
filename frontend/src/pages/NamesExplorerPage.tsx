import { useTranslation } from 'react-i18next';
import { useNames } from '../hooks/useNames';

export default function NamesExplorerPage() {
  const { t } = useTranslation('pages');
  const { data, isLoading, error } = useNames();

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">
            {t('namesExplorer.title')}
          </h1>
          <p className="text-lg text-gray-600">
            {t('namesExplorer.subtitle')}
          </p>
        </div>

        {/* Connection Status */}
        {data && (
          <div className="mb-6 p-3 bg-green-50 border border-green-200 rounded-lg">
            <p className="text-sm text-green-800">
              ✓ API Connected - {data.names.length} names loaded
            </p>
          </div>
        )}

        {isLoading && (
          <div className="mb-6 p-3 bg-blue-50 border border-blue-200 rounded-lg">
            <p className="text-sm text-blue-800">
              ⟳ {t('common:labels.loading')}
            </p>
          </div>
        )}

        {error && (
          <div className="mb-6 p-3 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-sm text-red-800">
              ✗ {t('common:labels.error')}: {error.message}
            </p>
          </div>
        )}

        {/* Filter Bar Placeholder */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Filters
          </h2>
          <div className="flex items-center justify-center py-12 border-2 border-dashed border-gray-300 rounded-lg">
            <p className="text-gray-500 text-lg">
              {t('namesExplorer.comingSoon.filters')}
            </p>
          </div>
        </div>

        {/* Names Table Placeholder */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Names
          </h2>
          <div className="flex items-center justify-center py-20 border-2 border-dashed border-gray-300 rounded-lg">
            <p className="text-gray-500 text-lg">
              {t('namesExplorer.comingSoon.table')}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}