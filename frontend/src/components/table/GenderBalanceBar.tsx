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
        <div className="flex-1 h-6 bg-gray-200 rounded-md flex items-center justify-center">
          <span className="text-xs text-gray-600 font-medium">Unknown</span>
        </div>
      </div>
    );
  }

  const femalePercent = 100 - genderBalance;
  const malePercent = genderBalance;

  return (
    <div className="space-y-1">
      {/* Visual Bar */}
      <div className="relative h-6 bg-gradient-to-r from-pink-200 via-purple-100 to-blue-200 rounded-md overflow-hidden">
        {/* Indicator Dot */}
        <div
          className="absolute top-1/2 -translate-y-1/2 w-3 h-3 bg-gray-800 rounded-full border-2 border-white shadow-md transition-all"
          style={{ left: `${genderBalance}%` }}
        />
      </div>

      {/* Percentages */}
      <div className="flex justify-between text-xs text-gray-600">
        <span className="font-medium">
          ♀ {formatPercentage(femalePercent, 0)}
        </span>
        <span className="font-medium">
          ♂ {formatPercentage(malePercent, 0)}
        </span>
      </div>
    </div>
  );
}