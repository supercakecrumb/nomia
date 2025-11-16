/**
 * FilterBar Component
 * 
 * Filter controls for the Names Explorer page.
 * Phase 1: UI only, without full state management.
 */

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useMetaYears } from '../../hooks/useMetaYears';
import { useCountries } from '../../hooks/useCountries';

interface FilterBarProps {
  /** Callback when filters change (not implemented in Phase 1) */
  onFiltersChange?: (filters: any) => void;
}

export default function FilterBar({ onFiltersChange }: FilterBarProps) {
  const { t } = useTranslation('filters');
  const { data: metaYears } = useMetaYears();
  const { data: countriesData } = useCountries();

  // Basic local state for Phase 1 (not yet connected to parent)
  const [yearMin, setYearMin] = useState('');
  const [yearMax, setYearMax] = useState('');
  const [selectedCountries, setSelectedCountries] = useState<string[]>([]);
  const [genderBalanceMin, setGenderBalanceMin] = useState('');
  const [genderBalanceMax, setGenderBalanceMax] = useState('');
  const [nameSearch, setNameSearch] = useState('');
  const [sortBy, setSortBy] = useState('rank');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [isExpanded, setIsExpanded] = useState(true);

  const countries = countriesData?.countries || [];

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
        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="lg:hidden px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900 bg-white rounded-lg hover:bg-gray-50 transition-colors shadow-sm"
        >
          {isExpanded ? (
            <span className="flex items-center gap-2">
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
              </svg>
              Hide
            </span>
          ) : (
            <span className="flex items-center gap-2">
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
              Show
            </span>
          )}
        </button>
      </div>

      {/* Filters Grid */}
      <div className={`p-6 ${isExpanded ? 'block' : 'hidden lg:block'}`}>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {/* Year Range Filter */}
          <div className="space-y-2">
            <label className="flex items-center gap-2 text-sm font-semibold text-gray-700">
              <svg className="w-4 h-4 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              {t('yearRange.label')}
            </label>
            <div className="flex gap-2">
              <input
                type="number"
                placeholder={t('yearRange.from')}
                value={yearMin}
                onChange={(e) => setYearMin(e.target.value)}
                min={metaYears?.min_year}
                max={metaYears?.max_year}
                className="flex-1 px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 hover:border-gray-400 bg-white placeholder:text-gray-400"
              />
              <span className="flex items-center text-gray-400 font-medium">—</span>
              <input
                type="number"
                placeholder={t('yearRange.to')}
                value={yearMax}
                onChange={(e) => setYearMax(e.target.value)}
                min={metaYears?.min_year}
                max={metaYears?.max_year}
                className="flex-1 px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 hover:border-gray-400 bg-white placeholder:text-gray-400"
              />
            </div>
            {metaYears && (
              <p className="text-xs text-gray-500 flex items-center gap-1">
                <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Available: {metaYears.min_year} - {metaYears.max_year}
              </p>
            )}
          </div>

          {/* Countries Multi-Select */}
          <div className="space-y-2">
            <label className="flex items-center gap-2 text-sm font-semibold text-gray-700">
              <svg className="w-4 h-4 text-secondary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {t('countries.label')}
            </label>
            <select
              multiple
              value={selectedCountries}
              onChange={(e) =>
                setSelectedCountries(
                  Array.from(e.target.selectedOptions, (option) => option.value)
                )
              }
              className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-secondary-500 focus:border-transparent transition-all duration-200 hover:border-gray-400 bg-white"
              size={4}
            >
              {countries.map((country) => (
                <option key={country.code} value={country.code} className="py-1">
                  {country.code} - {country.name}
                </option>
              ))}
            </select>
            <p className="text-xs text-gray-500 flex items-center gap-1">
              <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Hold Ctrl/Cmd to select multiple
            </p>
          </div>

          {/* Gender Balance Filter */}
          <div className="space-y-2">
            <label className="flex items-center gap-2 text-sm font-semibold text-gray-700">
              <svg className="w-4 h-4 text-accent-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 6l3 1m0 0l-3 9a5.002 5.002 0 006.001 0M6 7l3 9M6 7l6-2m6 2l3-1m-3 1l-3 9a5.002 5.002 0 006.001 0M18 7l3 9m-3-9l-6-2m0-2v2m0 16V5m0 16H9m3 0h3" />
              </svg>
              {t('genderBalance.label')}
            </label>
            <div className="space-y-2">
              <div className="flex gap-2">
                <input
                  type="number"
                  placeholder="Min (0=♀)"
                  value={genderBalanceMin}
                  onChange={(e) => setGenderBalanceMin(e.target.value)}
                  min="0"
                  max="100"
                  className="flex-1 px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent-500 focus:border-transparent transition-all duration-200 hover:border-gray-400 bg-white placeholder:text-gray-400"
                />
                <input
                  type="number"
                  placeholder="Max (100=♂)"
                  value={genderBalanceMax}
                  onChange={(e) => setGenderBalanceMax(e.target.value)}
                  min="0"
                  max="100"
                  className="flex-1 px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent-500 focus:border-transparent transition-all duration-200 hover:border-gray-400 bg-white placeholder:text-gray-400"
                />
              </div>
              <p className="text-xs text-gray-500">
                0 = {t('genderBalance.female')}, 100 = {t('genderBalance.male')}
              </p>
            </div>
          </div>

          {/* Name Search */}
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
              onChange={(e) => setNameSearch(e.target.value)}
              className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 hover:border-gray-400 bg-white placeholder:text-gray-400"
            />
            <p className="text-xs text-gray-500 flex items-center gap-1">
              <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Supports glob patterns: *, ?
            </p>
          </div>

          {/* Sort By */}
          <div className="space-y-2">
            <label className="flex items-center gap-2 text-sm font-semibold text-gray-700">
              <svg className="w-4 h-4 text-secondary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4h13M3 8h9m-9 4h9m5-4v12m0 0l-4-4m4 4l4-4" />
              </svg>
              {t('sort.label')}
            </label>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-secondary-500 focus:border-transparent transition-all duration-200 hover:border-gray-400 bg-white"
            >
              <option value="rank">{t('sort.options.rank')}</option>
              <option value="name">{t('sort.options.name')}</option>
              <option value="total_count">{t('sort.options.totalCount')}</option>
              <option value="gender_balance">{t('sort.options.genderBalance')}</option>
            </select>
          </div>

          {/* Sort Order */}
          <div className="space-y-2">
            <label className="flex items-center gap-2 text-sm font-semibold text-gray-700">
              <svg className="w-4 h-4 text-accent-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
              </svg>
              {t('sort.order')}
            </label>
            <div className="flex gap-2">
              <button
                onClick={() => setSortOrder('asc')}
                className={`flex-1 px-4 py-2.5 rounded-lg font-medium transition-all duration-200 ${
                  sortOrder === 'asc'
                    ? 'bg-gradient-to-r from-primary-600 to-primary-700 text-white shadow-md'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                {t('sort.ascending')}
              </button>
              <button
                onClick={() => setSortOrder('desc')}
                className={`flex-1 px-4 py-2.5 rounded-lg font-medium transition-all duration-200 ${
                  sortOrder === 'desc'
                    ? 'bg-gradient-to-r from-primary-600 to-primary-700 text-white shadow-md'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                {t('sort.descending')}
              </button>
            </div>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="mt-8 pt-6 border-t border-gray-200 flex flex-col sm:flex-row gap-3">
          <button
            className="flex-1 sm:flex-none px-6 py-3 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 font-medium transition-all duration-200 hover:shadow-md active:scale-95 flex items-center justify-center gap-2"
            onClick={() => {
              setYearMin('');
              setYearMax('');
              setSelectedCountries([]);
              setGenderBalanceMin('');
              setGenderBalanceMax('');
              setNameSearch('');
              setSortBy('rank');
              setSortOrder('asc');
            }}
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            {t('actions.reset')}
          </button>
          <button
            className="flex-1 px-6 py-3 bg-gradient-to-r from-primary-600 to-primary-700 text-white rounded-xl hover:from-primary-700 hover:to-primary-800 font-medium transition-all duration-200 shadow-md hover:shadow-lg active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:shadow-md flex items-center justify-center gap-2"
            disabled
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
            {t('actions.apply')} (Coming soon)
          </button>
        </div>
      </div>
    </div>
  );
}