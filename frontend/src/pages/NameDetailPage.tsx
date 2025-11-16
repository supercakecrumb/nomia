import { useTranslation } from 'react-i18next';
import { useParams, Link } from 'react-router-dom';
import { useNameTrend } from '../hooks/useNameTrend';

export default function NameDetailPage() {
  const { t } = useTranslation(['pages', 'common']);
  const { name } = useParams<{ name: string }>();
  const { data, isLoading, error } = useNameTrend({ name: name || '' });

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        {/* Back Link */}
        <div className="mb-6">
          <Link
            to="/names"
            className="inline-flex items-center text-blue-600 hover:text-blue-800 font-medium"
          >
            <span className="mr-2">←</span>
            {t('pages:nameDetail.backButton')}
          </Link>
        </div>

        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-5xl font-bold text-gray-900 mb-4">
            {name}
          </h1>
          <p className="text-lg text-gray-600">
            {t('pages:nameDetail.title')}
          </p>
        </div>

        {/* Loading State */}
        {isLoading && (
          <div className="bg-white rounded-lg shadow-md p-8">
            <div className="flex items-center justify-center py-12">
              <div className="text-center">
                <div className="animate-spin text-4xl mb-4">⟳</div>
                <p className="text-gray-600">
                  {t('pages:nameDetail.loading')}
                </p>
              </div>
            </div>
          </div>
        )}

        {/* Error State */}
        {error && (
          <div className="bg-white rounded-lg shadow-md p-8">
            <div className="flex items-center justify-center py-12">
              <div className="text-center">
                <div className="text-4xl mb-4 text-red-500">✗</div>
                <p className="text-red-600 font-semibold mb-2">
                  {t('pages:nameDetail.error')}
                </p>
                <p className="text-gray-600 text-sm">
                  {error.message}
                </p>
              </div>
            </div>
          </div>
        )}

        {/* Data State */}
        {data && (
          <>
            {/* Summary Card */}
            <div className="bg-white rounded-lg shadow-md p-8 mb-6">
              <h2 className="text-2xl font-bold text-gray-900 mb-6">
                {t('pages:nameDetail.statistics')}
              </h2>
              <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
                <div className="bg-blue-50 rounded-lg p-4">
                  <div className="text-sm text-gray-600 mb-1">
                    Total Count
                  </div>
                  <div className="text-2xl font-bold text-blue-900">
                    {data.summary.total_count.toLocaleString()}
                  </div>
                </div>
                <div className="bg-pink-50 rounded-lg p-4">
                  <div className="text-sm text-gray-600 mb-1">
                    Female Count
                  </div>
                  <div className="text-2xl font-bold text-pink-900">
                    {data.summary.female_count.toLocaleString()}
                  </div>
                </div>
                <div className="bg-blue-50 rounded-lg p-4">
                  <div className="text-sm text-gray-600 mb-1">
                    Male Count
                  </div>
                  <div className="text-2xl font-bold text-blue-900">
                    {data.summary.male_count.toLocaleString()}
                  </div>
                </div>
                <div className="bg-purple-50 rounded-lg p-4">
                  <div className="text-sm text-gray-600 mb-1">
                    Gender Balance
                  </div>
                  <div className="text-2xl font-bold text-purple-900">
                    {data.summary.gender_balance !== null 
                      ? data.summary.gender_balance.toFixed(2)
                      : 'N/A'
                    }
                  </div>
                </div>
              </div>
              <div className="mt-6 pt-6 border-t border-gray-200">
                <div className="flex flex-wrap gap-4 text-sm text-gray-600">
                  <div>
                    <span className="font-semibold">Years:</span> {data.summary.name_start} - {data.summary.name_end}
                  </div>
                  <div>
                    <span className="font-semibold">Countries:</span> {data.summary.countries.join(', ')}
                  </div>
                </div>
              </div>
            </div>

            {/* Time Series Placeholder */}
            <div className="bg-white rounded-lg shadow-md p-8 mb-6">
              <h2 className="text-2xl font-bold text-gray-900 mb-6">
                {t('pages:nameDetail.trends')}
              </h2>
              <div className="flex items-center justify-center py-16 border-2 border-dashed border-gray-300 rounded-lg">
                <p className="text-gray-500 text-lg">
                  {t('pages:nameDetail.chartsComingSoon')}
                </p>
              </div>
              <div className="mt-4 text-sm text-gray-600">
                <p>Data available for {data.time_series.length} years</p>
              </div>
            </div>

            {/* By Country Placeholder */}
            <div className="bg-white rounded-lg shadow-md p-8">
              <h2 className="text-2xl font-bold text-gray-900 mb-6">
                {t('pages:nameDetail.distribution')}
              </h2>
              <div className="flex items-center justify-center py-16 border-2 border-dashed border-gray-300 rounded-lg">
                <p className="text-gray-500 text-lg">
                  {t('pages:nameDetail.chartsComingSoon')}
                </p>
              </div>
              <div className="mt-4 text-sm text-gray-600">
                <p>Data available for {data.by_country.length} countries</p>
              </div>
            </div>
          </>
        )}

        {/* No Data State */}
        {!isLoading && !error && !data && (
          <div className="bg-white rounded-lg shadow-md p-8">
            <div className="flex items-center justify-center py-12">
              <div className="text-center">
                <div className="text-4xl mb-4">🔍</div>
                <p className="text-gray-600">
                  {t('pages:nameDetail.noData')}
                </p>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}