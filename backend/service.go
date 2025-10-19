package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-coldbrew/log"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/bluesg/transport-analytics/proto"
)

type RepositoryInterface interface {
	CreateLine(ctx context.Context, name string) (*Line, error)
	ListLines(ctx context.Context) ([]Line, error)
	GetLine(ctx context.Context, id uuid.UUID) (*Line, error)
	UpdateLine(ctx context.Context, id uuid.UUID, name string) (*Line, error)
	DeleteLine(ctx context.Context, id uuid.UUID) error
	GetOrCreateLine(ctx context.Context, name string) (*Line, error)

	CreateStation(ctx context.Context, name string, lineID uuid.UUID, status string) (*StationWithLine, error)
	ListStations(ctx context.Context, lineID *uuid.UUID) ([]StationWithLine, error)
	GetStation(ctx context.Context, id uuid.UUID) (*StationWithLine, error)
	UpdateStation(ctx context.Context, id uuid.UUID, name, status *string) (*StationWithLine, error)
	DeleteStation(ctx context.Context, id uuid.UUID) error
	GetOrCreateStation(ctx context.Context, name string, lineID uuid.UUID) (*Station, error)

	CreateIncident(ctx context.Context, stationID, lineID uuid.UUID, ts time.Time, durationMinutes int32, incidentType string) (*Incident, error)
	GetIncidentWithDetails(ctx context.Context, incidentID uuid.UUID) (*IncidentWithDetails, error)
	GetTopBreakdownsByLine(ctx context.Context, limit int32) ([]BreakdownCount, error)
	GetTopBreakdownsByStation(ctx context.Context, limit int32) ([]BreakdownCount, error)
	CalculateMTBF(ctx context.Context) ([]MTBFResult, error)
	GetRecentDisruptions(ctx context.Context, lineName, stationName string, limit int32) ([]IncidentWithDetails, error)
}

type Service struct {
	pb.UnimplementedTransportAnalyticsServer
	repo RepositoryInterface
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) HealthCheck(ctx context.Context, _ *emptypb.Empty) (*httpbody.HttpBody, error) {
	health := map[string]interface{}{
		"status":     "healthy",
		"assessment": "backend-analytics",
		"time":       time.Now().UTC().Format(time.RFC3339),
	}
	data, _ := json.Marshal(health)
	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        data,
	}, nil
}

func (s *Service) ReadyCheck(ctx context.Context, _ *emptypb.Empty) (*httpbody.HttpBody, error) {
	ready := map[string]interface{}{
		"status":     "ready",
		"assessment": "backend-analytics",
		"time":       time.Now().UTC().Format(time.RFC3339),
	}
	data, _ := json.Marshal(ready)
	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        data,
	}, nil
}

func (s *Service) CreateIncident(ctx context.Context, req *pb.CreateIncidentRequest) (*pb.IncidentResponse, error) {
	if err := s.validateIncidentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Info(ctx, "Creating incident", "line", req.Line, "station", req.Station)

	line, err := s.repo.GetOrCreateLine(ctx, strings.TrimSpace(req.Line))
	if err != nil {
		log.Error(ctx, "Failed to get/create line", "error", err)
		return nil, status.Error(codes.Internal, "failed to process line")
	}

	station, err := s.repo.GetOrCreateStation(ctx, strings.TrimSpace(req.Station), line.ID)
	if err != nil {
		log.Error(ctx, "Failed to get/create station", "error", err)
		return nil, status.Error(codes.Internal, "failed to process station")
	}

	ts := req.Timestamp.AsTime()
	incident, err := s.repo.CreateIncident(ctx, station.ID, line.ID, ts, req.DurationMinutes, req.IncidentType)
	if err != nil {
		log.Error(ctx, "Failed to create incident", "error", err)
		return nil, status.Error(codes.Internal, "failed to create incident")
	}

	log.Info(ctx, "Incident created successfully", "incident_id", incident.ID.String())

	return &pb.IncidentResponse{
		Id:              incident.ID.String(),
		Line:            req.Line,
		Station:         req.Station,
		Timestamp:       timestamppb.New(incident.Timestamp),
		DurationMinutes: incident.DurationMinutes,
		IncidentType:    incident.IncidentType,
		LineId:          line.ID.String(),
		StationId:       station.ID.String(),
		Status:          incident.Status,
	}, nil
}

