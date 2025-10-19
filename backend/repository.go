package backend

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrDatabaseError = errors.New("database error")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetOrCreateLine(ctx context.Context, name string) (*Line, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer func() { _ = tx.Rollback() }()

	var line Line
	err = tx.GetContext(ctx, &line, "SELECT id, name, created_at FROM lines WHERE name = $1", name)
	if err == nil {
		return &line, nil
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	err = tx.GetContext(ctx, &line,
		"INSERT INTO lines (name) VALUES ($1) RETURNING id, name, created_at",
		name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &line, nil
}

func (r *Repository) GetOrCreateStation(ctx context.Context, name string, lineID uuid.UUID) (*Station, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer func() { _ = tx.Rollback() }()

	var station Station
	err = tx.GetContext(ctx, &station,
		"SELECT id, name, line_id, status, created_at FROM stations WHERE name = $1 AND line_id = $2",
		name, lineID)
	if err == nil {
		return &station, nil
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	err = tx.GetContext(ctx, &station,
		"INSERT INTO stations (name, line_id) VALUES ($1, $2) RETURNING id, name, line_id, status, created_at",
		name, lineID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &station, nil
}

func (r *Repository) CreateIncident(ctx context.Context, stationID, lineID uuid.UUID, ts time.Time, durationMinutes int32, incidentType string) (*Incident, error) {
	var incident Incident
	err := r.db.GetContext(ctx, &incident,
		`INSERT INTO incidents (station_id, line_id, ts, duration_minutes, incident_type)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (station_id, line_id, ts) DO UPDATE
		 SET duration_minutes = EXCLUDED.duration_minutes, incident_type = EXCLUDED.incident_type
		 RETURNING id, station_id, line_id, ts, duration_minutes, incident_type, status, created_at`,
		stationID, lineID, ts, durationMinutes, incidentType)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &incident, nil
}

func (r *Repository) GetIncidentWithDetails(ctx context.Context, incidentID uuid.UUID) (*IncidentWithDetails, error) {
	var incident IncidentWithDetails
	err := r.db.GetContext(ctx, &incident,
		`SELECT i.id, i.station_id, i.line_id, i.ts, i.duration_minutes, i.incident_type, i.status,
		        l.name as line_name, s.name as station_name
		 FROM incidents i
		 JOIN lines l ON i.line_id = l.id
		 JOIN stations s ON i.station_id = s.id
		 WHERE i.id = $1`,
		incidentID)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &incident, nil
}

func (r *Repository) GetTopBreakdownsByLine(ctx context.Context, limit int32) ([]BreakdownCount, error) {
	var results []BreakdownCount
	err := r.db.SelectContext(ctx, &results,
		`SELECT l.name, COUNT(i.id)::int as count
		 FROM lines l
		 LEFT JOIN incidents i ON l.id = i.line_id
		 GROUP BY l.name
		 ORDER BY count DESC
		 LIMIT $1`,
		limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return results, nil
}

func (r *Repository) GetTopBreakdownsByStation(ctx context.Context, limit int32) ([]BreakdownCount, error) {
	var results []BreakdownCount
	err := r.db.SelectContext(ctx, &results,
		`SELECT s.name, COUNT(i.id)::int as count
		 FROM stations s
		 LEFT JOIN incidents i ON s.id = i.station_id
		 GROUP BY s.name
		 ORDER BY count DESC
		 LIMIT $1`,
		limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return results, nil
}

func (r *Repository) CalculateMTBF(ctx context.Context) ([]MTBFResult, error) {
	var results []MTBFResult
	query := `
		WITH line_incidents AS (
			SELECT
				l.name as line_name,
				i.ts,
				LAG(i.ts) OVER (PARTITION BY l.id ORDER BY i.ts) as prev_ts
			FROM incidents i
			JOIN lines l ON i.line_id = l.id
			ORDER BY l.id, i.ts
		),
		time_deltas AS (
			SELECT
				line_name,
				EXTRACT(EPOCH FROM (ts - prev_ts)) / 60.0 as minutes_between
			FROM line_incidents
			WHERE prev_ts IS NOT NULL
		),
		line_stats AS (
			SELECT
				line_name,
				COUNT(*) as incident_count,
				AVG(minutes_between) as avg_minutes_between
			FROM time_deltas
			GROUP BY line_name
		)
		SELECT
			line_name,
			ROUND(avg_minutes_between::numeric, 2)::float8 as mtbf_minutes
		FROM line_stats
		WHERE incident_count >= 1
		ORDER BY line_name`

	err := r.db.SelectContext(ctx, &results, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return results, nil
}

func (r *Repository) GetRecentDisruptions(ctx context.Context, lineName, stationName string, limit int32) ([]IncidentWithDetails, error) {
	var results []IncidentWithDetails

	query := `
		SELECT i.id, i.station_id, i.line_id, i.ts, i.duration_minutes, i.incident_type, i.status,
		       l.name as line_name, s.name as station_name
		FROM incidents i
		JOIN lines l ON i.line_id = l.id
		JOIN stations s ON i.station_id = s.id
		WHERE 1=1`

	args := []interface{}{}
	argPos := 1

	if lineName != "" {
		query += fmt.Sprintf(" AND l.name = $%d", argPos)
		args = append(args, lineName)
		argPos++
	}

	if stationName != "" {
		query += fmt.Sprintf(" AND s.name = $%d", argPos)
		args = append(args, stationName)
		argPos++
	}

	query += " ORDER BY i.ts DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, limit)
	}

	err := r.db.SelectContext(ctx, &results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return results, nil
}

func (r *Repository) CreateLine(ctx context.Context, name string) (*Line, error) {
	var line Line
	err := r.db.GetContext(ctx, &line,
		`INSERT INTO lines (name) VALUES ($1)
		 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		 RETURNING id, name, created_at`,
		name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &line, nil
}

func (r *Repository) ListLines(ctx context.Context) ([]Line, error) {
	var lines []Line
	err := r.db.SelectContext(ctx, &lines,
		"SELECT id, name, created_at FROM lines ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return lines, nil
}

func (r *Repository) GetLine(ctx context.Context, id uuid.UUID) (*Line, error) {
	var line Line
	err := r.db.GetContext(ctx, &line,
		"SELECT id, name, created_at FROM lines WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &line, nil
}

func (r *Repository) UpdateLine(ctx context.Context, id uuid.UUID, name string) (*Line, error) {
	var line Line
	err := r.db.GetContext(ctx, &line,
		"UPDATE lines SET name = $1 WHERE id = $2 RETURNING id, name, created_at",
		name, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &line, nil
}

func (r *Repository) DeleteLine(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM lines WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) CreateStation(ctx context.Context, name string, lineID uuid.UUID, status string) (*StationWithLine, error) {
	var lineExists bool
	err := r.db.GetContext(ctx, &lineExists, "SELECT EXISTS(SELECT 1 FROM lines WHERE id = $1)", lineID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if !lineExists {
		return nil, ErrNotFound
	}

	if status == "" {
		status = "active"
	}

	var station StationWithLine
	err = r.db.GetContext(ctx, &station,
		`INSERT INTO stations (name, line_id, status)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (name, line_id) DO UPDATE SET status = EXCLUDED.status
		 RETURNING id, name, line_id, status, created_at,
		 (SELECT name FROM lines WHERE id = $2) as line_name`,
		name, lineID, status)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &station, nil
}

func (r *Repository) ListStations(ctx context.Context, lineID *uuid.UUID) ([]StationWithLine, error) {
	var stations []StationWithLine
	query := `
		SELECT s.id, s.name, s.line_id, l.name as line_name, s.status, s.created_at
		FROM stations s
		JOIN lines l ON s.line_id = l.id`

	args := []interface{}{}
	if lineID != nil {
		query += " WHERE s.line_id = $1"
		args = append(args, *lineID)
	}
	query += " ORDER BY l.name, s.name"

	err := r.db.SelectContext(ctx, &stations, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return stations, nil
}

func (r *Repository) GetStation(ctx context.Context, id uuid.UUID) (*StationWithLine, error) {
	var station StationWithLine
	err := r.db.GetContext(ctx, &station,
		`SELECT s.id, s.name, s.line_id, l.name as line_name, s.status, s.created_at
		 FROM stations s
		 JOIN lines l ON s.line_id = l.id
		 WHERE s.id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &station, nil
}

func (r *Repository) UpdateStation(ctx context.Context, id uuid.UUID, name, status *string) (*StationWithLine, error) {
	current, err := r.GetStation(ctx, id)
	if err != nil {
		return nil, err
	}

	newName := current.Name
	if name != nil && *name != "" {
		newName = *name
	}
	newStatus := current.Status
	if status != nil && *status != "" {
		newStatus = *status
	}

	var station StationWithLine
	err = r.db.GetContext(ctx, &station,
		`UPDATE stations
		 SET name = $1, status = $2
		 WHERE id = $3
		 RETURNING id, name, line_id, status, created_at,
		 (SELECT name FROM lines WHERE id = line_id) as line_name`,
		newName, newStatus, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &station, nil
}

func (r *Repository) DeleteStation(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM stations WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
