/**
 * PopularityFilterTrio Component
 * 
 * Three mutually exclusive popularity filters: Min Count, Top N, Coverage Percent.
 * Only one can be active at a time (the "driver").
 * When user edits one, it becomes the driver.
 * The other two display derived values from API response.
 */

import { useTranslation } from 'react-i18next';

type PopularityDriver = 'min_count' | 'top_n' | 'coverage_percent' | null;

interface PopularityFilterTrioProps {
  minCount: number | null;
  topN: number | null;
  coveragePercent: number | null;
  activeDriver: PopularityDriver;
  onMinCountChange: (value: number | null) => void;
  onTopNChange: (value: number | null) => void;
  onCoveragePercentChange: (value: number | null) => void;
}

export default function PopularityFilterTrio({
  minCount,
  topN,
  coveragePercent,
  activeDriver,
  onMinCountChange,
  onTopNChange,
  onCoveragePercentChange,
}: PopularityFilterTrioProps) {
  const { t } = useTranslation('filters');

  const handleMinCountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value === '' ? null : parseInt(e.target.value, 10);
    onMinCountChange(value);
  };

  const handleTopNChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value === '' ? null : parseInt(e.target.value, 10);
    onTopNChange(value);
  };

  const handleCoveragePercentChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value === '' ? null : parseFloat(e.target.value);
    onCoveragePercentChange(value);
  };

  return (
    <div className="space-y-2">
      <label className="flex items-center gap-2 text-sm font-semibold text-gray-700">
        <svg className="w-4 h-4 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
        </svg>
        {t('popularity.label')}
      </label>
      
      <div className="grid grid-cols-3 gap-3">
        {/* Min Count */}
        <div className="relative">
          <input
            type="number"
            placeholder={t('popularity.minCount.placeholder')}
            value={minCount ?? ''}
            onChange={handleMinCountChange}
            min="0"
            className={`w-full px-3 py-2.5 border rounded-lg focus:ring-2 transition-all duration-200 ${
              activeDriver === 'min_count'
                ? 'border-primary-500 bg-primary-50 ring-2 ring-primary-200 font-semibold'
                : activeDriver
                ? 'border-gray-200 bg-gray-50 text-gray-600 cursor-default'
                : 'border-gray-300 bg-white hover:border-gray-400'
            }`}
            readOnly={activeDriver !== null && activeDriver !== 'min_count'}
          />
          <label className="block text-xs font-medium text-gray-600 mt-1">
            {t('popularity.minCount.label')}
            {activeDriver === 'min_count' && (
              <span className="ml-1 text-primary-600">●</span>
            )}
          </label>
        </div>

        {/* Top N */}
        <div className="relative">
          <input
            type="number"
            placeholder={t('popularity.topN.placeholder')}
            value={topN ?? ''}
            onChange={handleTopNChange}
            min="1"
            className={`w-full px-3 py-2.5 border rounded-lg focus:ring-2 transition-all duration-200 ${
              activeDriver === 'top_n'
                ? 'border-secondary-500 bg-secondary-50 ring-2 ring-secondary-200 font-semibold'
                : activeDriver
                ? 'border-gray-200 bg-gray-50 text-gray-600 cursor-default'
                : 'border-gray-300 bg-white hover:border-gray-400'
            }`}
            readOnly={activeDriver !== null && activeDriver !== 'top_n'}
          />
          <label className="block text-xs font-medium text-gray-600 mt-1">
            {t('popularity.topN.label')}
            {activeDriver === 'top_n' && (
              <span className="ml-1 text-secondary-600">●</span>
            )}
          </label>
        </div>

        {/* Coverage Percent */}
        <div className="relative">
          <input
            type="number"
            placeholder={t('popularity.coveragePercent.placeholder')}
            value={coveragePercent ?? ''}
            onChange={handleCoveragePercentChange}
            min="0"
            max="100"
            step="0.1"
            className={`w-full px-3 py-2.5 border rounded-lg focus:ring-2 transition-all duration-200 ${
              activeDriver === 'coverage_percent'
                ? 'border-accent-500 bg-accent-50 ring-2 ring-accent-200 font-semibold'
                : activeDriver
                ? 'border-gray-200 bg-gray-50 text-gray-600 cursor-default'
                : 'border-gray-300 bg-white hover:border-gray-400'
            }`}
            readOnly={activeDriver !== null && activeDriver !== 'coverage_percent'}
          />
          <label className="block text-xs font-medium text-gray-600 mt-1">
            {t('popularity.coveragePercent.label')}
            {activeDriver === 'coverage_percent' && (
              <span className="ml-1 text-accent-600">●</span>
            )}
          </label>
        </div>
      </div>

      {activeDriver && (
        <p className="text-xs text-gray-500 flex items-center gap-1 mt-2">
          <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Active filter: {activeDriver.replace('_', ' ')} • Other values are derived from results
        </p>
      )}
    </div>
  );
}