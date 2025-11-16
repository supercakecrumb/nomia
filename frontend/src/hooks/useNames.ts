/**
 * useNames Hook
 * 
 * React Query hook for fetching a paginated list of gender-neutral names.
 * Supports filtering by countries, years, gender balance, sorting, and search.
 * Uses keepPreviousData to prevent UI flickering during pagination.
 */

import { useQuery } from '@tanstack/react-query';
import { fetchNames } from '../api/client';
import type { NamesListResponse, NamesFilterParams } from '../types/api';

/**
 * Hook to fetch a paginated list of gender-neutral names with optional filters
 * 
 * @param filters - Optional filter parameters
 * @returns React Query result with names list, pagination, and metadata
 * 
 * @example
 * ```typescript
 * function NamesList() {
 *   const [page, setPage] = useState(1);
 *   const { data, isLoading, error } = useNames({
 *     page,
 *     page_size: 20,
 *     countries: ['US', 'UK'],
 *     year_min: 2000,
 *     sort_by: 'total_count',
 *     sort_order: 'desc'
 *   });
 *   
 *   if (isLoading) return <div>Loading...</div>;
 *   if (error) return <div>Error: {error.message}</div>;
 *   
 *   return (
 *     <div>
 *       <ul>
 *         {data.names.map(name => (
 *           <li key={name.name}>{name.name} - {name.total_count}</li>
 *         ))}
 *       </ul>
 *       <button onClick={() => setPage(p => p + 1)}>Next Page</button>
 *     </div>
 *   );
 * }
 * ```
 */
export const useNames = (filters: NamesFilterParams = {}) => {
  return useQuery<NamesListResponse, Error>({
    queryKey: ['names', filters],
    queryFn: () => fetchNames(filters),
    placeholderData: (previousData) => previousData, // Keep previous data while fetching new page
  });
};