/**
 * API Client
 * 
 * Provides functions for all API endpoints with support for both mock and real modes.
 * Mode is controlled via VITE_API_MODE environment variable.
 */

import type {
  MetaYearsResponse,
  CountriesResponse,
  NamesListResponse,
  NameTrendResponse,
  NamesFilterParams,
  NameTrendParams,
  ApiError,
} from '../types/api';

import {
  metaYearsFixture,
  countriesFixture,
  namesListFixture,
  nameDetailFixture,
  namesListEmptyFixture,
} from './fixtures';

// API configuration from environment variables
const API_MODE = import.meta.env.VITE_API_MODE || 'mock';
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000';
const MOCK_DELAY_MS = 200;

/**
 * Simulates network delay for mock responses
 */
const delay = (ms: number): Promise<void> => {
  return new Promise(resolve => setTimeout(resolve, ms));
};

/**
 * Builds query string from filter parameters
 */
const buildQueryString = (params: Record<string, any>): string => {
  const searchParams = new URLSearchParams();
  
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      if (Array.isArray(value)) {
        value.forEach(v => searchParams.append(key, String(v)));
      } else {
        searchParams.append(key, String(value));
      }
    }
  });
  
  const queryString = searchParams.toString();
  return queryString ? `?${queryString}` : '';
};

/**
 * Makes an HTTP request to the API with cancellation support
 */
const makeRequest = async <T>(endpoint: string, signal?: AbortSignal): Promise<T> => {
  const url = `${API_BASE_URL}${endpoint}`;
  
  console.log(`[API] Fetching: ${url}`);
  
  try {
    const response = await fetch(url, { signal });
    
    console.log(`[API] Response status: ${response.status} ${response.statusText}`);
    
    if (!response.ok) {
      const error: ApiError = await response.json().catch(() => ({
        code: 'UNKNOWN_ERROR',
        message: `HTTP ${response.status}: ${response.statusText}`,
      }));
      console.error(`[API] Error response:`, error);
      throw new Error(error.message);
    }
    
    const data = await response.json();
    console.log(`[API] Success response:`, data);
    return data;
  } catch (error) {
    console.error(`[API] Request failed:`, error);
    if (error instanceof Error) {
      throw error;
    }
    throw new Error('An unknown error occurred');
  }
};

/**
 * Fetches the available year range from the database
 * 
 * @returns Promise resolving to min and max years
 * 
 * @example
 * ```typescript
 * const { min_year, max_year } = await fetchMetaYears();
 * console.log(`Data available from ${min_year} to ${max_year}`);
 * ```
 */
export const fetchMetaYears = async (): Promise<MetaYearsResponse> => {
  if (API_MODE === 'mock') {
    await delay(MOCK_DELAY_MS);
    return metaYearsFixture;
  }
  
  return makeRequest<MetaYearsResponse>('/api/meta/years');
};

/**
 * Fetches the list of available countries with their data sources
 * 
 * @returns Promise resolving to array of countries with metadata
 * 
 * @example
 * ```typescript
 * const { countries } = await fetchCountries();
 * countries.forEach(c => console.log(`${c.name} (${c.code})`));
 * ```
 */
export const fetchCountries = async (): Promise<CountriesResponse> => {
  if (API_MODE === 'mock') {
    await delay(MOCK_DELAY_MS);
    return countriesFixture;
  }
  
  return makeRequest<CountriesResponse>('/api/meta/countries');
};

/**
 * Fetches a paginated list of gender-neutral names with optional filters
 * 
 * @param filters - Optional filter parameters for the query
 * @returns Promise resolving to paginated names list with metadata
 * 
 * @example
 * ```typescript
 * const response = await fetchNames({
 *   page: 1,
 *   page_size: 20,
 *   countries: ['US', 'UK'],
 *   year_min: 2000,
 *   sort_by: 'total_count',
 *   sort_order: 'desc'
 * });
 * ```
 */
export const fetchNames = async (
  filters: NamesFilterParams = {},
  signal?: AbortSignal
): Promise<NamesListResponse> => {
  if (API_MODE === 'mock') {
    await delay(MOCK_DELAY_MS);
    
    // Simulate empty results for specific filter combinations
    if (filters.search && filters.search.toLowerCase() === 'xyz123') {
      return namesListEmptyFixture;
    }
    
    return namesListFixture;
  }
  
  const queryString = buildQueryString(filters);
  return makeRequest<NamesListResponse>(`/api/names${queryString}`, signal);
};

/**
 * Fetches detailed trend data for a specific name
 * 
 * @param params - Parameters including name and optional filters
 * @returns Promise resolving to name trend data with time series and country breakdown
 * 
 * @example
 * ```typescript
 * const trend = await fetchNameTrend({
 *   name: 'Alex',
 *   countries: ['US', 'UK'],
 *   year_min: 1980,
 *   year_max: 2023
 * });
 * console.log(`${trend.name}: ${trend.summary.total_count} total occurrences`);
 * ```
 */
export const fetchNameTrend = async (
  params: NameTrendParams
): Promise<NameTrendResponse> => {
  if (API_MODE === 'mock') {
    await delay(MOCK_DELAY_MS);
    return nameDetailFixture;
  }
  
  const queryString = buildQueryString(params);
  return makeRequest<NameTrendResponse>(`/api/names/trend${queryString}`);
};

/**
 * Get current API mode
 * @returns 'mock' or 'real'
 */
export const getApiMode = (): string => API_MODE;

/**
 * Get API base URL (only relevant in real mode)
 * @returns API base URL
 */
export const getApiBaseUrl = (): string => API_BASE_URL;