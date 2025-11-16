/**
 * Pagination Component
 * 
 * Displays pagination controls for navigating through pages of data.
 */

import { useTranslation } from 'react-i18next';
import { formatNumber } from '../../utils/formatters';

interface PaginationProps {
  /** Current page number (1-based) */
  currentPage: number;
  /** Total number of pages */
  totalPages: number;
  /** Total count of items */
  totalCount: number;
  /** Number of items per page */
  pageSize: number;
  /** Callback when page changes */
  onPageChange: (page: number) => void;
  /** Callback when page size changes */
  onPageSizeChange?: (pageSize: number) => void;
}

export default function Pagination({
  currentPage,
  totalPages,
  totalCount,
  pageSize,
  onPageChange,
  onPageSizeChange,
}: PaginationProps) {
  const { t } = useTranslation('common');

  const startItem = (currentPage - 1) * pageSize + 1;
  const endItem = Math.min(currentPage * pageSize, totalCount);

  const canGoPrevious = currentPage > 1;
  const canGoNext = currentPage < totalPages;

  const pageSizeOptions = [20, 50, 100, 200];

  return (
    <div className="bg-white rounded-2xl shadow-medium border border-gray-100 overflow-hidden">
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 px-6 py-4">
        {/* Results info with page size selector */}
        <div className="flex items-center gap-3">
          <div className="text-sm text-gray-700 flex items-center gap-2">
            <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span>
              {t('pagination.showing')} <span className="font-semibold text-gray-900">{formatNumber(startItem)}</span>{' '}
              {t('pagination.to')} <span className="font-semibold text-gray-900">{formatNumber(endItem)}</span>{' '}
              {t('pagination.of')} <span className="font-semibold text-gray-900">{formatNumber(totalCount)}</span>{' '}
              {t('pagination.results')}
            </span>
          </div>
          
          {/* Page Size Selector */}
          {onPageSizeChange && (
            <div className="flex items-center gap-2">
              <span className="text-xs text-gray-500">|</span>
              <select
                value={pageSize}
                onChange={(e) => onPageSizeChange(parseInt(e.target.value))}
                className="px-3 py-1.5 text-sm font-medium bg-white border-2 border-gray-300 rounded-lg hover:border-primary-400 focus:border-primary-500 focus:ring-2 focus:ring-primary-200 transition-all duration-200 cursor-pointer"
              >
                {pageSizeOptions.map((size) => (
                  <option key={size} value={size}>
                    {size} per page
                  </option>
                ))}
              </select>
            </div>
          )}
        </div>

        {/* Navigation buttons */}
        <div className="flex items-center gap-2">
          {/* Previous Button */}
          <button
            onClick={() => onPageChange(currentPage - 1)}
            disabled={!canGoPrevious}
            className={`
              flex items-center gap-2 px-4 py-2.5 text-sm font-medium rounded-xl transition-all duration-200
              ${
                canGoPrevious
                  ? 'bg-white text-gray-700 border-2 border-gray-300 hover:border-primary-400 hover:bg-primary-50 hover:text-primary-700 shadow-sm hover:shadow-md active:scale-95'
                  : 'bg-gray-100 text-gray-400 border-2 border-gray-200 cursor-not-allowed'
              }
            `}
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
            {t('pagination.previous')}
          </button>

          {/* Page Info */}
          <div className="hidden sm:flex items-center gap-2 px-4 py-2.5 bg-gradient-to-r from-primary-50 to-secondary-50 rounded-xl border-2 border-primary-200">
            <svg className="w-4 h-4 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14" />
            </svg>
            <span className="text-sm font-semibold text-gray-700">
              {t('pagination.page')} <span className="text-primary-700">{currentPage}</span>{' '}
              {t('pagination.of')} <span className="text-primary-700">{totalPages}</span>
            </span>
          </div>

          {/* Mobile Page Info */}
          <div className="sm:hidden flex items-center gap-2 px-3 py-2 bg-gradient-to-r from-primary-50 to-secondary-50 rounded-xl border-2 border-primary-200">
            <span className="text-sm font-semibold text-primary-700">
              {currentPage}/{totalPages}
            </span>
          </div>

          {/* Next Button */}
          <button
            onClick={() => onPageChange(currentPage + 1)}
            disabled={!canGoNext}
            className={`
              flex items-center gap-2 px-4 py-2.5 text-sm font-medium rounded-xl transition-all duration-200
              ${
                canGoNext
                  ? 'bg-gradient-to-r from-primary-600 to-primary-700 text-white hover:from-primary-700 hover:to-primary-800 shadow-md hover:shadow-lg active:scale-95'
                  : 'bg-gray-100 text-gray-400 border-2 border-gray-200 cursor-not-allowed'
              }
            `}
          >
            {t('pagination.next')}
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}