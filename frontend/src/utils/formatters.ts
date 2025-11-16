/**
 * Utility Formatters
 * 
 * Common formatting functions for numbers, percentages, and dates.
 */

/**
 * Format a number with thousand separators
 * @param num - The number to format
 * @returns Formatted string with commas (e.g., "1,234,567")
 */
export function formatNumber(num: number): string {
  return num.toLocaleString('en-US');
}

/**
 * Format a number as a percentage
 * @param num - The number to format (0-100)
 * @param decimals - Number of decimal places (default: 1)
 * @returns Formatted percentage string (e.g., "45.5%")
 */
export function formatPercentage(num: number, decimals: number = 1): string {
  return `${num.toFixed(decimals)}%`;
}

/**
 * Format a year range
 * @param start - Start year
 * @param end - End year
 * @returns Formatted year range (e.g., "1990-2023" or "2020" if same)
 */
export function formatYearRange(start: number, end: number): string {
  return start === end ? `${start}` : `${start}-${end}`;
}