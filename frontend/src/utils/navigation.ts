/**
 * Navigation helper utilities for type-safe routing
 */

/**
 * Build a URL for viewing a specific name's details
 */
export function buildNameDetailUrl(name: string): string {
  return `/name/${encodeURIComponent(name)}`;
}

/**
 * Build a URL with preserved query parameters
 */
export function buildUrlWithParams(
  path: string,
  params?: Record<string, string | number | boolean | undefined>
): string {
  if (!params || Object.keys(params).length === 0) {
    return path;
  }

  const searchParams = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined) {
      searchParams.append(key, String(value));
    }
  });

  const queryString = searchParams.toString();
  return queryString ? `${path}?${queryString}` : path;
}

/**
 * Get current query parameters as an object
 */
export function getQueryParams(search: string): Record<string, string> {
  const params: Record<string, string> = {};
  const searchParams = new URLSearchParams(search);
  
  searchParams.forEach((value, key) => {
    params[key] = value;
  });
  
  return params;
}

/**
 * Route paths as constants for type safety
 */
export const ROUTES = {
  HOME: '/',
  NAMES_EXPLORER: '/names',
  NAME_DETAIL: '/name/:name',
} as const;

/**
 * Build name detail path from route constant
 */
export function getNameDetailPath(name: string): string {
  return ROUTES.NAME_DETAIL.replace(':name', encodeURIComponent(name));
}