/**
 * API Response Types
 * 
 * Type definitions for all API responses matching the backend contract.
 */

/**
 * Meta years response - provides the available year range in the database
 */
export interface MetaYearsResponse {
  min_year: number;
  max_year: number;
}

/**
 * Country information with data source details
 */
export interface Country {
  code: string;
  name: string;
  data_source_name: string;
  data_source_url: string;
  data_source_description: string | null;
  data_source_requires_manual_download: boolean;
}

/**
 * Countries response
 */
export interface CountriesResponse {
  countries: Country[];
}

/**
 * Pagination metadata
 */
export interface PaginationMeta {
  page: number;
  page_size: number;
  total_count: number;
  total_pages: number;
  db_start: number;
  db_end: number;
}

/**
 * Popularity summary statistics
 */
export interface PopularitySummary {
  population_total: number;
  active_driver: string | null;
  active_value: number | null;
  derived_min_count: number;
  derived_top_n: number;
  derived_coverage_percent: number;
}

/**
 * Name list metadata including pagination and popularity summary
 */
export interface NamesListMeta extends PaginationMeta {
  popularity_summary: PopularitySummary;
}

/**
 * Name entry in the list
 */
export interface NameEntry {
  name: string;
  total_count: number;
  female_count: number;
  male_count: number;
  unknown_count: number;
  gender_balance: number | null;
  has_unknown_data: boolean;
  rank: number;
  cumulative_share: number;
  name_start: number;
  name_end: number;
  countries: string[];
}

/**
 * Names list response with pagination and names array
 */
export interface NamesListResponse {
  meta: NamesListMeta;
  names: NameEntry[];
}

/**
 * Time series data point for a name's trend over time
 */
export interface TimeSeriesPoint {
  year: number;
  total_count: number;
  female_count: number;
  male_count: number;
  unknown_count: number;
  gender_balance: number | null;
}

/**
 * Country-specific statistics for a name
 */
export interface CountryStats {
  country_code: string;
  country_name: string;
  total_count: number;
  female_count: number;
  male_count: number;
  unknown_count: number;
  gender_balance: number | null;
}

/**
 * Summary statistics for a name
 */
export interface NameSummary {
  total_count: number;
  female_count: number;
  male_count: number;
  unknown_count: number;
  gender_balance: number | null;
  has_unknown_data: boolean;
  name_start: number;
  name_end: number;
  countries: string[];
}

/**
 * Basic metadata for name trend response
 */
export interface NameTrendMeta {
  db_start: number;
  db_end: number;
}

/**
 * Name trend response with detailed statistics and time series
 */
export interface NameTrendResponse {
  name: string;
  meta: NameTrendMeta;
  summary: NameSummary;
  time_series: TimeSeriesPoint[];
  by_country: CountryStats[];
}

/**
 * API error response
 */
export interface ApiError {
  code: string;
  message: string;
}

/**
 * Filter parameters for names list endpoint
 */
export interface NamesFilterParams {
  page?: number;
  page_size?: number;
  countries?: string[];
  year_min?: number;
  year_max?: number;
  gender_balance_min?: number;
  gender_balance_max?: number;
  min_count?: number;
  top_n?: number;
  coverage_percent?: number;
  sort_by?: 'name' | 'total_count' | 'gender_balance' | 'rank';
  sort_order?: 'asc' | 'desc';
  search?: string;
}

/**
 * Parameters for name trend endpoint
 */
export interface NameTrendParams {
  name: string;
  countries?: string[];
  year_min?: number;
  year_max?: number;
}