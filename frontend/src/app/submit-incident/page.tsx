'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { getLines, getStations, createIncident } from '@/lib/api';
import type { Line, Station } from '@/types';

export default function SubmitIncident() {
  const router = useRouter();
  const [lines, setLines] = useState<Line[]>([]);
  const [stations, setStations] = useState<Station[]>([]);
  const [filteredStations, setFilteredStations] = useState<Station[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const [formData, setFormData] = useState({
    line: '',
    station: '',
    date: '',
    time: '',
    duration_minutes: '',
    incident_type: 'mechanical',
  });

  useEffect(() => {
    async function fetchData() {
      try {
        const [linesData, stationsData] = await Promise.all([
          getLines(),
          getStations(),
        ]);
        setLines(linesData.lines || []);
        setStations(stationsData.stations || []);
      } catch (err) {
        setError('Failed to load form data');
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    }
    fetchData();
  }, []);

  useEffect(() => {
    if (formData.line) {
      const filtered = stations.filter((s) => s.lineName === formData.line);
      setFilteredStations(filtered);
      // Reset station if current selection is not in the new line
      if (formData.station && !filtered.some((s) => s.name === formData.station)) {
        setFormData((prev) => ({ ...prev, station: '' }));
      }
    } else {
      setFilteredStations([]);
    }
  }, [formData.line, stations, formData.station]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);

    try {
      // Combine date and time into ISO timestamp
      const timestamp = new Date(`${formData.date}T${formData.time}`).toISOString();

      await createIncident({
        line: formData.line,
        station: formData.station,
        timestamp,
        durationMinutes: parseInt(formData.duration_minutes, 10),
        incidentType: formData.incident_type,
      });

      setSuccess(true);
      // Reset form
      setFormData({
        line: '',
        station: '',
        date: '',
        time: '',
        duration_minutes: '',
        incident_type: 'mechanical',
      });

      // Redirect to dashboard after 2 seconds with cache refresh
      setTimeout(() => {
        router.push('/?refresh=true');
      }, 2000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to submit incident');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold text-gray-900">
              Submit Incident Report
            </h1>
            <button
              onClick={() => router.push('/')}
              className="text-blue-600 hover:text-blue-800 font-medium"
            >
              Back to Dashboard
            </button>
          </div>
        </div>
      </header>

      <main className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {success && (
          <div className="mb-6 bg-green-50 border border-green-200 rounded-lg p-4">
            <div className="flex items-center">
              <svg
                className="w-5 h-5 text-green-600 mr-2"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                  clipRule="evenodd"
                />
              </svg>
              <p className="text-green-800 font-medium">
                Incident submitted successfully! Redirecting to dashboard...
              </p>
            </div>
          </div>
        )}

        {error && (
          <div className="mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
            <div className="flex items-center">
              <svg
                className="w-5 h-5 text-red-600 mr-2"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  clipRule="evenodd"
                />
              </svg>
              <p className="text-red-800">{error}</p>
            </div>
          </div>
        )}

        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Line Selection */}
            <div>
              <label
                htmlFor="line"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Line <span className="text-red-500">*</span>
              </label>
              <select
                id="line"
                name="line"
                value={formData.line}
                onChange={handleInputChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">Select a line</option>
                {lines.map((line) => (
                  <option key={line.id} value={line.name}>
                    {line.name}
                  </option>
                ))}
              </select>
            </div>

            {/* Station Selection */}
            <div>
              <label
                htmlFor="station"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Station <span className="text-red-500">*</span>
              </label>
              <select
                id="station"
                name="station"
                value={formData.station}
                onChange={handleInputChange}
                required
                disabled={!formData.line}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
              >
                <option value="">
                  {formData.line ? 'Select a station' : 'Select a line first'}
                </option>
                {filteredStations.map((station) => (
                  <option key={station.id} value={station.name}>
                    {station.name}
                  </option>
                ))}
              </select>
            </div>

            {/* Date and Time */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label
                  htmlFor="date"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Date <span className="text-red-500">*</span>
                </label>
                <input
                  type="date"
                  id="date"
                  name="date"
                  value={formData.date}
                  onChange={handleInputChange}
                  max={new Date().toISOString().split('T')[0]}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              <div>
                <label
                  htmlFor="time"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Time <span className="text-red-500">*</span>
                </label>
                <input
                  type="time"
                  id="time"
                  name="time"
                  value={formData.time}
                  onChange={handleInputChange}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            </div>

            {/* Duration */}
            <div>
              <label
                htmlFor="duration_minutes"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Duration (minutes) <span className="text-red-500">*</span>
              </label>
              <input
                type="number"
                id="duration_minutes"
                name="duration_minutes"
                value={formData.duration_minutes}
                onChange={handleInputChange}
                min="0"
                max="1440"
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter duration in minutes (0-1440)"
              />
              <p className="mt-1 text-sm text-gray-500">
                Maximum 1440 minutes (24 hours)
              </p>
            </div>

            {/* Incident Type */}
            <div>
              <label
                htmlFor="incident_type"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Incident Type <span className="text-red-500">*</span>
              </label>
              <select
                id="incident_type"
                name="incident_type"
                value={formData.incident_type}
                onChange={handleInputChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="mechanical">Mechanical</option>
                <option value="signal">Signal</option>
                <option value="power">Power</option>
                <option value="weather">Weather</option>
                <option value="other">Other</option>
              </select>
            </div>

            {/* Submit Button */}
            <div className="pt-4">
              <button
                type="submit"
                disabled={isSubmitting}
                className="w-full bg-blue-600 text-white py-3 px-4 rounded-md font-medium hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:bg-blue-400 disabled:cursor-not-allowed transition-colors"
              >
                {isSubmitting ? (
                  <span className="flex items-center justify-center">
                    <svg
                      className="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <circle
                        className="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        strokeWidth="4"
                      ></circle>
                      <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                      ></path>
                    </svg>
                    Submitting...
                  </span>
                ) : (
                  'Submit Incident Report'
                )}
              </button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}