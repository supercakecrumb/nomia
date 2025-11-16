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
    <div className="bg-white rounded-lg shadow-md">
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b border-gray-200">
        <h2 className="text-xl font-semibold text-gray-900">{t('title')}</h2>
        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="lg:hidden px-3 py-1 text-sm text-gray-600 hover:text-gray-900"
        >
          {isExpanded ? '▲ Hide' : '▼ Show'}
        </button>
      </div>

      {/* Filters Grid */}
      <div className={`p-6 ${isExpanded ? 'block' : 'hidden lg:block'}`}>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {/* Year Range Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
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
                className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <span className="flex items-center text-gray-500">—</span>
              <input
                type="number"
                placeholder={t('yearRange.to')}
                value={yearMax}
                onChange={(e) => setYearMax(e.target.value)}
                min={metaYears?.min_year}
                max={metaYears?.max_year}
                className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            {metaYears && (
              <p className="mt-1 text-xs text-gray-500">
                Available: {metaYears.min_year} - {metaYears.max_year}
              </p>
            )}
          </div>

          {/* Countries Multi-Select */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
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
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              size={4}
            >
              {countries.map((country) => (
                <option key={country.code} value={country.code}>
                  {country.code} - {country.name}
                </option>
              ))}
            </select>
            <p className="mt-1 text-xs text-gray-500">
              Hold Ctrl/Cmd to select multiple
            </p>
          </div>

          {/* Gender Balance Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
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
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                <input
                  type="number"
                  placeholder="Max (100=♂)"
                  value={genderBalanceMax}
                  onChange={(e) => setGenderBalanceMax(e.target.value)}
                  min="0"
                  max="100"
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <p className="text-xs text-gray-500">
                0 = {t('genderBalance.female')}, 100 = {t('genderBalance.male')}
              </p>
            </div>
          </div>

          {/* Name Search */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {t('nameSearch.label')}
            </label>
            <input
              type="text"
              placeholder={t('nameSearch.placeholder')}
              value={nameSearch}
              onChange={(e) => setNameSearch(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            <p className="mt-1 text-xs text-gray-500">
              Supports glob patterns: *, ?
            </p>
          </div>

          {/* Sort By */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {t('sort.label')}
            </label>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="rank">{t('sort.options.rank')}</option>
              <option value="name">{t('sort.options.name')}</option>
              <option value="total_count">{t('sort.options.totalCount')}</option>
              <option value="gender_balance">{t('sort.options.genderBalance')}</option>
            </select>
          </div>

          {/* Sort Order */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {t('sort.order')}
            </label>
            <div className="flex gap-2">
              <button
                onClick={() => setSortOrder('asc')}
                className={`flex-1 px-3 py-2 rounded-md font-medium ${
                  sortOrder === 'asc'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                {t('sort.ascending')}
              </button>
              <button
                onClick={() => setSortOrder('desc')}
                className={`flex-1 px-3 py-2 rounded-md font-medium ${
                  sortOrder === 'desc'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                {t('sort.descending')}
              </button>
            </div>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="mt-6 flex gap-3">
          <button
            className="px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 font-medium"
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
            {t('actions.reset')}
          </button>
          <button
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 font-medium disabled:bg-gray-300 disabled:cursor-not-allowed"
            disabled
          >
            {t('actions.apply')} (Coming soon)
          </button>
        </div>
      </div>
    </div>
  );
}