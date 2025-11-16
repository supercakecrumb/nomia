/**
 * GenderBalanceBar Component
 * 
 * Visual representation of gender balance for a name.
 * Shows a horizontal bar with female on left, male on right, and a dot indicator.
 */

import { formatPercentage } from '../../utils/formatters';

interface GenderBalanceBarProps {
  /** Gender balance value (0-100, where 0=100% female, 100=100% male) */
  genderBalance: number | null;
  /** Female count */
  femaleCount: number;
  /** Male count */
  maleCount: number;
  /** Whether the name has unknown gender data */
  hasUnknownData?: boolean;
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
        <div className="flex-1 h-7 bg-gradient-to-r from-gray-200 to-gray-300 rounded-lg flex items-center justify-center border border-gray-300 shadow-sm">
          <span className="text-xs text-gray-600 font-semibold">Unknown</span>
        </div>
      </div>
    );
  }

  const femalePercent = 100 - genderBalance;
  const malePercent = genderBalance;

  return (
    <div className="space-y-2">
      {/* Visual Bar with Soft Gradient */}
      <div className="relative h-7 rounded-lg overflow-hidden shadow-sm border border-gray-200">
        {/* Soft base gradient background */}
        <div
          className="absolute inset-0 transition-all duration-300"
          style={{
            background: `linear-gradient(to right,
              #FFB6D9 0%,
              #F5C6E8 25%,
              #E6D9FF 50%,
              #D6E8F5 75%,
              #B6D9FF 100%)`
          }}
        ></div>
        
        {/* Gradient Overlay for Depth */}
        <div className="absolute inset-0 bg-gradient-to-b from-white/20 to-transparent"></div>
        
        {/* Female Side Accent (soft pink) */}
        <div
          className="absolute left-0 top-0 h-full opacity-50 transition-all duration-300"
          style={{
            width: `${femalePercent}%`,
            background: 'linear-gradient(to right, #FFB6D9, #FFCCE5)'
          }}
        ></div>
        
        {/* Male Side Accent (soft blue) */}
        <div
          className="absolute right-0 top-0 h-full opacity-50 transition-all duration-300"
          style={{
            width: `${malePercent}%`,
            background: 'linear-gradient(to left, #B6D9FF, #CCE5FF)'
          }}
        ></div>

        {/* Center Neutral Zone (lavender) */}
        <div
          className="absolute top-0 h-full opacity-30 transition-all duration-300"
          style={{
            left: `${Math.max(0, genderBalance - 10)}%`,
            width: '20%',
            background: 'linear-gradient(to right, #E6D9FF, #F0E6FF, #E6D9FF)'
          }}
        ></div>
        
        {/* Position Indicator Dot */}
        <div
          className="absolute top-1/2 -translate-y-1/2 w-4 h-4 bg-gray-900 rounded-full border-2 border-white shadow-lg transition-all duration-300 hover:scale-125 z-10"
          style={{ left: `calc(${genderBalance}% - 8px)` }}
          title={`Gender Balance: ${genderBalance.toFixed(1)}%`}
        >
          {/* Inner glow effect */}
          <div className="absolute inset-0.5 bg-gradient-to-br from-white/50 to-transparent rounded-full"></div>
        </div>
      </div>

      {/* Percentages with Icons - softer colors */}
      <div className="flex justify-between items-center text-xs">
        <div className="flex items-center gap-1.5 px-2 py-1 rounded-md border" style={{ backgroundColor: '#FFE5F0', borderColor: '#FFB6D9' }}>
          <span className="font-bold" style={{ color: '#D9006C' }}>♀</span>
          <span className="font-semibold" style={{ color: '#A60052' }}>
            {formatPercentage(femalePercent, 0)}
          </span>
        </div>
        
        <div className="flex items-center gap-1 px-2 py-0.5 rounded-md border" style={{ backgroundColor: '#F5F0FF', borderColor: '#E6D9FF' }}>
          <span className="text-xs font-medium" style={{ color: '#8B6FBD' }}>Neutral</span>
        </div>
        
        <div className="flex items-center gap-1.5 px-2 py-1 rounded-md border" style={{ backgroundColor: '#E5F2FF', borderColor: '#B6D9FF' }}>
          <span className="font-bold" style={{ color: '#006CD9' }}>♂</span>
          <span className="font-semibold" style={{ color: '#0052A6' }}>
            {formatPercentage(malePercent, 0)}
          </span>
        </div>
      </div>
    </div>
  );
}