/**
 * YearRangeSlider Component
 * 
 * Dual-handle range slider for selecting year range.
 * Shows min/max values from useMetaYears (e.g., 1880-2024).
 */

import { useTranslation } from 'react-i18next';
import Slider from 'rc-slider';
import 'rc-slider/assets/index.css';

interface YearRangeSliderProps {
  min: number;
  max: number;
  value: [number, number];
  onChange: (value: [number, number]) => void;
}

export default function YearRangeSlider({
  min,
  max,
  value,
  onChange,
}: YearRangeSliderProps) {
  const { t } = useTranslation('filters');

  return (
    <div className="space-y-3">
      <label className="flex items-center gap-2 text-sm font-semibold text-gray-700">
        <svg className="w-4 h-4 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        {t('yearRange.label')}
      </label>

      <div className="px-2">
        <Slider
          range
          min={min}
          max={max}
          value={value}
          onChange={(val) => onChange(val as [number, number])}
          trackStyle={[{ backgroundColor: '#3b82f6', height: 6 }]}
          handleStyle={[
            {
              backgroundColor: '#3b82f6',
              borderColor: '#3b82f6',
              width: 20,
              height: 20,
              marginTop: -7,
              boxShadow: '0 2px 8px rgba(59, 130, 246, 0.3)',
            },
            {
              backgroundColor: '#3b82f6',
              borderColor: '#3b82f6',
              width: 20,
              height: 20,
              marginTop: -7,
              boxShadow: '0 2px 8px rgba(59, 130, 246, 0.3)',
            },
          ]}
          railStyle={{ backgroundColor: '#e5e7eb', height: 6 }}
        />
      </div>

      <div className="flex items-center justify-between text-sm">
        <span className="font-semibold text-gray-900 bg-primary-50 px-3 py-1 rounded-lg border border-primary-200">
          {value[0]}
        </span>
        <span className="text-gray-400 font-medium">—</span>
        <span className="font-semibold text-gray-900 bg-primary-50 px-3 py-1 rounded-lg border border-primary-200">
          {value[1]}
        </span>
      </div>

      <p className="text-xs text-gray-500 flex items-center gap-1">
        <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Available: {min} - {max}
      </p>
    </div>
  );
}