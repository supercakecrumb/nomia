/**
 * CountriesDropdown Component
 * 
 * Compact, elegant multi-select dropdown for countries.
 * Shows country flags (emoji), selected count when collapsed.
 * Pretty checkboxes when expanded.
 */

import { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import type { Country } from '../../types/api';

interface CountriesDropdownProps {
  countries: Country[];
  selectedCountries: string[];
  onChange: (selected: string[]) => void;
}

// Country code to flag emoji mapping
const countryFlags: Record<string, string> = {
  'US': '🇺🇸',
  'UK': '🇬🇧',
  'CA': '🇨🇦',
  'AU': '🇦🇺',
  'IE': '🇮🇪',
  'NZ': '🇳🇿',
  'SE': '🇸🇪',
  'NO': '🇳🇴',
  'DK': '🇩🇰',
  'FI': '🇫🇮',
};

export default function CountriesDropdown({
  countries,
  selectedCountries,
  onChange,
}: CountriesDropdownProps) {
  const { t } = useTranslation('filters');
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const toggleCountry = (code: string) => {
    if (selectedCountries.includes(code)) {
      onChange(selectedCountries.filter(c => c !== code));
    } else {
      onChange([...selectedCountries, code]);
    }
  };

  const getButtonLabel = () => {
    if (selectedCountries.length === 0) {
      return t('countries.all');
    }
    
    if (selectedCountries.length <= 2) {
      return selectedCountries.map(code => countryFlags[code] || code).join(' ');
    }
    
    const firstTwo = selectedCountries.slice(0, 2).map(code => countryFlags[code] || code).join(' ');
    return `${firstTwo} +${selectedCountries.length - 2}`;
  };

  return (
    <div className="space-y-2" ref={dropdownRef}>
      <label className="flex items-center gap-2 text-sm font-semibold text-gray-700">
        <svg className="w-4 h-4 text-secondary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {t('countries.label')}
      </label>

      <div className="relative">
        <button
          type="button"
          onClick={() => setIsOpen(!isOpen)}
          className="w-full px-4 py-2.5 bg-white border border-gray-300 rounded-lg hover:border-gray-400 focus:ring-2 focus:ring-secondary-500 focus:border-transparent transition-all duration-200 flex items-center justify-between"
        >
          <span className="flex items-center gap-2 text-sm font-medium text-gray-900">
            {getButtonLabel()}
            {selectedCountries.length > 0 && (
              <span className="ml-2 px-2 py-0.5 bg-secondary-100 text-secondary-700 rounded-full text-xs font-semibold">
                {selectedCountries.length}
              </span>
            )}
          </span>
          <svg
            className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${isOpen ? 'rotate-180' : ''}`}
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
        </button>

        {isOpen && (
          <div className="absolute z-50 mt-2 w-full bg-white border border-gray-200 rounded-lg shadow-lg max-h-64 overflow-y-auto">
            <div className="p-2 space-y-1">
              {countries.map((country) => (
                <label
                  key={country.code}
                  className="flex items-center gap-3 px-3 py-2 hover:bg-gray-50 rounded-lg cursor-pointer transition-colors group"
                >
                  <input
                    type="checkbox"
                    checked={selectedCountries.includes(country.code)}
                    onChange={() => toggleCountry(country.code)}
                    className="w-4 h-4 text-secondary-600 border-gray-300 rounded focus:ring-2 focus:ring-secondary-500 cursor-pointer"
                  />
                  <span className="text-xl">{countryFlags[country.code] || '🌐'}</span>
                  <div className="flex-1 min-w-0">
                    <div className="text-sm font-medium text-gray-900 group-hover:text-secondary-700 transition-colors">
                      {country.name}
                    </div>
                    <div className="text-xs text-gray-500">{country.code}</div>
                  </div>
                  {selectedCountries.includes(country.code) && (
                    <svg className="w-5 h-5 text-secondary-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                  )}
                </label>
              ))}
            </div>

            {selectedCountries.length > 0 && (
              <div className="border-t border-gray-200 p-2">
                <button
                  type="button"
                  onClick={() => onChange([])}
                  className="w-full px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 rounded-lg transition-colors"
                >
                  Clear all
                </button>
              </div>
            )}
          </div>
        )}
      </div>

      <p className="text-xs text-gray-500 flex items-center gap-1">
        <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {selectedCountries.length === 0 ? 'All countries selected' : `${selectedCountries.length} of ${countries.length} selected`}
      </p>
    </div>
  );
}