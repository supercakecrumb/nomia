/**
 * NamesTable Component
 * 
 * Displays a table of gender-neutral names with various statistics.
 * Supports loading, error, and empty states.
 */

import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import type { NameEntry } from '../../types/api';
import { formatNumber, formatYearRange } from '../../utils/formatters';
import GenderBalanceBar from './GenderBalanceBar';
import { getNameDetailPath } from '../../utils/navigation';

interface NamesTableProps {
  /** Array of name entries to display */
  names: NameEntry[];
  /** Loading state */
  isLoading?: boolean;
  /** Error state */
  error?: Error | null;
}

export default function NamesTable({ names, isLoading, error }: NamesTableProps) {
  const { t } = useTranslation(['common', 'pages']);

  // Loading state
  if (isLoading) {
    return (
      <div className="bg-white rounded-lg shadow-md">
        <div className="p-8 text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-4 border-blue-600 border-t-transparent"></div>
          <p className="mt-4 text-gray-600">{t('common:labels.loading')}</p>
        </div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="bg-white rounded-lg shadow-md">
        <div className="p-8 text-center">
          <div className="text-red-500 text-5xl mb-4">⚠️</div>
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
      <div className="bg-white rounded-lg shadow-md">
        <div className="p-8 text-center">
          <div className="text-gray-400 text-5xl mb-4">🔍</div>
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
    <div className="bg-white rounded-lg shadow-md overflow-hidden">
      {/* Desktop Table */}
      <div className="hidden md:block overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Rank
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Name
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Total Count
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Gender Balance
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Period
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Countries
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {names.map((name) => (
              <tr key={name.name} className="hover:bg-gray-50 transition-colors">
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  #{name.rank}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <Link
                    to={getNameDetailPath(name.name)}
                    className="text-blue-600 hover:text-blue-800 font-medium text-lg"
                  >
                    {name.name}
                  </Link>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 font-medium">
                  {formatNumber(name.total_count)}
                </td>
                <td className="px-6 py-4">
                  <div className="w-48">
                    <GenderBalanceBar
                      genderBalance={name.gender_balance}
                      femaleCount={name.female_count}
                      maleCount={name.male_count}
                      hasUnknownData={name.has_unknown_data}
                    />
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {formatYearRange(name.name_start, name.name_end)}
                </td>
                <td className="px-6 py-4">
                  <div className="flex flex-wrap gap-1">
                    {name.countries.map((country) => (
                      <span
                        key={country}
                        className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800"
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
          <div key={name.name} className="p-4 hover:bg-gray-50 transition-colors">
            <div className="flex items-start justify-between mb-3">
              <div>
                <Link
                  to={getNameDetailPath(name.name)}
                  className="text-blue-600 hover:text-blue-800 font-medium text-lg"
                >
                  {name.name}
                </Link>
                <p className="text-sm text-gray-500 mt-1">
                  Rank #{name.rank} • {formatNumber(name.total_count)} total
                </p>
              </div>
            </div>

            <div className="space-y-3">
              <div>
                <p className="text-xs text-gray-500 mb-1">Gender Balance</p>
                <GenderBalanceBar
                  genderBalance={name.gender_balance}
                  femaleCount={name.female_count}
                  maleCount={name.male_count}
                  hasUnknownData={name.has_unknown_data}
                />
              </div>

              <div className="flex justify-between text-sm">
                <div>
                  <p className="text-xs text-gray-500">Period</p>
                  <p className="text-gray-900 font-medium">
                    {formatYearRange(name.name_start, name.name_end)}
                  </p>
                </div>
                <div>
                  <p className="text-xs text-gray-500">Countries</p>
                  <div className="flex gap-1 mt-1">
                    {name.countries.map((country) => (
                      <span
                        key={country}
                        className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800"
                      >
                        {country}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}