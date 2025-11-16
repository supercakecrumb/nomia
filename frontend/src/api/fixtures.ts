/**
 * API Fixtures Module
 * 
 * Imports and exports JSON fixture files from spec-examples directory
 * for use in mock mode API client.
 */

import type {
  MetaYearsResponse,
  CountriesResponse,
  NamesListResponse,
  NameTrendResponse,
} from '../types/api';

// Import JSON fixtures using the @spec-examples path alias
import metaYearsJson from '@spec-examples/meta-years.json';
import countriesJson from '@spec-examples/countries.json';
import namesListJson from '@spec-examples/names-list.json';
import nameDetailJson from '@spec-examples/name-detail.json';
import namesListEmptyJson from '@spec-examples/names-list-empty.json';

/**
 * Meta years fixture - provides available year range
 */
export const metaYearsFixture: MetaYearsResponse = metaYearsJson;

/**
 * Countries fixture - list of available countries with data sources
 */
export const countriesFixture: CountriesResponse = countriesJson;

/**
 * Names list fixture - paginated list of gender-neutral names
 */
export const namesListFixture: NamesListResponse = namesListJson;

/**
 * Name detail fixture - detailed trend data for a specific name
 */
export const nameDetailFixture: NameTrendResponse = nameDetailJson;

/**
 * Empty names list fixture - response when no names match filters
 */
export const namesListEmptyFixture: NamesListResponse = namesListEmptyJson;