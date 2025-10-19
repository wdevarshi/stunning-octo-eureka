
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS incidents CASCADE;
DROP TABLE IF EXISTS stations CASCADE;
DROP TABLE IF EXISTS lines CASCADE;

CREATE TABLE lines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT UNIQUE NOT NULL CHECK (LENGTH(TRIM(name)) > 0 AND LENGTH(name) <= 100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE stations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL CHECK (LENGTH(TRIM(name)) > 0 AND LENGTH(name) <= 100),
    line_id UUID NOT NULL REFERENCES lines(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'maintenance', 'closed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(name, line_id)
);

CREATE TABLE incidents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    station_id UUID NOT NULL REFERENCES stations(id) ON DELETE CASCADE,
    line_id UUID NOT NULL REFERENCES lines(id) ON DELETE CASCADE,
    ts TIMESTAMPTZ NOT NULL,
    duration_minutes INT NOT NULL,
    incident_type TEXT NOT NULL CHECK (incident_type IN ('mechanical', 'power', 'signal', 'weather', 'other')),
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'investigating', 'resolved', 'closed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(station_id, line_id, ts)
);

CREATE INDEX idx_incidents_ts ON incidents(ts DESC);
CREATE INDEX idx_incidents_station_id ON incidents(station_id);
CREATE INDEX idx_incidents_line_id ON incidents(line_id);
CREATE INDEX idx_incidents_line_ts ON incidents(line_id, ts DESC);
CREATE INDEX idx_incidents_station_ts ON incidents(station_id, ts DESC);
CREATE INDEX idx_incidents_status ON incidents(status);
CREATE INDEX idx_stations_status ON stations(status);

INSERT INTO lines (name) VALUES
    ('North South Line'),
    ('East West Line'),
    ('Circle Line'),
    ('Downtown Line'),
    ('Thomson East Coast Line'),
    ('North East Line');

INSERT INTO stations (name, line_id)
SELECT 'Jurong East', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Bukit Batok', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Bukit Gombak', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Choa Chu Kang', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Yew Tee', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Kranji', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Marsiling', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Woodlands', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Admiralty', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Sembawang', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Canberra', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Yishun', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Khatib', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Yio Chu Kang', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Ang Mo Kio', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Bishan', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Braddell', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Toa Payoh', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Novena', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Newton', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Orchard', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Somerset', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Dhoby Ghaut', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'City Hall', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Raffles Place', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Marina Bay', id FROM lines WHERE name = 'North South Line'
UNION ALL
SELECT 'Pasir Ris', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Tampines', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Simei', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Tanah Merah', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Bedok', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Kembangan', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Eunos', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Paya Lebar', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Aljunied', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Kallang', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Lavender', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Bugis', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'City Hall', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Raffles Place', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Tanjong Pagar', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Outram Park', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Tiong Bahru', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Redhill', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Queenstown', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Commonwealth', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Buona Vista', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Dover', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Clementi', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Jurong East', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Chinese Garden', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Lakeside', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Boon Lay', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Pioneer', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Joo Koon', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Gul Circle', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Tuas Crescent', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Tuas West Road', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Tuas Link', id FROM lines WHERE name = 'East West Line'
UNION ALL
SELECT 'Dhoby Ghaut', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Bras Basah', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Esplanade', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Promenade', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Nicoll Highway', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Stadium', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Mountbatten', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Dakota', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Paya Lebar', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'MacPherson', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Tai Seng', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Bartley', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Serangoon', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Lorong Chuan', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Bishan', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Marymount', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Caldecott', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Botanic Gardens', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Farrer Road', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Holland Village', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Buona Vista', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'one-north', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Kent Ridge', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Haw Par Villa', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Pasir Panjang', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Labrador Park', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Telok Blangah', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'HarbourFront', id FROM lines WHERE name = 'Circle Line'
UNION ALL
SELECT 'Chinatown', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Telok Ayer', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Downtown', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Bayfront', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Promenade', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Bugis', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Rochor', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Little India', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Farrer Park', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Boon Keng', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Bendemeer', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Geylang Bahru', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Mattar', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'MacPherson', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Ubi', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Kaki Bukit', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Bedok North', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Bedok Reservoir', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Tampines West', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Tampines', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Tampines East', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Upper Changi', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Expo', id FROM lines WHERE name = 'Downtown Line'
UNION ALL
SELECT 'Woodlands', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Woodlands South', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Springleaf', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Lentor', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Mayflower', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Bright Hill', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Upper Thomson', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Caldecott', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Stevens', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Napier', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Orchard Boulevard', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Orchard', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Great World', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Havelock', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Outram Park', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Maxwell', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Shenton Way', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'Marina Bay', id FROM lines WHERE name = 'Thomson East Coast Line'
UNION ALL
SELECT 'HarbourFront', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Outram Park', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Chinatown', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Clarke Quay', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Dhoby Ghaut', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Little India', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Farrer Park', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Boon Keng', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Potong Pasir', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Woodleigh', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Serangoon', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Kovan', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Hougang', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Buangkok', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Sengkang', id FROM lines WHERE name = 'North East Line'
UNION ALL
SELECT 'Punggol', id FROM lines WHERE name = 'North East Line';


SELECT
    'Lines' as entity,
    COUNT(*) as count
FROM lines
UNION ALL
SELECT
    'Stations' as entity,
    COUNT(*) as count
FROM stations
UNION ALL
SELECT
    'Incidents' as entity,
    COUNT(*) as count
FROM incidents;
