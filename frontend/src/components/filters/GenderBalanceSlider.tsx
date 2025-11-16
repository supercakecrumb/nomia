/**
 * GenderBalanceSlider Component
 * 
 * Dual-handle range slider for selecting gender balance range.
 * Scale: 0-100 (0 = all female, 100 = all male).
 * Shows labels and visual indicators.
 */

import { useTranslation } from 'react-i18next';
import Slider from 'rc-slider';
import 'rc-slider/assets/index.css';

interface GenderBalanceSliderProps {
  value: [number, number];
  onChange: (value: [number, number]) => void;
}

export default function GenderBalanceSlider({
  value,
  onChange,
}: GenderBalanceSliderProps) {
  const { t } = useTranslation('filters');

  return (
    <div className="space-y-3">
      <label className="flex items-center gap-2 text-sm font-semibold text-gray-700">
        <svg className="w-4 h-4 text-accent-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 6l3 1m0 0l-3 9a5.002 5.002 0 006.001 0M6 7l3 9M6 7l6-2m6 2l3-1m-3 1l-3 9a5.002 5.002 0 006.001 0M18 7l3 9m-3-9l-6-2m0-2v2m0 16V5m0 16H9m3 0h3" />
        </svg>
        {t('genderBalance.label')}
      </label>

      <div className="px-2">
        <Slider
          range
          min={0}
          max={100}
          value={value}
          onChange={(val) => onChange(val as [number, number])}
          trackStyle={[{ backgroundColor: '#8b5cf6', height: 6 }]}
          handleStyle={[
            {
              backgroundColor: '#ec4899',
              borderColor: '#ec4899',
              width: 20,
              height: 20,
              marginTop: -7,
              boxShadow: '0 2px 8px rgba(236, 72, 153, 0.3)',
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

      <div className="flex items-center justify-between text-sm pt-4">
        <div className="flex items-center gap-2">
          <span className="text-pink-600 font-semibold">♀</span>
          <span className="font-medium text-gray-700">{t('genderBalance.moreFemale')}</span>
        </div>
        <div className="flex items-center gap-2">
          <span className="font-semibold text-gray-900 bg-accent-50 px-3 py-1 rounded-lg border border-accent-200">
            {value[0]}% - {value[1]}%
          </span>
        </div>
        <div className="flex items-center gap-2">
          <span className="font-medium text-gray-700">{t('genderBalance.moreMale')}</span>
          <span className="text-blue-600 font-semibold">♂</span>
        </div>
      </div>

      <p className="text-xs text-gray-500 flex items-center gap-1">
        <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        0 = {t('genderBalance.female')}, 50 = {t('genderBalance.neutral')}, 100 = {t('genderBalance.male')}
      </p>
    </div>
  );
}