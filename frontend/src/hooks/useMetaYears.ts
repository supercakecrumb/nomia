/**
 * useMetaYears Hook
 * 
 * React Query hook for fetching the available year range from the API.
 * Data is cached indefinitely (staleTime: Infinity) as it rarely changes.
 */

import { useQuery } from '@tanstack/react-query';
import { fetchMetaYears } from '../api/client';
import type { MetaYearsResponse } from '../types/api';

/**
 * Hook to fetch the available year range in the database
 * 
 * @returns React Query result with min_year and max_year data
 * 
 * @example
 * ```typescript
 * function YearRangeDisplay() {
 *   const { data, isLoading, error } = useMetaYears();
 *   
 *   if (isLoading) return <div>Loading...</div>;
 *   if (error) return <div>Error: {error.message}</div>;
 *   
 *   return <div>Data from {data.min_year} to {data.max_year}</div>;
 * }
 * ```
 */
export const useMetaYears = () => {
  return useQuery<MetaYearsResponse, Error>({
    queryKey: ['meta', 'years'],
    queryFn: fetchMetaYears,
    staleTime: Infinity, // Data rarely changes, keep it fresh indefinitely
  });
};