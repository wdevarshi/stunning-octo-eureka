'use client';

import { Clock } from 'lucide-react';
import type { MTBFResponse } from '@/types';

interface MTBFMetricsProps {
  data: MTBFResponse | null;
  isLoading?: boolean;
}

function getColorClass(mtbfMinutes: number): string {
  const mtbfHours = mtbfMinutes / 60;
  if (mtbfHours >= 720) return 'bg-green-50 border-green-200 text-green-800';
  if (mtbfHours >= 360) return 'bg-yellow-50 border-yellow-200 text-yellow-800';
  return 'bg-red-50 border-red-200 text-red-800';
}

function formatMTBF(minutes: number): string {
  const hours = minutes / 60;
  if (hours >= 24) {
    const days = Math.floor(hours / 24);
    const remainingHours = Math.round(hours % 24);
    return `${days}d ${remainingHours}h`;
  }
  return `${Math.round(hours)}h`;
}

export default function MTBFMetrics({ data, isLoading }: MTBFMetricsProps) {
  if (isLoading) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/3 mb-4"></div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <div key={i} className="h-24 bg-gray-100 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  const sortedLines = data?.lines
    ? [...data.lines].sort((a, b) => b.mtbfMinutes - a.mtbfMinutes)
    : [];

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-xl font-semibold text-gray-900 mb-6 flex items-center">
        <Clock className="h-5 w-5 mr-2" />
        Mean Time Between Failures (MTBF)
      </h2>

      {sortedLines.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {sortedLines.map((line) => (
            <div
              key={line.name}
              className={`border-2 rounded-lg p-4 transition-all hover:shadow-md ${getColorClass(
                line.mtbfMinutes
              )}`}
            >
              <div className="text-sm font-medium mb-1">{line.name}</div>
              <div className="text-2xl font-bold">
                {formatMTBF(line.mtbfMinutes)}
              </div>
              <div className="text-xs mt-1 opacity-75">
                {Math.round(line.mtbfMinutes)} minutes
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="h-48 flex items-center justify-center text-gray-500">
          No MTBF data available
        </div>
      )}
    </div>
  );
}
