package backend

import (
	"time"

	"github.com/google/uuid"
)

type Line struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type Station struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	LineID    uuid.UUID `db:"line_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

type Incident struct {
	ID              uuid.UUID `db:"id"`
	StationID       uuid.UUID `db:"station_id"`
	LineID          uuid.UUID `db:"line_id"`
	Timestamp       time.Time `db:"ts"`
	DurationMinutes int32     `db:"duration_minutes"`
	IncidentType    string    `db:"incident_type"`
	Status          string    `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
}

type IncidentWithDetails struct {
	ID              uuid.UUID `db:"id"`
	StationID       uuid.UUID `db:"station_id"`
	LineID          uuid.UUID `db:"line_id"`
	Timestamp       time.Time `db:"ts"`
	DurationMinutes int32     `db:"duration_minutes"`
	IncidentType    string    `db:"incident_type"`
	Status          string    `db:"status"`
	LineName        string    `db:"line_name"`
	StationName     string    `db:"station_name"`
}

type BreakdownCount struct {
	Name  string `db:"name"`
	Count int32  `db:"count"`
}

type MTBFResult struct {
	LineName    string  `db:"line_name"`
	MTBFMinutes float64 `db:"mtbf_minutes"`
}

type StationWithLine struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	LineID    uuid.UUID `db:"line_id"`
	LineName  string    `db:"line_name"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}
