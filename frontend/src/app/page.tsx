'use client';

import { useEffect, useState, useCallback, Suspense } from 'react';
import { useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { Plus } from 'lucide-react';
import DashboardHeader from '@/components/DashboardHeader';
import TopBreakdownsChart from '@/components/TopBreakdownsChart';
import MTBFMetrics from '@/components/MTBFMetrics';
import RecentDisruptionsTable from '@/components/RecentDisruptionsTable';
import { getTopBreakdowns, getMTBF, getRecentDisruptions } from '@/lib/api';
import type { DashboardData } from '@/types';

const REFRESH_INTERVAL = 30000;

function DashboardContent() {
  const searchParams = useSearchParams();
  const [data, setData] = useState<DashboardData>({
    topBreakdownsByLine: null,
    topBreakdownsByStation: null,
    mtbf: null,
    recentDisruptions: null,
  });
  const [isLoading, setIsLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);
  const [error, setError] = useState<string | null>(null);

  const fetchDashboardData = useCallback(async (showRefreshing = false) => {
    try {
      if (showRefreshing) {
        setIsRefreshing(true);
      }
      setError(null);

      const [topByLine, topByStation, mtbfData, disruptions] = await Promise.all([
        getTopBreakdowns('line', 5),
        getTopBreakdowns('station', 5),
        getMTBF(),
        getRecentDisruptions(undefined, undefined, 20),
      ]);

      setData({
        topBreakdownsByLine: topByLine,
        topBreakdownsByStation: topByStation,
        mtbf: mtbfData,
        recentDisruptions: disruptions,
      });
      setLastUpdated(new Date());
    } catch (err) {
      console.error('Failed to fetch dashboard data:', err);
      setError(
        err instanceof Error
          ? err.message
          : 'Failed to load dashboard data. Please try again.'
      );
    } finally {
      setIsLoading(false);
      setIsRefreshing(false);
    }
  }, []);

  useEffect(() => {
    // Check if we need to refresh (e.g., returning from submit-incident page)
    const shouldRefresh = searchParams.get('refresh') === 'true';
    fetchDashboardData(shouldRefresh);

    const interval = setInterval(() => {
      fetchDashboardData(true);
    }, REFRESH_INTERVAL);

    return () => clearInterval(interval);
  }, [fetchDashboardData, searchParams]);

  const handleManualRefresh = () => {
    fetchDashboardData(true);
  };

  if (error && !data.mtbf) {
    return (
      <div className="min-h-screen bg-gray-50">
        <DashboardHeader
          lastUpdated={lastUpdated}
          isRefreshing={isRefreshing}
          onRefresh={handleManualRefresh}
        />
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
            <h3 className="text-lg font-semibold text-red-800 mb-2">
              Error Loading Dashboard
            </h3>
            <p className="text-red-600">{error}</p>
            <button
              onClick={handleManualRefresh}
              className="mt-4 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <DashboardHeader
        lastUpdated={lastUpdated}
        isRefreshing={isRefreshing}
        onRefresh={handleManualRefresh}
      />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Quick Action Banner */}
        <div className="mb-8 bg-gradient-to-r from-blue-50 to-green-50 border border-blue-200 rounded-lg p-6">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-lg font-semibold text-gray-900 mb-1">
                Report a New Incident
              </h2>
              <p className="text-sm text-gray-600">
                Submit incident reports to help track transport reliability
              </p>
            </div>
            <Link
              href="/submit-incident"
              className="inline-flex items-center px-6 py-3 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition-colors"
            >
              <Plus className="h-5 w-5 mr-2" />
              Submit Incident
            </Link>
          </div>
        </div>

        <div className="space-y-8">
          <TopBreakdownsChart
            dataByLine={data.topBreakdownsByLine}
            dataByStation={data.topBreakdownsByStation}
            isLoading={isLoading}
          />

          <MTBFMetrics data={data.mtbf} isLoading={isLoading} />

          <RecentDisruptionsTable
            data={data.recentDisruptions}
            isLoading={isLoading}
          />
        </div>

        {error && (
          <div className="mt-4 bg-yellow-50 border border-yellow-200 rounded-lg p-4">
            <p className="text-sm text-yellow-800">
              <strong>Warning:</strong> {error}
            </p>
          </div>
        )}
      </main>
    </div>
  );
}

export default function Dashboard() {
  return (
    <Suspense fallback={
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading dashboard...</p>
        </div>
      </div>
    }>
      <DashboardContent />
    </Suspense>
  );
}
