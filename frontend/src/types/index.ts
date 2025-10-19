export interface TopBreakdownItem {
  name: string;
  count: number;
}

export interface TopBreakdownsResponse {
  scope: string;
  items: TopBreakdownItem[];
}

export interface MTBFLineItem {
  name: string;
  mtbfMinutes: number;
}

export interface MTBFResponse {
  lines: MTBFLineItem[];
}

export interface RecentDisruptionItem {
  line: string;
  station: string;
  timestamp: string;
  durationMinutes: number;
  incidentType: string;
  status: string;
}

export interface RecentDisruptionsResponse {
  items: RecentDisruptionItem[];
}

export interface DashboardData {
  topBreakdownsByLine: TopBreakdownsResponse | null;
  topBreakdownsByStation: TopBreakdownsResponse | null;
  mtbf: MTBFResponse | null;
  recentDisruptions: RecentDisruptionsResponse | null;
}

export interface Line {
  id: string;
  name: string;
  createdAt: string;
}

export interface LinesResponse {
  lines: Line[];
}

export interface Station {
  id: string;
  name: string;
  lineId: string;
  lineName: string;
  status: string;
  createdAt: string;
}

export interface StationsResponse {
  stations: Station[];
}

export interface CreateIncidentRequest {
  line: string;
  station: string;
  timestamp: string;
  durationMinutes: number;
  incidentType: string;
}

export interface CreateIncidentResponse {
  id: string;
  line: string;
  station: string;
  timestamp: string;
  durationMinutes: number;
  incidentType: string;
  status: string;
}
