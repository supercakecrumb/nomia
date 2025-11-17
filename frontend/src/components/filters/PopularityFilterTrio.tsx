/**
 * PopularityFilterTrio Component
 *
 * Three mutually exclusive popularity filters: Min Count, Top N, Coverage Percent.
 * Only one can be active at a time (the "driver").
 * When user edits one, it becomes the driver.
 * The other two display derived values from API response.
 * Each filter includes both a slider and a numeric input that sync with each other.
 */

import { useTranslation } from 'react-i18next';
import Slider from 'rc-slider';
import 'rc-slider/assets/index.css';

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

  const handleMinCountFocus = () => {
    // When this field is focused, clear the other two if they have values
    if (topN !== null || coveragePercent !== null) {
      onTopNChange(null);
      onCoveragePercentChange(null);
    }
  };

  const handleTopNFocus = () => {
    // When this field is focused, clear the other two if they have values
    if (minCount !== null || coveragePercent !== null) {
      onMinCountChange(null);
      onCoveragePercentChange(null);
    }
  };

  const handleCoveragePercentFocus = () => {
    // When this field is focused, clear the other two if they have values
    if (minCount !== null || topN !== null) {
      onMinCountChange(null);
      onTopNChange(null);
    }
  };

  const handleMinCountInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value === '' ? null : parseInt(e.target.value, 10);
    onMinCountChange(value);
  };

  const handleTopNInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value === '' ? null : parseInt(e.target.value, 10);
    onTopNChange(value);
  };

  const handleCoveragePercentInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value === '' ? null : parseFloat(e.target.value);
    onCoveragePercentChange(value);
  };

  const handleMinCountSliderChange = (value: number | number[]) => {
    const numValue = Array.isArray(value) ? value[0] : value;
    onMinCountChange(numValue);
  };

  const handleTopNSliderChange = (value: number | number[]) => {
    const numValue = Array.isArray(value) ? value[0] : value;
    onTopNChange(numValue);
  };

  const handleCoveragePercentSliderChange = (value: number | number[]) => {
    const numValue = Array.isArray(value) ? value[0] : value;
    onCoveragePercentChange(numValue);
  };

  return (
    <div className="space-y-2">
      <label className="flex items-center gap-2 text-sm font-semibold text-gray-700">
        <svg className="w-4 h-4 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
        </svg>
        {t('popularity.label')}
      </label>
      
      <div className="grid grid-cols-3 gap-4">
        {/* Min Count */}
        <div className="space-y-2">
          <label className="block text-xs font-medium text-gray-600">
            {t('popularity.minCount.label')}
            {activeDriver === 'min_count' && (
              <span className="ml-1 text-primary-600">●</span>
            )}
          </label>
          <input
            type="number"
            placeholder={t('popularity.minCount.placeholder')}
            value={minCount ?? ''}
            onChange={handleMinCountInputChange}
            onFocus={handleMinCountFocus}
            min="0"
            className={`w-full px-3 py-2 border rounded-lg focus:ring-2 transition-all duration-200 text-sm ${
              activeDriver === 'min_count'
                ? 'border-primary-500 bg-primary-50 ring-2 ring-primary-200 font-semibold'
                : activeDriver
                ? 'border-gray-200 bg-gray-50 text-gray-600'
                : 'border-gray-300 bg-white hover:border-gray-400'
            }`}
          />
          <div className="px-2">
            <Slider
              min={0}
              max={10000}
              value={minCount ?? 0}
              onChange={handleMinCountSliderChange}
              onFocus={handleMinCountFocus}
              disabled={activeDriver !== null && activeDriver !== 'min_count'}
              railStyle={{ backgroundColor: '#e5e7eb', height: 6 }}
              trackStyle={{ backgroundColor: activeDriver === 'min_count' ? '#3b82f6' : '#94a3b8', height: 6 }}
              handleStyle={{
                borderColor: activeDriver === 'min_count' ? '#3b82f6' : '#94a3b8',
                backgroundColor: '#ffffff',
                opacity: 1,
                width: 16,
                height: 16,
                marginTop: -5,
              }}
            />
          </div>
        </div>

        {/* Top N */}
        <div className="space-y-2">
          <label className="block text-xs font-medium text-gray-600">
            {t('popularity.topN.label')}
            {activeDriver === 'top_n' && (
              <span className="ml-1 text-secondary-600">●</span>
            )}
          </label>
          <input
            type="number"
            placeholder={t('popularity.topN.placeholder')}
            value={topN ?? ''}
            onChange={handleTopNInputChange}
            onFocus={handleTopNFocus}
            min="1"
            className={`w-full px-3 py-2 border rounded-lg focus:ring-2 transition-all duration-200 text-sm ${
              activeDriver === 'top_n'
                ? 'border-secondary-500 bg-secondary-50 ring-2 ring-secondary-200 font-semibold'
                : activeDriver
                ? 'border-gray-200 bg-gray-50 text-gray-600'
                : 'border-gray-300 bg-white hover:border-gray-400'
            }`}
          />
          <div className="px-2">
            <Slider
              min={1}
              max={1000}
              value={topN ?? 1}
              onChange={handleTopNSliderChange}
              onFocus={handleTopNFocus}
              disabled={activeDriver !== null && activeDriver !== 'top_n'}
              railStyle={{ backgroundColor: '#e5e7eb', height: 6 }}
              trackStyle={{ backgroundColor: activeDriver === 'top_n' ? '#14b8a6' : '#94a3b8', height: 6 }}
              handleStyle={{
                borderColor: activeDriver === 'top_n' ? '#14b8a6' : '#94a3b8',
                backgroundColor: '#ffffff',
                opacity: 1,
                width: 16,
                height: 16,
                marginTop: -5,
              }}
            />
          </div>
        </div>

        {/* Coverage Percent */}
        <div className="space-y-2">
          <label className="block text-xs font-medium text-gray-600">
            {t('popularity.coveragePercent.label')}
            {activeDriver === 'coverage_percent' && (
              <span className="ml-1 text-accent-600">●</span>
            )}
          </label>
          <input
            type="number"
            placeholder={t('popularity.coveragePercent.placeholder')}
            value={coveragePercent ?? ''}
            onChange={handleCoveragePercentInputChange}
            onFocus={handleCoveragePercentFocus}
            min="0"
            max="100"
            step="0.1"
            className={`w-full px-3 py-2 border rounded-lg focus:ring-2 transition-all duration-200 text-sm ${
              activeDriver === 'coverage_percent'
                ? 'border-accent-500 bg-accent-50 ring-2 ring-accent-200 font-semibold'
                : activeDriver
                ? 'border-gray-200 bg-gray-50 text-gray-600'
                : 'border-gray-300 bg-white hover:border-gray-400'
            }`}
          />
          <div className="px-2">
            <Slider
              min={0}
              max={100}
              step={0.1}
              value={coveragePercent ?? 0}
              onChange={handleCoveragePercentSliderChange}
              onFocus={handleCoveragePercentFocus}
              disabled={activeDriver !== null && activeDriver !== 'coverage_percent'}
              railStyle={{ backgroundColor: '#e5e7eb', height: 6 }}
              trackStyle={{ backgroundColor: activeDriver === 'coverage_percent' ? '#f59e0b' : '#94a3b8', height: 6 }}
              handleStyle={{
                borderColor: activeDriver === 'coverage_percent' ? '#f59e0b' : '#94a3b8',
                backgroundColor: '#ffffff',
                opacity: 1,
                width: 16,
                height: 16,
                marginTop: -5,
              }}
            />
          </div>
        </div>
      </div>

      {activeDriver && (
        <p className="text-xs text-gray-500 flex items-center gap-1 mt-2">
          <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          {t('filters:popularity.activeDescription', { driver: activeDriver.replace('_', ' ') })}
        </p>
      )}
    </div>
  );
}