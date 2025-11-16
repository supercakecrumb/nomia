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
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      <div className="container mx-auto px-4 py-8 md:py-12">
        {/* Page Header */}
        <div className="mb-8 md:mb-12 animate-fade-in">
          <div className="flex items-center gap-4 mb-4">
            <div className="w-12 h-12 md:w-14 md:h-14 bg-gradient-to-br from-primary-600 to-secondary-600 rounded-2xl flex items-center justify-center shadow-lg">
              <svg className="w-6 h-6 md:w-7 md:h-7 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <div>
              <h1 className="text-3xl md:text-4xl font-bold text-gray-900 tracking-tight">
                {t('namesExplorer.title')}
              </h1>
              <p className="text-base md:text-lg text-gray-600 mt-1">
                {t('namesExplorer.subtitle')}
              </p>
            </div>
          </div>
        </div>

        {/* Connection Status Banner */}
        {data && (
          <div className="mb-6 animate-slide-down">
            <div className="bg-gradient-to-r from-green-50 to-emerald-50 border-2 border-green-200 rounded-xl p-4 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="flex-shrink-0">
                  <div className="w-10 h-10 bg-green-500 rounded-full flex items-center justify-center shadow-md">
                    <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  </div>
                </div>
                <div className="flex-1">
                  <p className="text-sm font-semibold text-green-900">
                    ✓ API Connected Successfully
                  </p>
                  <p className="text-xs text-green-700 mt-0.5">
                    {data.names.length} names loaded • Page {data.meta.page} of {data.meta.total_pages} • {data.meta.total_count.toLocaleString()} total results
                  </p>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Filter Bar */}
        <div className="mb-8 animate-slide-up">
          <FilterBar />
        </div>

        {/* Names Table */}
        <div className="mb-6 animate-slide-up animation-delay-100">
          <NamesTable
            names={data?.names || []}
            isLoading={isLoading}
            error={error}
          />
        </div>

        {/* Pagination */}
        {data && data.names.length > 0 && (
          <div className="animate-slide-up animation-delay-200">
            <Pagination
              currentPage={data.meta.page}
              totalPages={data.meta.total_pages}
              totalCount={data.meta.total_count}
              pageSize={data.meta.page_size}
              onPageChange={setPage}
            />
          </div>
        )}

        {/* Empty State Encouragement */}
        {!isLoading && !error && data && data.names.length === 0 && (
          <div className="mt-8 text-center">
            <div className="inline-flex items-center gap-2 px-6 py-3 bg-primary-50 border-2 border-primary-200 rounded-xl text-primary-700 font-medium">
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Try adjusting your filters to see more results
            </div>
          </div>
        )}
      </div>
    </div>
  );
}