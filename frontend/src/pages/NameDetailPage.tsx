import { useTranslation } from 'react-i18next';
import { useParams, Link } from 'react-router-dom';
import { useNameTrend } from '../hooks/useNameTrend';

export default function NameDetailPage() {
  const { t } = useTranslation(['pages', 'common']);
  const { name } = useParams<{ name: string }>();
  const { data, isLoading, error } = useNameTrend({ name: name || '' });

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      <div className="container mx-auto px-4 py-8 md:py-12">
        {/* Back Link */}
        <div className="mb-6 animate-fade-in">
          <Link
            to="/names"
            className="inline-flex items-center gap-2 text-primary-600 hover:text-primary-700 font-medium transition-all duration-200 hover:gap-3 group"
          >
            <svg className="w-5 h-5 transition-transform group-hover:-translate-x-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            {t('pages:nameDetail.backButton')}
          </Link>
        </div>

        {/* Hero Section */}
        <div className="mb-12 animate-slide-up">
          <div className="relative overflow-hidden bg-gradient-to-r from-primary-600 via-secondary-600 to-accent-600 rounded-3xl shadow-strong p-12 md:p-16">
            {/* Decorative Elements */}
            <div className="absolute top-0 right-0 -mt-4 -mr-4 w-40 h-40 bg-white/10 rounded-full blur-3xl"></div>
            <div className="absolute bottom-0 left-0 -mb-4 -ml-4 w-40 h-40 bg-white/10 rounded-full blur-3xl"></div>
            
            <div className="relative text-center">
              <div className="inline-flex items-center justify-center w-16 h-16 bg-white/20 backdrop-blur-sm rounded-2xl mb-6 shadow-lg">
                <span className="text-3xl">✨</span>
              </div>
              <h1 className="text-5xl md:text-7xl font-bold text-white mb-4 tracking-tight">
                {name}
              </h1>
              <p className="text-xl md:text-2xl text-white/90 max-w-2xl mx-auto">
                {t('pages:nameDetail.title')}
              </p>
            </div>
          </div>
        </div>

        {/* Loading State */}
        {isLoading && (
          <div className="bg-white rounded-2xl shadow-medium p-12 border border-gray-100 animate-pulse">
            <div className="flex items-center justify-center py-12">
              <div className="text-center">
                <div className="relative w-16 h-16 mx-auto mb-6">
                  <div className="absolute inset-0 border-4 border-primary-200 border-t-primary-600 rounded-full animate-spin"></div>
                  <div className="absolute inset-2 bg-primary-50 rounded-full"></div>
                </div>
                <p className="text-gray-600 font-medium">
                  {t('pages:nameDetail.loading')}
                </p>
              </div>
            </div>
          </div>
        )}

        {/* Error State */}
        {error && (
          <div className="bg-white rounded-2xl shadow-medium p-12 border border-red-100 animate-slide-up">
            <div className="flex items-center justify-center py-12">
              <div className="text-center">
                <div className="inline-flex items-center justify-center w-16 h-16 bg-red-100 rounded-full mb-6">
                  <svg className="w-8 h-8 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <p className="text-red-600 font-semibold text-lg mb-2">
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
            {/* Summary Statistics Cards */}
            <div className="mb-8 animate-slide-up">
              <div className="flex items-center gap-3 mb-6">
                <svg className="w-6 h-6 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
                <h2 className="text-2xl md:text-3xl font-bold text-gray-900">
                  {t('pages:nameDetail.statistics')}
                </h2>
              </div>

              <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
                {/* Total Count Card */}
                <div className="group bg-white rounded-2xl shadow-medium hover:shadow-strong transition-all duration-300 p-6 border border-gray-100 hover:border-primary-200">
                  <div className="flex items-center gap-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-primary-100 to-primary-200 rounded-xl flex items-center justify-center shadow-sm group-hover:scale-110 transition-transform">
                      <svg className="w-6 h-6 text-primary-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                      </svg>
                    </div>
                    <div className="flex-1">
                      <div className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-1">
                        Total Count
                      </div>
                      <div className="text-2xl font-bold text-gray-900">
                        {data.summary.total_count.toLocaleString()}
                      </div>
                    </div>
                  </div>
                </div>

                {/* Female Count Card */}
                <div className="group bg-white rounded-2xl shadow-medium hover:shadow-strong transition-all duration-300 p-6 border border-gray-100 hover:border-pink-200">
                  <div className="flex items-center gap-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-pink-100 to-pink-200 rounded-xl flex items-center justify-center shadow-sm group-hover:scale-110 transition-transform">
                      <span className="text-2xl">♀</span>
                    </div>
                    <div className="flex-1">
                      <div className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-1">
                        Female Count
                      </div>
                      <div className="text-2xl font-bold text-pink-900">
                        {data.summary.female_count.toLocaleString()}
                      </div>
                    </div>
                  </div>
                </div>

                {/* Male Count Card */}
                <div className="group bg-white rounded-2xl shadow-medium hover:shadow-strong transition-all duration-300 p-6 border border-gray-100 hover:border-blue-200">
                  <div className="flex items-center gap-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-blue-100 to-blue-200 rounded-xl flex items-center justify-center shadow-sm group-hover:scale-110 transition-transform">
                      <span className="text-2xl">♂</span>
                    </div>
                    <div className="flex-1">
                      <div className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-1">
                        Male Count
                      </div>
                      <div className="text-2xl font-bold text-blue-900">
                        {data.summary.male_count.toLocaleString()}
                      </div>
                    </div>
                  </div>
                </div>

                {/* Gender Balance Card */}
                <div className="group bg-white rounded-2xl shadow-medium hover:shadow-strong transition-all duration-300 p-6 border border-gray-100 hover:border-purple-200">
                  <div className="flex items-center gap-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-purple-100 to-purple-200 rounded-xl flex items-center justify-center shadow-sm group-hover:scale-110 transition-transform">
                      <svg className="w-6 h-6 text-purple-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 6l3 1m0 0l-3 9a5.002 5.002 0 006.001 0M6 7l3 9M6 7l6-2m6 2l3-1m-3 1l-3 9a5.002 5.002 0 006.001 0M18 7l3 9m-3-9l-6-2m0-2v2m0 16V5m0 16H9m3 0h3" />
                      </svg>
                    </div>
                    <div className="flex-1">
                      <div className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-1">
                        Gender Balance
                      </div>
                      <div className="text-2xl font-bold text-purple-900">
                        {data.summary.gender_balance !== null 
                          ? data.summary.gender_balance.toFixed(1)
                          : 'N/A'
                        }
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Additional Info */}
              <div className="mt-6 bg-gradient-to-r from-gray-50 to-gray-100 rounded-xl p-6 border border-gray-200">
                <div className="flex flex-wrap gap-6 text-sm">
                  <div className="flex items-center gap-2">
                    <svg className="w-5 h-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                    <span className="font-semibold text-gray-700">Years:</span>
                    <span className="text-gray-900">{data.summary.name_start} - {data.summary.name_end}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <svg className="w-5 h-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span className="font-semibold text-gray-700">Countries:</span>
                    <span className="text-gray-900">{data.summary.countries.join(', ')}</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Time Series Placeholder */}
            <div className="mb-8 animate-slide-up animation-delay-100">
              <div className="bg-white rounded-2xl shadow-medium p-8 border border-gray-100">
                <div className="flex items-center gap-3 mb-6">
                  <svg className="w-6 h-6 text-secondary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" />
                  </svg>
                  <h2 className="text-2xl md:text-3xl font-bold text-gray-900">
                    {t('pages:nameDetail.trends')}
                  </h2>
                </div>
                <div className="flex flex-col items-center justify-center py-20 border-2 border-dashed border-gray-300 rounded-xl bg-gray-50">
                  <div className="w-16 h-16 bg-gradient-to-br from-secondary-100 to-secondary-200 rounded-full flex items-center justify-center mb-4 shadow-md">
                    <svg className="w-8 h-8 text-secondary-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" />
                    </svg>
                  </div>
                  <p className="text-gray-500 text-lg font-medium mb-2">
                    {t('pages:nameDetail.chartsComingSoon')}
                  </p>
                  <p className="text-sm text-gray-400">
                    Data available for {data.time_series.length} years
                  </p>
                </div>
              </div>
            </div>

            {/* By Country Placeholder */}
            <div className="animate-slide-up animation-delay-200">
              <div className="bg-white rounded-2xl shadow-medium p-8 border border-gray-100">
                <div className="flex items-center gap-3 mb-6">
                  <svg className="w-6 h-6 text-accent-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <h2 className="text-2xl md:text-3xl font-bold text-gray-900">
                    {t('pages:nameDetail.distribution')}
                  </h2>
                </div>
                <div className="flex flex-col items-center justify-center py-20 border-2 border-dashed border-gray-300 rounded-xl bg-gray-50">
                  <div className="w-16 h-16 bg-gradient-to-br from-accent-100 to-accent-200 rounded-full flex items-center justify-center mb-4 shadow-md">
                    <svg className="w-8 h-8 text-accent-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
                    </svg>
                  </div>
                  <p className="text-gray-500 text-lg font-medium mb-2">
                    {t('pages:nameDetail.chartsComingSoon')}
                  </p>
                  <p className="text-sm text-gray-400">
                    Data available for {data.by_country.length} countries
                  </p>
                </div>
              </div>
            </div>
          </>
        )}

        {/* No Data State */}
        {!isLoading && !error && !data && (
          <div className="bg-white rounded-2xl shadow-medium p-12 border border-gray-100 animate-slide-up">
            <div className="flex items-center justify-center py-12">
              <div className="text-center">
                <div className="inline-flex items-center justify-center w-20 h-20 bg-gray-100 rounded-full mb-6">
                  <svg className="w-10 h-10 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                </div>
                <p className="text-gray-600 text-lg">
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