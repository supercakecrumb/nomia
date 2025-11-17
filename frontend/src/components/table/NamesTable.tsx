/**
 * NamesTable Component
 * 
 * Displays a table of gender-neutral names with various statistics.
 * Supports loading, error, and empty states.
 */

import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import type { NameEntry } from '../../types/api';
import { formatNumber } from '../../utils/formatters';
import GenderBalanceBar from './GenderBalanceBar';
import { getNameDetailPath } from '../../utils/navigation';

interface NamesTableProps {
  /** Array of name entries to display */
  names: NameEntry[];
  /** Loading state */
  isLoading?: boolean;
  /** Error state */
  error?: Error | null;
  /** Current sort column */
  sortBy?: 'name' | 'total_count' | 'gender_balance' | 'rank' | null;
  /** Current sort order */
  sortOrder?: 'asc' | 'desc' | null;
  /** Callback when sort changes */
  onSortChange?: (sortBy: 'name' | 'total_count' | 'gender_balance' | 'rank' | null, sortOrder: 'asc' | 'desc' | null) => void;
}

export default function NamesTable({
  names,
  isLoading,
  error,
  sortBy = null,
  sortOrder = null,
  onSortChange,
}: NamesTableProps) {
  const { t } = useTranslation(['common', 'pages', 'filters']);

  const handleSort = (column: 'name' | 'total_count' | 'gender_balance' | 'rank') => {
    if (!onSortChange) return;

    // Cycle through: no sort -> asc -> desc -> no sort
    if (sortBy !== column) {
      onSortChange(column, 'asc');
    } else if (sortOrder === 'asc') {
      onSortChange(column, 'desc');
    } else {
      onSortChange(null, null);
    }
  };

  const getSortIcon = (column: 'name' | 'total_count' | 'gender_balance' | 'rank') => {
    if (sortBy !== column) {
      return (
        <svg className="w-4 h-4 text-gray-300 group-hover:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
        </svg>
      );
    }
    
    if (sortOrder === 'asc') {
      return (
        <svg className="w-4 h-4 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
        </svg>
      );
    }
    
    return (
      <svg className="w-4 h-4 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
      </svg>
    );
  };

  // Loading state with shimmer skeleton
  if (isLoading) {
    return (
      <div className="bg-white rounded-2xl shadow-medium overflow-hidden border border-gray-100">
        <div className="p-8">
          <div className="flex items-center justify-center mb-6">
            <div className="relative">
              <div className="w-16 h-16 border-4 border-primary-200 border-t-primary-600 rounded-full animate-spin"></div>
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="w-8 h-8 bg-primary-100 rounded-full"></div>
              </div>
            </div>
          </div>
          <p className="text-center text-gray-600 font-medium">{t('common:labels.loading')}</p>
          
          {/* Skeleton Rows */}
          <div className="mt-8 space-y-4">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="skeleton-shimmer h-20 rounded-xl"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="bg-white rounded-2xl shadow-medium overflow-hidden border border-red-100">
        <div className="p-12 text-center">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-red-100 rounded-full mb-4">
            <svg className="w-8 h-8 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            {t('common:labels.error')}
          </h3>
          <p className="text-gray-600">{error.message}</p>
        </div>
      </div>
    );
  }

  // Empty state
  if (!names || names.length === 0) {
    return (
      <div className="bg-white rounded-2xl shadow-medium overflow-hidden border border-gray-100">
        <div className="p-12 text-center">
          <div className="inline-flex items-center justify-center w-20 h-20 bg-gray-100 rounded-full mb-4">
            <svg className="w-10 h-10 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            {t('pages:namesExplorer.emptyState.title')}
          </h3>
          <p className="text-gray-600">
            {t('pages:namesExplorer.emptyState.description')}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-2xl shadow-medium hover:shadow-strong transition-shadow duration-300 overflow-hidden border border-gray-100">
      {/* Desktop Table */}
      <div className="hidden md:block overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gradient-to-r from-gray-50 to-gray-100 border-b-2 border-gray-200">
            <tr>
              <th className="px-6 py-4 text-left text-xs font-bold text-gray-700 uppercase tracking-wider">
                <button
                  onClick={() => handleSort('rank')}
                  className={`flex items-center gap-2 group hover:text-primary-600 transition-colors ${
                    sortBy === 'rank' ? 'text-primary-600' : ''
                  } ${onSortChange ? 'cursor-pointer' : 'cursor-default'}`}
                  disabled={!onSortChange}
                >
                  <svg className="w-4 h-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14" />
                  </svg>
                  {t('pages:namesTable.rank')}
                  {onSortChange && getSortIcon('rank')}
                </button>
              </th>
              <th className="px-6 py-4 text-left text-xs font-bold text-gray-700 uppercase tracking-wider">
                <button
                  onClick={() => handleSort('name')}
                  className={`flex items-center gap-2 group hover:text-primary-600 transition-colors ${
                    sortBy === 'name' ? 'text-primary-600' : ''
                  } ${onSortChange ? 'cursor-pointer' : 'cursor-default'}`}
                  disabled={!onSortChange}
                >
                  <svg className="w-4 h-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  {t('pages:namesTable.name')}
                  {onSortChange && getSortIcon('name')}
                </button>
              </th>
              <th className="px-6 py-4 text-left text-xs font-bold text-gray-700 uppercase tracking-wider">
                <button
                  onClick={() => handleSort('total_count')}
                  className={`flex items-center gap-2 group hover:text-primary-600 transition-colors ${
                    sortBy === 'total_count' ? 'text-primary-600' : ''
                  } ${onSortChange ? 'cursor-pointer' : 'cursor-default'}`}
                  disabled={!onSortChange}
                >
                  <svg className="w-4 h-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                  {t('pages:namesTable.totalCount')}
                  {onSortChange && getSortIcon('total_count')}
                </button>
              </th>
              <th className="px-6 py-4 text-left text-xs font-bold text-gray-700 uppercase tracking-wider">
                <button
                  onClick={() => handleSort('gender_balance')}
                  className={`flex items-center gap-2 group hover:text-primary-600 transition-colors ${
                    sortBy === 'gender_balance' ? 'text-primary-600' : ''
                  } ${onSortChange ? 'cursor-pointer' : 'cursor-default'}`}
                  disabled={!onSortChange}
                >
                  <svg className="w-4 h-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 6l3 1m0 0l-3 9a5.002 5.002 0 006.001 0M6 7l3 9M6 7l6-2m6 2l3-1m-3 1l-3 9a5.002 5.002 0 006.001 0M18 7l3 9m-3-9l-6-2m0-2v2m0 16V5m0 16H9m3 0h3" />
                  </svg>
                  {t('pages:namesTable.genderBalance')}
                  {onSortChange && getSortIcon('gender_balance')}
                </button>
              </th>
              <th className="px-6 py-4 text-left text-xs font-bold text-gray-700 uppercase tracking-wider">
                <div className="flex items-center gap-2">
                  <svg className="w-4 h-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  {t('pages:namesTable.countries')}
                </div>
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-100">
            {names.map((name, index) => (
              <tr 
                key={name.name} 
                className={`
                  transition-all duration-200 hover:bg-gradient-to-r hover:from-primary-50 hover:to-transparent
                  ${index % 2 === 0 ? 'bg-white' : 'bg-gray-50/50'}
                `}
              >
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center gap-2">
                    <div className="w-8 h-8 bg-gradient-to-br from-primary-100 to-secondary-100 rounded-lg flex items-center justify-center shadow-sm">
                      <span className="text-xs font-bold text-primary-700">#{name.rank}</span>
                    </div>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <Link
                    to={getNameDetailPath(name.name)}
                    className="text-primary-600 hover:text-primary-700 font-semibold text-lg transition-colors hover:underline"
                  >
                    {name.name}
                  </Link>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center gap-2">
                    <span className="text-gray-900 font-semibold">
                      {formatNumber(name.total_count)}
                    </span>
                  </div>
                </td>
                <td className="px-6 py-4">
                  <div className="w-56">
                    <GenderBalanceBar
                      genderBalance={name.gender_balance}
                      femaleCount={name.female_count}
                      maleCount={name.male_count}
                      hasUnknownData={name.has_unknown_data}
                    />
                  </div>
                </td>
                <td className="px-6 py-4">
                  <div className="flex flex-wrap gap-1.5">
                    {name.countries.map((country) => (
                      <span
                        key={country}
                        className="inline-flex items-center px-2.5 py-1 rounded-lg text-xs font-semibold bg-gradient-to-r from-secondary-100 to-secondary-200 text-secondary-800 border border-secondary-300 shadow-sm hover:shadow-md transition-shadow"
                      >
                        {country}
                      </span>
                    ))}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Mobile Cards */}
      <div className="md:hidden divide-y divide-gray-200">
        {names.map((name) => (
          <div key={name.name} className="p-6 hover:bg-gradient-to-r hover:from-primary-50 hover:to-transparent transition-colors">
            <div className="flex items-start justify-between mb-4">
              <div className="flex-1">
                <Link
                  to={getNameDetailPath(name.name)}
                  className="text-primary-600 hover:text-primary-700 font-semibold text-xl transition-colors"
                >
                  {name.name}
                </Link>
                <div className="flex items-center gap-3 mt-2 text-sm text-gray-600">
                  <div className="flex items-center gap-1">
                    <div className="w-6 h-6 bg-gradient-to-br from-primary-100 to-secondary-100 rounded-md flex items-center justify-center">
                      <span className="text-xs font-bold text-primary-700">#{name.rank}</span>
                    </div>
                    <span>{t('pages:namesTable.mobile.rankLabel')}:</span>
                  </div>
                  <span className="font-medium">{formatNumber(name.total_count)} {t('pages:namesTable.mobile.totalLabel')}</span>
                </div>
              </div>
            </div>

            <div className="space-y-4">
              <div>
                <p className="text-xs font-semibold text-gray-500 mb-2 uppercase tracking-wide">{t('pages:namesTable.genderBalance')}</p>
                <GenderBalanceBar
                  genderBalance={name.gender_balance}
                  femaleCount={name.female_count}
                  maleCount={name.male_count}
                  hasUnknownData={name.has_unknown_data}
                />
              </div>

              <div>
                <p className="text-xs font-semibold text-gray-500 mb-1 uppercase tracking-wide">{t('pages:namesTable.countries')}</p>
                <div className="flex flex-wrap gap-1">
                  {name.countries.map((country) => (
                    <span
                      key={country}
                      className="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-semibold bg-gradient-to-r from-secondary-100 to-secondary-200 text-secondary-800"
                    >
                      {country}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}