func (s *Service) GetTopBreakdowns(ctx context.Context, req *pb.TopBreakdownsRequest) (*pb.TopBreakdownsResponse, error) {
	scope := strings.ToLower(strings.TrimSpace(req.Scope))
	if scope != "line" && scope != "station" {
		return nil, status.Error(codes.InvalidArgument, "scope must be 'line' or 'station'")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 5
	}
	if limit > 100 {
		limit = 100
	}

	log.Info(ctx, "Getting top breakdowns", "scope", scope, "limit", limit)

	var breakdowns []BreakdownCount
	var err error

	if scope == "line" {
		breakdowns, err = s.repo.GetTopBreakdownsByLine(ctx, limit)
	} else {
		breakdowns, err = s.repo.GetTopBreakdownsByStation(ctx, limit)
	}

	if err != nil {
		log.Error(ctx, "Failed to get top breakdowns", "error", err)
		return nil, status.Error(codes.Internal, "failed to get breakdowns")
	}

	items := make([]*pb.TopBreakdownItem, len(breakdowns))
	for i, b := range breakdowns {
		items[i] = &pb.TopBreakdownItem{
			Name:  b.Name,
			Count: b.Count,
		}
	}

	return &pb.TopBreakdownsResponse{
		Scope: scope,
		Items: items,
	}, nil
}

func (s *Service) GetMTBF(ctx context.Context, _ *emptypb.Empty) (*pb.MTBFResponse, error) {
	log.Info(ctx, "Calculating MTBF for all lines")

	results, err := s.repo.CalculateMTBF(ctx)
	if err != nil {
		log.Error(ctx, "Failed to calculate MTBF", "error", err)
		return nil, status.Error(codes.Internal, "failed to calculate MTBF")
	}

	lines := make([]*pb.MTBFLineItem, len(results))
	for i, r := range results {
		lines[i] = &pb.MTBFLineItem{
			Name:        r.LineName,
			MtbfMinutes: r.MTBFMinutes,
		}
	}

	return &pb.MTBFResponse{
		Lines: lines,
	}, nil
}

func (s *Service) GetRecentDisruptions(ctx context.Context, req *pb.RecentDisruptionsRequest) (*pb.RecentDisruptionsResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	lineName := strings.TrimSpace(req.Line)
	stationName := strings.TrimSpace(req.Station)

	log.Info(ctx, "Getting recent disruptions",
		"line", lineName,
		"station", stationName,
		"limit", limit)

	incidents, err := s.repo.GetRecentDisruptions(ctx, lineName, stationName, limit)
	if err != nil {
		log.Error(ctx, "Failed to get recent disruptions", "error", err)
		return nil, status.Error(codes.Internal, "failed to get disruptions")
	}

	items := make([]*pb.RecentDisruptionItem, len(incidents))
	for i, inc := range incidents {
		items[i] = &pb.RecentDisruptionItem{
			Line:            inc.LineName,
			Station:         inc.StationName,
			Timestamp:       timestamppb.New(inc.Timestamp),
			DurationMinutes: inc.DurationMinutes,
			IncidentType:    inc.IncidentType,
			Status:          inc.Status,
		}
	}

	return &pb.RecentDisruptionsResponse{
		Items: items,
	}, nil
}

func (s *Service) validateIncidentRequest(req *pb.CreateIncidentRequest) error {
	line := strings.TrimSpace(req.Line)
	if line == "" {
		return fmt.Errorf("line must not be empty")
	}
	if len(line) > 100 {
		return fmt.Errorf("line must not exceed 100 characters")
	}

	station := strings.TrimSpace(req.Station)
	if station == "" {
		return fmt.Errorf("station must not be empty")
	}
	if len(station) > 100 {
		return fmt.Errorf("station must not exceed 100 characters")
	}

	if req.Timestamp == nil {
		return fmt.Errorf("timestamp is required")
	}

	ts := req.Timestamp.AsTime()
	if ts.After(time.Now().UTC()) {
		return fmt.Errorf("timestamp cannot be in the future")
	}

	if req.DurationMinutes < 0 || req.DurationMinutes > 1440 {
		return fmt.Errorf("duration_minutes must be between 0 and 1440")
	}

	validTypes := map[string]bool{
		"mechanical": true,
		"power":      true,
		"signal":     true,
		"weather":    true,
		"other":      true,
	}
	if !validTypes[req.IncidentType] {
		return fmt.Errorf("incident_type must be one of: mechanical, power, signal, weather, other")
	}

	return nil
}

