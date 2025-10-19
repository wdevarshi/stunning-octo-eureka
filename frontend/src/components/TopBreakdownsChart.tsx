'use client';

import { useState } from 'react';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from 'recharts';
import type { TopBreakdownsResponse } from '@/types';

interface TopBreakdownsChartProps {
  dataByLine: TopBreakdownsResponse | null;
  dataByStation: TopBreakdownsResponse | null;
  isLoading?: boolean;
}

const COLORS = ['#3b82f6', '#60a5fa', '#93c5fd', '#bfdbfe', '#dbeafe'];

export default function TopBreakdownsChart({
  dataByLine,
  dataByStation,
  isLoading,
}: TopBreakdownsChartProps) {
  const [view, setView] = useState<'line' | 'station'>('line');

  const currentData = view === 'line' ? dataByLine : dataByStation;

  if (isLoading) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-4"></div>
          <div className="h-64 bg-gray-100 rounded"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-semibold text-gray-900">
          Top Breakdowns by {view === 'line' ? 'Line' : 'Station'}
        </h2>
        <div className="flex space-x-2">
          <button
            onClick={() => setView('line')}
            className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
              view === 'line'
                ? 'bg-primary-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            By Line
          </button>
          <button
            onClick={() => setView('station')}
            className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
              view === 'station'
                ? 'bg-primary-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            By Station
          </button>
        </div>
      </div>

      {currentData && currentData.items.length > 0 ? (
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={currentData.items}>
            <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
            <XAxis dataKey="name" stroke="#6b7280" />
            <YAxis stroke="#6b7280" />
            <Tooltip
              contentStyle={{
                backgroundColor: '#fff',
                border: '1px solid #e5e7eb',
                borderRadius: '0.375rem',
              }}
            />
            <Bar dataKey="count" radius={[8, 8, 0, 0]}>
              {currentData.items.map((_, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      ) : (
        <div className="h-64 flex items-center justify-center text-gray-500">
          No data available
        </div>
      )}
    </div>
  );
}
