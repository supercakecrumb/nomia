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
      {/* Visual Bar with Enhanced Gradient */}
      <div className="relative h-7 bg-gradient-to-r from-pink-200 via-purple-200 to-blue-200 rounded-lg overflow-hidden shadow-sm border border-gray-200">
        {/* Gradient Overlay for Depth */}
        <div className="absolute inset-0 bg-gradient-to-b from-white/20 to-transparent"></div>
        
        {/* Female Side Indicator */}
        <div 
          className="absolute left-0 top-0 h-full bg-gradient-to-r from-pink-400 to-pink-300 opacity-60 transition-all duration-300"
          style={{ width: `${femalePercent}%` }}
        ></div>
        
        {/* Male Side Indicator */}
        <div 
          className="absolute right-0 top-0 h-full bg-gradient-to-l from-blue-400 to-blue-300 opacity-60 transition-all duration-300"
          style={{ width: `${malePercent}%` }}
        ></div>

        {/* Center Neutral Zone */}
        <div 
          className="absolute top-0 h-full bg-gradient-to-r from-purple-300 via-purple-400 to-purple-300 opacity-40"
          style={{ 
            left: `${Math.max(0, genderBalance - 10)}%`,
            width: '20%'
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

      {/* Percentages with Icons */}
      <div className="flex justify-between items-center text-xs">
        <div className="flex items-center gap-1.5 px-2 py-1 bg-pink-50 rounded-md border border-pink-200">
          <span className="text-pink-700 font-bold">♀</span>
          <span className="font-semibold text-pink-900">
            {formatPercentage(femalePercent, 0)}
          </span>
        </div>
        
        <div className="flex items-center gap-1 px-2 py-0.5 bg-purple-50 rounded-md border border-purple-200">
          <span className="text-xs text-purple-600 font-medium">Neutral</span>
        </div>
        
        <div className="flex items-center gap-1.5 px-2 py-1 bg-blue-50 rounded-md border border-blue-200">
          <span className="text-blue-700 font-bold">♂</span>
          <span className="font-semibold text-blue-900">
            {formatPercentage(malePercent, 0)}
          </span>
        </div>
      </div>
    </div>
  );
}