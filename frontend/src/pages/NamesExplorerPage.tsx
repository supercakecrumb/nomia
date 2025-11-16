/**
 * Names Explorer Page
 * 
 * Main page for browsing and filtering gender-neutral names.
 * Includes comprehensive filters with URL synchronization.
 */

import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNames } from '../hooks/useNames';
import { useFilters } from '../hooks/useFilters';
import { useMetaYears } from '../hooks/useMetaYears';
import FilterBar from '../components/filters/FilterBar';
import NamesTable from '../components/table/NamesTable';
import Pagination from '../components/table/Pagination';

export default function NamesExplorerPage() {
  const { t } = useTranslation('pages');
  const { data: metaYears } = useMetaYears();
  
  // Initialize filters with meta years
  const {
    filters,
    setYearRange,
    setCountries,
    setGenderBalance,
    setMinCount,
    setTopN,
    setCoveragePercent,
    setSearch,
    setSort,
    setPage,
    setPageSize,
    resetFilters,
    updateDerivedValues,
    getApiParams,
  } = useFilters(metaYears?.min_year || 1880, metaYears?.max_year || 2025);
  
  // Fetch names with current filters
  const { data, isLoading, error, isFetching } = useNames(getApiParams());

  // Debug logging for API response
  useEffect(() => {
    if (data) {
      console.log('[NamesExplorerPage] API Response:', {
        hasData: !!data,
        hasNames: !!data.names,
        namesIsArray: Array.isArray(data.names),
        namesLength: data?.names?.length ?? 'null/undefined',
        hasMeta: !!data.meta,
        meta: data.meta,
        fullResponse: data,
      });
    }
  }, [data]);

  // Update derived popularity values when API response changes
  useEffect(() => {
    if (data?.meta?.popularity_summary) {
      const { derived_min_count, derived_top_n, derived_coverage_percent } = data.meta.popularity_summary;
      updateDerivedValues(derived_min_count, derived_top_n, derived_coverage_percent);
    }
  }, [data, updateDerivedValues]);

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      {/* Loading Banner - Shows when fetching new data */}
      {isFetching && data && (
        <div className="fixed top-0 left-0 right-0 bg-primary-600 text-white py-3 px-4 flex items-center justify-center gap-3 z-50 shadow-lg">
          <div className="inline-block animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
          <span className="text-sm font-medium">Loading names... This may take a few seconds</span>
        </div>
      )}
      
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

        {/* Filter Bar */}
        <div className="mb-8 animate-slide-up">
          <FilterBar
            yearMin={filters.yearMin}
            yearMax={filters.yearMax}
            onYearRangeChange={setYearRange}
            selectedCountries={filters.countries}
            onCountriesChange={setCountries}
            genderBalanceMin={filters.genderBalanceMin}
            genderBalanceMax={filters.genderBalanceMax}
            onGenderBalanceChange={setGenderBalance}
            minCount={filters.minCount}
            topN={filters.topN}
            coveragePercent={filters.coveragePercent}
            popularityDriver={filters.popularityDriver}
            onMinCountChange={setMinCount}
            onTopNChange={setTopN}
            onCoveragePercentChange={setCoveragePercent}
            nameSearch={filters.search}
            onNameSearchChange={setSearch}
            onReset={resetFilters}
          />
        </div>

        {/* Names Table */}
        <div className="mb-6 animate-slide-up animation-delay-100">
          <NamesTable
            names={data?.names || []}
            isLoading={isLoading}
            error={error}
            sortBy={filters.sortBy}
            sortOrder={filters.sortOrder}
            onSortChange={setSort}
          />
        </div>

        {/* Pagination */}
        {data?.names && data?.meta && data.names.length > 0 && (
          <div className="animate-slide-up animation-delay-200">
            <Pagination
              currentPage={data.meta.page}
              totalPages={data.meta.total_pages}
              totalCount={data.meta.total_count}
              pageSize={data.meta.page_size}
              onPageChange={setPage}
              onPageSizeChange={setPageSize}
            />
          </div>
        )}

        {/* Empty State Encouragement */}
        {!isLoading && !error && data?.names && data.names.length === 0 && (
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