func (s *Service) CreateLine(ctx context.Context, req *pb.CreateLineRequest) (*pb.LineResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name must not be empty")
	}
	if len(name) > 100 {
		return nil, status.Error(codes.InvalidArgument, "name must not exceed 100 characters")
	}

	log.Info(ctx, "Creating line", "name", name)

	line, err := s.repo.CreateLine(ctx, name)
	if err != nil {
		log.Error(ctx, "Failed to create line", "error", err)
		return nil, status.Error(codes.Internal, "failed to create line")
	}

	log.Info(ctx, "Line created successfully", "line_id", line.ID.String())

	return &pb.LineResponse{
		Id:        line.ID.String(),
		Name:      line.Name,
		CreatedAt: timestamppb.New(line.CreatedAt),
	}, nil
}

func (s *Service) ListLines(ctx context.Context, _ *emptypb.Empty) (*pb.ListLinesResponse, error) {
	log.Info(ctx, "Listing all lines")

	lines, err := s.repo.ListLines(ctx)
	if err != nil {
		log.Error(ctx, "Failed to list lines", "error", err)
		return nil, status.Error(codes.Internal, "failed to list lines")
	}

	responses := make([]*pb.LineResponse, len(lines))
	for i, line := range lines {
		responses[i] = &pb.LineResponse{
			Id:        line.ID.String(),
			Name:      line.Name,
			CreatedAt: timestamppb.New(line.CreatedAt),
		}
	}

	return &pb.ListLinesResponse{Lines: responses}, nil
}

func (s *Service) GetLine(ctx context.Context, req *pb.GetLineRequest) (*pb.LineResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid line ID")
	}

	log.Info(ctx, "Getting line", "id", id.String())

	line, err := s.repo.GetLine(ctx, id)
	if err == ErrNotFound {
		return nil, status.Error(codes.NotFound, "line not found")
	}
	if err != nil {
		log.Error(ctx, "Failed to get line", "error", err)
		return nil, status.Error(codes.Internal, "failed to get line")
	}

	return &pb.LineResponse{
		Id:        line.ID.String(),
		Name:      line.Name,
		CreatedAt: timestamppb.New(line.CreatedAt),
	}, nil
}

func (s *Service) UpdateLine(ctx context.Context, req *pb.UpdateLineRequest) (*pb.LineResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid line ID")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name must not be empty")
	}
	if len(name) > 100 {
		return nil, status.Error(codes.InvalidArgument, "name must not exceed 100 characters")
	}

	log.Info(ctx, "Updating line", "id", id.String(), "name", name)

	line, err := s.repo.UpdateLine(ctx, id, name)
	if err == ErrNotFound {
		return nil, status.Error(codes.NotFound, "line not found")
	}
	if err != nil {
		log.Error(ctx, "Failed to update line", "error", err)
		return nil, status.Error(codes.Internal, "failed to update line")
	}

	log.Info(ctx, "Line updated successfully", "line_id", line.ID.String())

	return &pb.LineResponse{
		Id:        line.ID.String(),
		Name:      line.Name,
		CreatedAt: timestamppb.New(line.CreatedAt),
	}, nil
}

func (s *Service) DeleteLine(ctx context.Context, req *pb.DeleteLineRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid line ID")
	}

	log.Info(ctx, "Deleting line", "id", id.String())

	err = s.repo.DeleteLine(ctx, id)
	if err == ErrNotFound {
		return nil, status.Error(codes.NotFound, "line not found")
	}
	if err != nil {
		log.Error(ctx, "Failed to delete line", "error", err)
		return nil, status.Error(codes.Internal, "failed to delete line")
	}

	log.Info(ctx, "Line deleted successfully", "line_id", id.String())

	return &emptypb.Empty{}, nil
}

func (s *Service) CreateStation(ctx context.Context, req *pb.CreateStationRequest) (*pb.StationResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name must not be empty")
	}
	if len(name) > 100 {
		return nil, status.Error(codes.InvalidArgument, "name must not exceed 100 characters")
	}

	lineID, err := uuid.Parse(req.LineId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid line ID")
	}

	statusVal := strings.TrimSpace(req.Status)
	if statusVal == "" {
		statusVal = "active"
	}
	validStatuses := map[string]bool{
		"active":      true,
		"inactive":    true,
		"maintenance": true,
		"closed":      true,
	}
	if !validStatuses[statusVal] {
		return nil, status.Error(codes.InvalidArgument, "status must be one of: active, inactive, maintenance, closed")
	}

	log.Info(ctx, "Creating station", "name", name, "line_id", lineID.String(), "status", statusVal)

	station, err := s.repo.CreateStation(ctx, name, lineID, statusVal)
	if err == ErrNotFound {
		return nil, status.Error(codes.NotFound, "line not found")
	}
	if err != nil {
		log.Error(ctx, "Failed to create station", "error", err)
		return nil, status.Error(codes.Internal, "failed to create station")
	}

	log.Info(ctx, "Station created successfully", "station_id", station.ID.String())

	return &pb.StationResponse{
		Id:        station.ID.String(),
		Name:      station.Name,
		LineId:    station.LineID.String(),
		LineName:  station.LineName,
		Status:    station.Status,
		CreatedAt: timestamppb.New(station.CreatedAt),
	}, nil
}

