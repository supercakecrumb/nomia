/**
 * GenderBalanceBar Component
 * 
 * Visual representation of gender balance for a name.
 * Shows a smooth gradient bar with a dot indicator and percentage labels.
 */

import { formatPercentage } from '../../utils/formatters';

interface GenderBalanceBarProps {
  /** Gender balance value (0–100, where 0 = 100% female, 100 = 100% male) */
  genderBalance: number | null;
  /** Female count */
  femaleCount: number;
  /** Male count */
  maleCount: number;
  /** Whether the name has unknown gender data */
  hasUnknownData?: boolean;
}

function clampPercent(value: number): number {
  if (Number.isNaN(value)) return 50;
  return Math.min(100, Math.max(0, value));
}

export default function GenderBalanceBar({
  genderBalance,
  femaleCount,
  maleCount,
  hasUnknownData = false,
}: GenderBalanceBarProps) {
  // Handle null/unknown gender balance
  if (genderBalance === null || hasUnknownData) {
    return (
      <div className="flex items-center gap-2">
        <div className="flex-1 h-7 rounded-full bg-gradient-to-r from-gray-200 to-gray-300 shadow-sm flex items-center justify-center">
          <span className="text-xs text-gray-600 font-semibold">Unknown</span>
        </div>
      </div>
    );
  }

  const clamped = clampPercent(genderBalance);
  const femalePercent = 100 - clamped;
  const malePercent = clamped;

  const ariaLabel = `Gender balance: ${formatPercentage(
    femalePercent,
    0,
  )} female and ${formatPercentage(malePercent, 0)} male.`;

  return (
    <div className="space-y-2">
      {/* Smooth gradient bar with no internal borders */}
      <div
        className="relative h-7 rounded-full overflow-hidden shadow-[0_10px_30px_rgba(157,178,255,0.35)]"
        role="img"
        aria-label={ariaLabel}
      >
        {/* Single smooth gradient - no overlays to create seams */}
        <div
          aria-hidden="true"
          className="absolute inset-0"
          style={{
            background: `linear-gradient(
              90deg,
              #ffc3dd 0%,
              #f7c8e6 20%,
              #ecd4f1 40%,
              #e3dcff 50%,
              #d3e3ff 70%,
              #bfd9ff 100%
            )`,
          }}
        />

        {/* Soft vertical highlight (full-width, no edges) */}
        <div
          aria-hidden="true"
          className="absolute inset-0 pointer-events-none"
          style={{
            background:
              'radial-gradient(circle at 50% 0%, rgba(255,255,255,0.55) 0, rgba(255,255,255,0.15) 45%, transparent 70%)',
            mixBlendMode: 'soft-light',
          }}
        />

        {/* Indicator dot */}
        <button
          type="button"
          className="absolute top-1/2 -translate-y-1/2 -translate-x-1/2 w-5 h-5 rounded-full bg-white shadow-[0_0_0_3px_rgba(255,255,255,0.9),0_4px_10px_rgba(0,0,0,0.25)] flex items-center justify-center transition-transform duration-200 ease-out hover:scale-110 focus-visible:scale-110 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-1 focus-visible:ring-indigo-300"
          style={{ left: `${clamped}%` }}
          title={`Gender balance: ${clamped.toFixed(1)}% male`}
        >
          <span className="block w-3 h-3 rounded-full bg-gray-900" />
        </button>
      </div>

      {/* Percentages with Icons */}
      <div className="flex justify-between items-center text-xs">
        <div
          className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl"
          style={{ backgroundColor: '#FFE5F0' }}
        >
          <span className="font-bold" style={{ color: '#D9006C' }}>
            ♀
          </span>
          <span className="font-semibold" style={{ color: '#A60052' }}>
            {formatPercentage(femalePercent, 0)}
          </span>
        </div>

        <div
          className="px-3 py-1.5 rounded-xl font-medium"
          style={{ backgroundColor: '#F5F0FF', color: '#8B6FBD' }}
        >
          Neutral
        </div>

        <div
          className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl"
          style={{ backgroundColor: '#E5F2FF' }}
        >
          <span className="font-bold" style={{ color: '#006CD9' }}>
            ♂
          </span>
          <span className="font-semibold" style={{ color: '#0052A6' }}>
            {formatPercentage(malePercent, 0)}
          </span>
        </div>
      </div>
    </div>
  );
}