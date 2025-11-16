/**
 * useFilters Hook
 * 
 * Manages all filter state with URL synchronization and debouncing.
 * Handles year range, countries, gender balance, popularity trio, search, and sort.
 */

import { useState, useEffect, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useDebounce } from './useDebounce';
import type { NamesFilterParams } from '../types/api';

type PopularityDriver = 'min_count' | 'top_n' | 'coverage_percent' | null;

interface FilterState {
  // Year range
  yearMin: number;
  yearMax: number;
  
  // Countries
  countries: string[];
  
  // Gender balance
  genderBalanceMin: number;
  genderBalanceMax: number;
  
  // Popularity trio
  minCount: number | null;
  topN: number | null;
  coveragePercent: number | null;
  popularityDriver: PopularityDriver;
  
  // Name search
  search: string;
  
  // Sort
  sortBy: 'name' | 'total_count' | 'gender_balance' | 'rank' | null;
  sortOrder: 'asc' | 'desc' | null;
  
  // Pagination
  page: number;
  pageSize: number;
}

export const useFilters = (defaultYearMin: number = 1880, defaultYearMax: number = 2023) => {
  const [searchParams, setSearchParams] = useSearchParams();
  
  // Initialize state from URL params
  const [filters, setFilters] = useState<FilterState>(() => {
    const yearMin = parseInt(searchParams.get('year_min') || String(defaultYearMin));
    const yearMax = parseInt(searchParams.get('year_max') || String(defaultYearMax));
    const countries = searchParams.get('countries')?.split(',').filter(Boolean) || [];
    const genderBalanceMin = parseInt(searchParams.get('gender_balance_min') || '0');
    const genderBalanceMax = parseInt(searchParams.get('gender_balance_max') || '100');
    const minCount = searchParams.get('min_count') ? parseInt(searchParams.get('min_count')!) : null;
    const topN = searchParams.get('top_n') ? parseInt(searchParams.get('top_n')!) : null;
    const coveragePercent = searchParams.get('coverage_percent') ? parseFloat(searchParams.get('coverage_percent')!) : null;
    const search = searchParams.get('search') || '';
    const sortBy = (searchParams.get('sort_by') as any) || null;
    const sortOrder = (searchParams.get('sort_order') as 'asc' | 'desc') || null;
    const page = parseInt(searchParams.get('page') || '1');
    const pageSize = parseInt(searchParams.get('page_size') || '100');
    
    // Determine active driver
    let popularityDriver: PopularityDriver = null;
    if (minCount !== null) popularityDriver = 'min_count';
    else if (topN !== null) popularityDriver = 'top_n';
    else if (coveragePercent !== null) popularityDriver = 'coverage_percent';
    
    return {
      yearMin,
      yearMax,
      countries,
      genderBalanceMin,
      genderBalanceMax,
      minCount,
      topN,
      coveragePercent,
      popularityDriver,
      search,
      sortBy,
      sortOrder,
      page,
      pageSize,
    };
  });
  
  // Debounce slider values (300ms)
  const debouncedYearMin = useDebounce(filters.yearMin, 300);
  const debouncedYearMax = useDebounce(filters.yearMax, 300);
  const debouncedGenderBalanceMin = useDebounce(filters.genderBalanceMin, 300);
  const debouncedGenderBalanceMax = useDebounce(filters.genderBalanceMax, 300);
  
  // Update URL params when filters change
  useEffect(() => {
    const params = new URLSearchParams();
    
    // Only add non-default values
    if (debouncedYearMin !== defaultYearMin) params.set('year_min', String(debouncedYearMin));
    if (debouncedYearMax !== defaultYearMax) params.set('year_max', String(debouncedYearMax));
    if (filters.countries.length > 0) params.set('countries', filters.countries.join(','));
    if (debouncedGenderBalanceMin !== 0) params.set('gender_balance_min', String(debouncedGenderBalanceMin));
    if (debouncedGenderBalanceMax !== 100) params.set('gender_balance_max', String(debouncedGenderBalanceMax));
    if (filters.minCount !== null) params.set('min_count', String(filters.minCount));
    if (filters.topN !== null) params.set('top_n', String(filters.topN));
    if (filters.coveragePercent !== null) params.set('coverage_percent', String(filters.coveragePercent));
    if (filters.search) params.set('search', filters.search);
    if (filters.sortBy) params.set('sort_by', filters.sortBy);
    if (filters.sortOrder) params.set('sort_order', filters.sortOrder);
    if (filters.page > 1) params.set('page', String(filters.page));
    if (filters.pageSize !== 100) params.set('page_size', String(filters.pageSize));
    
    setSearchParams(params, { replace: true });
  }, [
    debouncedYearMin,
    debouncedYearMax,
    filters.countries,
    debouncedGenderBalanceMin,
    debouncedGenderBalanceMax,
    filters.minCount,
    filters.topN,
    filters.coveragePercent,
    filters.search,
    filters.sortBy,
    filters.sortOrder,
    filters.page,
    filters.pageSize,
    defaultYearMin,
    defaultYearMax,
    setSearchParams,
  ]);
  
  // Convert to API params
  const getApiParams = useCallback((): NamesFilterParams => {
    const params: NamesFilterParams = {
      page: filters.page,
      page_size: filters.pageSize,
    };
    
    if (debouncedYearMin !== defaultYearMin) params.year_min = debouncedYearMin;
    if (debouncedYearMax !== defaultYearMax) params.year_max = debouncedYearMax;
    if (filters.countries.length > 0) params.countries = filters.countries;
    if (debouncedGenderBalanceMin !== 0) params.gender_balance_min = debouncedGenderBalanceMin;
    if (debouncedGenderBalanceMax !== 100) params.gender_balance_max = debouncedGenderBalanceMax;
    if (filters.minCount !== null) params.min_count = filters.minCount;
    if (filters.topN !== null) params.top_n = filters.topN;
    if (filters.coveragePercent !== null) params.coverage_percent = filters.coveragePercent;
    if (filters.search) params.search = filters.search;
    if (filters.sortBy) params.sort_by = filters.sortBy;
    if (filters.sortOrder) params.sort_order = filters.sortOrder;
    
    return params;
  }, [
    filters,
    debouncedYearMin,
    debouncedYearMax,
    debouncedGenderBalanceMin,
    debouncedGenderBalanceMax,
    defaultYearMin,
    defaultYearMax,
  ]);
  
  // Update methods
  const setYearRange = useCallback((min: number, max: number) => {
    setFilters(prev => ({ ...prev, yearMin: min, yearMax: max, page: 1 }));
  }, []);
  
  const setCountries = useCallback((countries: string[]) => {
    setFilters(prev => ({ ...prev, countries, page: 1 }));
  }, []);
  
  const setGenderBalance = useCallback((min: number, max: number) => {
    setFilters(prev => ({ ...prev, genderBalanceMin: min, genderBalanceMax: max, page: 1 }));
  }, []);
  
  const setMinCount = useCallback((value: number | null) => {
    setFilters(prev => ({
      ...prev,
      minCount: value,
      topN: null,
      coveragePercent: null,
      popularityDriver: value !== null ? 'min_count' : null,
      page: 1,
    }));
  }, []);
  
  const setTopN = useCallback((value: number | null) => {
    setFilters(prev => ({
      ...prev,
      minCount: null,
      topN: value,
      coveragePercent: null,
      popularityDriver: value !== null ? 'top_n' : null,
      page: 1,
    }));
  }, []);
  
  const setCoveragePercent = useCallback((value: number | null) => {
    setFilters(prev => ({
      ...prev,
      minCount: null,
      topN: null,
      coveragePercent: value,
      popularityDriver: value !== null ? 'coverage_percent' : null,
      page: 1,
    }));
  }, []);
  
  const setSearch = useCallback((value: string) => {
    setFilters(prev => ({ ...prev, search: value, page: 1 }));
  }, []);
  
  const setSort = useCallback((
    sortBy: 'name' | 'total_count' | 'gender_balance' | 'rank' | null,
    sortOrder: 'asc' | 'desc' | null
  ) => {
    setFilters(prev => ({ ...prev, sortBy, sortOrder, page: 1 }));
  }, []);
  
  const setPage = useCallback((page: number) => {
    setFilters(prev => ({ ...prev, page }));
  }, []);
  
  const setPageSize = useCallback((pageSize: number) => {
    setFilters(prev => ({ ...prev, pageSize, page: 1 }));
  }, []);
  
  const resetFilters = useCallback(() => {
    setFilters({
      yearMin: defaultYearMin,
      yearMax: defaultYearMax,
      countries: [],
      genderBalanceMin: 0,
      genderBalanceMax: 100,
      minCount: null,
      topN: null,
      coveragePercent: null,
      popularityDriver: null,
      search: '',
      sortBy: null,
      sortOrder: null,
      page: 1,
      pageSize: 100,
    });
  }, [defaultYearMin, defaultYearMax]);
  
  // Update derived values from API response
  const updateDerivedValues = useCallback((
    derivedMinCount: number,
    derivedTopN: number,
    derivedCoveragePercent: number
  ) => {
    setFilters(prev => {
      const updates: Partial<FilterState> = {};
      
      if (prev.popularityDriver === 'min_count') {
        updates.topN = derivedTopN;
        updates.coveragePercent = derivedCoveragePercent;
      } else if (prev.popularityDriver === 'top_n') {
        updates.minCount = derivedMinCount;
        updates.coveragePercent = derivedCoveragePercent;
      } else if (prev.popularityDriver === 'coverage_percent') {
        updates.minCount = derivedMinCount;
        updates.topN = derivedTopN;
      }
      
      return { ...prev, ...updates };
    });
  }, []);
  
  return {
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
  };
};