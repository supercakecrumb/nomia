/**
 * useNameTrend Hook
 * 
 * React Query hook for fetching detailed trend data for a specific name.
 * Includes time series data and country-specific statistics.
 */

import { useQuery } from '@tanstack/react-query';
import { fetchNameTrend } from '../api/client';
import type { NameTrendResponse, NameTrendParams } from '../types/api';

/**
 * Hook to fetch detailed trend data for a specific name
 * 
 * @param params - Parameters including name and optional filters
 * @param enabled - Whether the query should run (defaults to true if name is provided)
 * @returns React Query result with name trend data, time series, and country breakdown
 * 
 * @example
 * ```typescript
 * function NameDetail({ name }: { name: string }) {
 *   const { data, isLoading, error } = useNameTrend({ name });
 *   
 *   if (isLoading) return <div>Loading...</div>;
 *   if (error) return <div>Error: {error.message}</div>;
 *   
 *   return (
 *     <div>
 *       <h1>{data.name}</h1>
 *       <p>Total occurrences: {data.summary.total_count}</p>
 *       <p>Gender balance: {data.summary.gender_balance}% male</p>
 *       <h2>By Country</h2>
 *       <ul>
 *         {data.by_country.map(c => (
 *           <li key={c.country_code}>
 *             {c.country_name}: {c.total_count}
 *           </li>
 *         ))}
 *       </ul>
 *     </div>
 *   );
 * }
 * ```
 */
export const useNameTrend = (
  params: NameTrendParams,
  options?: { enabled?: boolean }
) => {
  const enabled = options?.enabled ?? Boolean(params.name);
  
  return useQuery<NameTrendResponse, Error>({
    queryKey: ['names', 'trend', params],
    queryFn: () => fetchNameTrend(params),
    enabled, // Only fetch if name is provided and enabled is true
  });
};