func (s *Service) ListStations(ctx context.Context, req *pb.ListStationsRequest) (*pb.ListStationsResponse, error) {
	var lineID *uuid.UUID
	if req.LineId != "" {
		id, err := uuid.Parse(req.LineId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid line ID")
		}
		lineID = &id
	}

	log.Info(ctx, "Listing stations", "line_id", req.LineId)

	stations, err := s.repo.ListStations(ctx, lineID)
	if err != nil {
		log.Error(ctx, "Failed to list stations", "error", err)
		return nil, status.Error(codes.Internal, "failed to list stations")
	}

	responses := make([]*pb.StationResponse, len(stations))
	for i, station := range stations {
		responses[i] = &pb.StationResponse{
			Id:        station.ID.String(),
			Name:      station.Name,
			LineId:    station.LineID.String(),
			LineName:  station.LineName,
			Status:    station.Status,
			CreatedAt: timestamppb.New(station.CreatedAt),
		}
	}

	return &pb.ListStationsResponse{Stations: responses}, nil
}

func (s *Service) GetStation(ctx context.Context, req *pb.GetStationRequest) (*pb.StationResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid station ID")
	}

	log.Info(ctx, "Getting station", "id", id.String())

	station, err := s.repo.GetStation(ctx, id)
	if err == ErrNotFound {
		return nil, status.Error(codes.NotFound, "station not found")
	}
	if err != nil {
		log.Error(ctx, "Failed to get station", "error", err)
		return nil, status.Error(codes.Internal, "failed to get station")
	}

	return &pb.StationResponse{
		Id:        station.ID.String(),
		Name:      station.Name,
		LineId:    station.LineID.String(),
		LineName:  station.LineName,
		Status:    station.Status,
		CreatedAt: timestamppb.New(station.CreatedAt),
	}, nil
}

func (s *Service) UpdateStation(ctx context.Context, req *pb.UpdateStationRequest) (*pb.StationResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid station ID")
	}

	var name *string
	if req.Name != "" {
		trimmed := strings.TrimSpace(req.Name)
		if len(trimmed) > 100 {
			return nil, status.Error(codes.InvalidArgument, "name must not exceed 100 characters")
		}
		name = &trimmed
	}

	var statusVal *string
	if req.Status != "" {
		trimmed := strings.TrimSpace(req.Status)
		validStatuses := map[string]bool{
			"active":      true,
			"inactive":    true,
			"maintenance": true,
			"closed":      true,
		}
		if !validStatuses[trimmed] {
			return nil, status.Error(codes.InvalidArgument, "status must be one of: active, inactive, maintenance, closed")
		}
		statusVal = &trimmed
	}

	if name == nil && statusVal == nil {
		return nil, status.Error(codes.InvalidArgument, "at least one field (name or status) must be provided")
	}

	log.Info(ctx, "Updating station", "id", id.String())

	station, err := s.repo.UpdateStation(ctx, id, name, statusVal)
	if err == ErrNotFound {
		return nil, status.Error(codes.NotFound, "station not found")
	}
	if err != nil {
		log.Error(ctx, "Failed to update station", "error", err)
		return nil, status.Error(codes.Internal, "failed to update station")
	}

	log.Info(ctx, "Station updated successfully", "station_id", station.ID.String())

	return &pb.StationResponse{
		Id:        station.ID.String(),
		Name:      station.Name,
		LineId:    station.LineID.String(),
		LineName:  station.LineName,
		Status:    station.Status,
		CreatedAt: timestamppb.New(station.CreatedAt),
	}, nil
}

func (s *Service) DeleteStation(ctx context.Context, req *pb.DeleteStationRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid station ID")
	}

	log.Info(ctx, "Deleting station", "id", id.String())

	err = s.repo.DeleteStation(ctx, id)
	if err == ErrNotFound {
		return nil, status.Error(codes.NotFound, "station not found")
	}
	if err != nil {
		log.Error(ctx, "Failed to delete station", "error", err)
		return nil, status.Error(codes.Internal, "failed to delete station")
	}

	log.Info(ctx, "Station deleted successfully", "station_id", id.String())

	return &emptypb.Empty{}, nil
}
