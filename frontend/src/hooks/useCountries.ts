/**
 * useCountries Hook
 * 
 * React Query hook for fetching the list of available countries.
 * Data is cached indefinitely (staleTime: Infinity) as it rarely changes.
 */

import { useQuery } from '@tanstack/react-query';
import { fetchCountries } from '../api/client';
import type { CountriesResponse } from '../types/api';

/**
 * Hook to fetch the list of available countries with their data sources
 * 
 * @returns React Query result with countries array
 * 
 * @example
 * ```typescript
 * function CountrySelector() {
 *   const { data, isLoading, error } = useCountries();
 *   
 *   if (isLoading) return <div>Loading countries...</div>;
 *   if (error) return <div>Error: {error.message}</div>;
 *   
 *   return (
 *     <select>
 *       {data.countries.map(c => (
 *         <option key={c.code} value={c.code}>{c.name}</option>
 *       ))}
 *     </select>
 *   );
 * }
 * ```
 */
export const useCountries = () => {
  return useQuery<CountriesResponse, Error>({
    queryKey: ['meta', 'countries'],
    queryFn: fetchCountries,
    staleTime: Infinity, // Country list rarely changes
  });
};