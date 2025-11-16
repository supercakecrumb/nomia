/**
 * FilterBar Component
 * 
 * Comprehensive filter controls for the Names Explorer page.
 * Includes year range slider, countries dropdown, gender balance slider,
 * popularity trio, and name search.
 */

import { useTranslation } from 'react-i18next';
import { useMetaYears } from '../../hooks/useMetaYears';
import { useCountries } from '../../hooks/useCountries';
import YearRangeSlider from './YearRangeSlider';
import CountriesDropdown from './CountriesDropdown';
import GenderBalanceSlider from './GenderBalanceSlider';
// TODO: Phase 2 - Uncomment when popularity filter is re-enabled
// import PopularityFilterTrio from './PopularityFilterTrio';

interface FilterBarProps {
  // Year range
  yearMin: number;
  yearMax: number;
  onYearRangeChange: (min: number, max: number) => void;
  
  // Countries
  selectedCountries: string[];
  onCountriesChange: (countries: string[]) => void;
  
  // Gender balance
  genderBalanceMin: number;
  genderBalanceMax: number;
  onGenderBalanceChange: (min: number, max: number) => void;
  
  // TODO: Phase 2 - Add back popularity trio filters
  // Requires backend optimization to handle <500ms response times
  // minCount: number | null;
  // topN: number | null;
  // coveragePercent: number | null;
  // popularityDriver: 'min_count' | 'top_n' | 'coverage_percent' | null;
  // onMinCountChange: (value: number | null) => void;
  // onTopNChange: (value: number | null) => void;
  // onCoveragePercentChange: (value: number | null) => void;
  
  // TODO: Phase 2 - Add name search with glob pattern support
  // Backend needs to implement efficient name matching
  // nameSearch: string;
  // onNameSearchChange: (value: string) => void;
  
  // Actions
  onReset: () => void;
}

export default function FilterBar({
  yearMin,
  yearMax,
  onYearRangeChange,
  selectedCountries,
  onCountriesChange,
  genderBalanceMin,
  genderBalanceMax,
  onGenderBalanceChange,
  // minCount,
  // topN,
  // coveragePercent,
  // popularityDriver,
  // onMinCountChange,
  // onTopNChange,
  // onCoveragePercentChange,
  // nameSearch,
  // onNameSearchChange,
  onReset,
}: FilterBarProps) {
  const { t } = useTranslation('filters');
  const { data: metaYears } = useMetaYears();
  const { data: countriesData } = useCountries();

  const countries = countriesData?.countries || [];
  
  // Use API values when available, permissive fallbacks otherwise
  // Using 0 and 9999 ensures no artificial restrictions until API loads
  const minYear = metaYears?.min_year ?? 0;
  const maxYear = metaYears?.max_year ?? 9999;

  return (
    <div className="bg-white rounded-2xl shadow-medium hover:shadow-strong transition-shadow duration-300 border border-gray-100">
      {/* Header */}
      <div className="flex items-center justify-between p-6 border-b border-gray-200 bg-gradient-to-r from-primary-50 to-secondary-50 rounded-t-2xl">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-gradient-to-br from-primary-600 to-secondary-600 rounded-xl flex items-center justify-center shadow-md">
            <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
            </svg>
          </div>
          <h2 className="text-xl font-bold text-gray-900">{t('title')}</h2>
        </div>
      </div>

      {/* Filters */}
      <div className="p-6 space-y-6">
        {/* Year Range Slider - Full Width */}
        <div className="pb-4 border-b border-gray-100">
          <YearRangeSlider
            min={minYear}
            max={maxYear}
            value={[yearMin, yearMax]}
            onChange={([min, max]) => onYearRangeChange(min, max)}
          />
        </div>

        {/* Countries and Gender Balance - Side by Side */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 pb-4 border-b border-gray-100">
          <CountriesDropdown
            countries={countries}
            selectedCountries={selectedCountries}
            onChange={onCountriesChange}
          />
          
          <GenderBalanceSlider
            value={[genderBalanceMin, genderBalanceMax]}
            onChange={([min, max]) => onGenderBalanceChange(min, max)}
          />
        </div>

        {/* TODO: Phase 2 - Popularity Trio Filters
            Uncomment when backend is optimized (<500ms response times)
            
        <div className="pb-4 border-b border-gray-100">
          <PopularityFilterTrio
            minCount={minCount}
            topN={topN}
            coveragePercent={coveragePercent}
            activeDriver={popularityDriver}
            onMinCountChange={onMinCountChange}
            onTopNChange={onTopNChange}
            onCoveragePercentChange={onCoveragePercentChange}
          />
        </div>
        */}

        {/* TODO: Phase 2 - Name Search
            Uncomment when backend implements efficient glob pattern search
            
        <div className="space-y-2">
          <label className="flex items-center gap-2 text-sm font-semibold text-gray-700">
            <svg className="w-4 h-4 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            {t('nameSearch.label')}
          </label>
          <input
            type="text"
            placeholder={t('nameSearch.placeholder')}
            value={nameSearch}
            onChange={(e) => onNameSearchChange(e.target.value)}
            className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 hover:border-gray-400 bg-white placeholder:text-gray-400"
          />
          <p className="text-xs text-gray-500 flex items-center gap-1">
            <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Supports glob patterns: *, ?
          </p>
        </div>
        */}

        {/* Action Buttons */}
        <div className="pt-4 flex gap-3">
          <button
            onClick={onReset}
            className="px-6 py-3 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 font-medium transition-all duration-200 hover:shadow-md active:scale-95 flex items-center justify-center gap-2"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            {t('actions.reset')}
          </button>
        </div>
      </div>
    </div>
  );
}