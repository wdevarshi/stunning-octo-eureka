import type {
  TopBreakdownsResponse,
  MTBFResponse,
  RecentDisruptionsResponse,
  LinesResponse,
  StationsResponse,
  CreateIncidentRequest,
  CreateIncidentResponse,
} from '@/types';

// Use different API URLs for server-side (inside Docker) vs client-side (browser)
const API_BASE_URL = typeof window === 'undefined'
  ? (process.env.API_URL || 'http://nginx:8080') // Server-side: use Docker network
  : (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'); // Client-side: use host machine

class APIError extends Error {
  constructor(message: string, public status?: number) {
    super(message);
    this.name = 'APIError';
  }
}

async function fetchAPI<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      cache: 'no-store',
      headers: {
        'Content-Type': 'application/json',
      },
      ...options,
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new APIError(
        errorData.message || `API request failed: ${response.statusText}`,
        response.status
      );
    }

    const data = await response.json();
    return data;
  } catch (error) {
    if (error instanceof APIError) {
      throw error;
    }
    throw new APIError('Network error: Unable to connect to the API');
  }
}

export async function getTopBreakdowns(
  scope: 'line' | 'station' = 'line',
  limit: number = 5
): Promise<TopBreakdownsResponse> {
  return fetchAPI<TopBreakdownsResponse>(
    `/analytics/top_breakdowns?scope=${scope}&limit=${limit}`
  );
}

export async function getMTBF(): Promise<MTBFResponse> {
  return fetchAPI<MTBFResponse>('/analytics/mean_time_between_failures');
}

export async function getRecentDisruptions(
  line?: string,
  station?: string,
  limit: number = 20
): Promise<RecentDisruptionsResponse> {
  const params = new URLSearchParams();
  if (line) params.append('line', line);
  if (station) params.append('station', station);
  params.append('limit', limit.toString());

  const query = params.toString();
  return fetchAPI<RecentDisruptionsResponse>(
    `/analytics/recent_disruptions${query ? `?${query}` : ''}`
  );
}

export async function getLines(): Promise<LinesResponse> {
  return fetchAPI<LinesResponse>('/lines');
}

export async function getStations(lineId?: string): Promise<StationsResponse> {
  const params = lineId ? `?line_id=${lineId}` : '';
  return fetchAPI<StationsResponse>(`/stations${params}`);
}

export async function createIncident(
  data: CreateIncidentRequest
): Promise<CreateIncidentResponse> {
  // Convert camelCase to snake_case for API
  const payload = {
    line: data.line,
    station: data.station,
    timestamp: data.timestamp,
    duration_minutes: data.durationMinutes,
    incident_type: data.incidentType,
  };

  return fetchAPI<CreateIncidentResponse>('/incidents', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export { APIError };
