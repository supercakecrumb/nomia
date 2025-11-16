import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  Brush,
} from 'recharts';
import type { TimeSeriesPoint } from '../../types/api';

interface Props {
  data: TimeSeriesPoint[];
  title?: string;
}

/**
 * TrendChart Component
 * 
 * Displays name popularity trends over time using Recharts AreaChart.
 * Features:
 * - Two area series for female and male counts
 * - Brush component for zooming and panning the time window
 * - Smart tick formatting for readability with 200+ years of data
 * - Responsive container that adapts to parent width
 */
export function TrendChart({ data, title = 'Popularity Trend Over Time' }: Props) {
  // Transform data to match Recharts format
  const chartData = data.map(point => ({
    year: point.year,
    female: point.female_count,
    male: point.male_count,
  }));

  // Custom tooltip formatter
  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="bg-white border border-gray-200 rounded-lg shadow-lg p-3">
          <p className="font-semibold text-gray-900 mb-2">{label}</p>
          {payload.map((entry: any, index: number) => (
            <p key={index} className="text-sm" style={{ color: entry.color }}>
              <span className="font-medium">{entry.name}:</span>{' '}
              {entry.value.toLocaleString()}
            </p>
          ))}
          <p className="text-sm text-gray-600 mt-1 pt-1 border-t">
            <span className="font-medium">Total:</span>{' '}
            {(payload[0].value + payload[1].value).toLocaleString()}
          </p>
        </div>
      );
    }
    return null;
  };

  return (
    <div className="rounded-2xl bg-white p-6 md:p-8 border border-gray-100">
      <h3 className="text-xl md:text-2xl font-bold text-gray-900 mb-6">
        {title}
      </h3>

      <ResponsiveContainer width="100%" height={400}>
        <AreaChart 
          data={chartData} 
          margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
        >
          <defs>
            <linearGradient id="colorFemale" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#ec4899" stopOpacity={0.3}/>
              <stop offset="95%" stopColor="#ec4899" stopOpacity={0}/>
            </linearGradient>
            <linearGradient id="colorMale" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3}/>
              <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
            </linearGradient>
          </defs>
          
          <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
          
          <XAxis
            dataKey="year"
            tickFormatter={(year) => (year % 10 === 0 ? year.toString() : '')}
            interval={0}
            tick={{ fill: '#6b7280', fontSize: 12 }}
            stroke="#d1d5db"
          />
          
          <YAxis
            tickFormatter={(value) => value.toLocaleString()}
            tick={{ fill: '#6b7280', fontSize: 12 }}
            stroke="#d1d5db"
          />
          
          <Tooltip content={<CustomTooltip />} />
          
          <Legend 
            wrapperStyle={{ paddingTop: '20px' }}
            iconType="circle"
          />

          <Area
            type="monotone"
            dataKey="female"
            name="Girls"
            stroke="#ec4899"
            strokeWidth={2}
            fill="url(#colorFemale)"
          />
          
          <Area
            type="monotone"
            dataKey="male"
            name="Boys"
            stroke="#3b82f6"
            strokeWidth={2}
            fill="url(#colorMale)"
          />

          <Brush
            dataKey="year"
            height={30}
            stroke="#6366f1"
            fill="#f3f4f6"
            travellerWidth={10}
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}