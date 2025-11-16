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
}

export default function Pagination({
  currentPage,
  totalPages,
  totalCount,
  pageSize,
  onPageChange,
}: PaginationProps) {
  const { t } = useTranslation('common');

  const startItem = (currentPage - 1) * pageSize + 1;
  const endItem = Math.min(currentPage * pageSize, totalCount);

  const canGoPrevious = currentPage > 1;
  const canGoNext = currentPage < totalPages;

  return (
    <div className="flex flex-col sm:flex-row items-center justify-between gap-4 px-4 py-3 bg-white border-t border-gray-200">
      {/* Results info */}
      <div className="text-sm text-gray-700">
        {t('pagination.showing')} <span className="font-medium">{formatNumber(startItem)}</span>{' '}
        {t('pagination.to')} <span className="font-medium">{formatNumber(endItem)}</span>{' '}
        {t('pagination.of')} <span className="font-medium">{formatNumber(totalCount)}</span>{' '}
        {t('pagination.results')}
      </div>

      {/* Navigation buttons */}
      <div className="flex items-center gap-2">
        <button
          onClick={() => onPageChange(currentPage - 1)}
          disabled={!canGoPrevious}
          className={`
            px-3 py-2 text-sm font-medium rounded-md
            ${
              canGoPrevious
                ? 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50'
                : 'bg-gray-100 text-gray-400 border border-gray-200 cursor-not-allowed'
            }
          `}
        >
          {t('pagination.previous')}
        </button>

        <div className="text-sm text-gray-700">
          {t('pagination.page')} <span className="font-medium">{currentPage}</span>{' '}
          {t('pagination.of')} <span className="font-medium">{totalPages}</span>
        </div>

        <button
          onClick={() => onPageChange(currentPage + 1)}
          disabled={!canGoNext}
          className={`
            px-3 py-2 text-sm font-medium rounded-md
            ${
              canGoNext
                ? 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50'
                : 'bg-gray-100 text-gray-400 border border-gray-200 cursor-not-allowed'
            }
          `}
        >
          {t('pagination.next')}
        </button>
      </div>
    </div>
  );
}