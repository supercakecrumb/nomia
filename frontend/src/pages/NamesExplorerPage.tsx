/**
 * Names Explorer Page
 * 
 * Main page for browsing and filtering gender-neutral names.
 * Phase 1: Display fixture data with basic UI components.
 */

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNames } from '../hooks/useNames';
import FilterBar from '../components/filters/FilterBar';
import NamesTable from '../components/table/NamesTable';
import Pagination from '../components/table/Pagination';

export default function NamesExplorerPage() {
  const { t } = useTranslation('pages');
  const [page, setPage] = useState(1);
  
  // Fetch names with pagination (no filters for Phase 1)
  const { data, isLoading, error } = useNames({ page, page_size: 20 });

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">
            {t('namesExplorer.title')}
          </h1>
          <p className="text-lg text-gray-600">
            {t('namesExplorer.subtitle')}
          </p>
        </div>

        {/* Connection Status */}
        {data && (
          <div className="mb-6 p-3 bg-green-50 border border-green-200 rounded-lg">
            <p className="text-sm text-green-800">
              ✓ API Connected - {data.names.length} names loaded (Page {data.meta.page} of {data.meta.total_pages})
            </p>
          </div>
        )}

        {/* Filter Bar */}
        <div className="mb-6">
          <FilterBar />
        </div>

        {/* Names Table */}
        <div className="mb-6">
          <NamesTable
            names={data?.names || []}
            isLoading={isLoading}
            error={error}
          />
        </div>

        {/* Pagination */}
        {data && data.names.length > 0 && (
          <Pagination
            currentPage={data.meta.page}
            totalPages={data.meta.total_pages}
            totalCount={data.meta.total_count}
            pageSize={data.meta.page_size}
            onPageChange={setPage}
          />
        )}
      </div>
    </div>
  );
}