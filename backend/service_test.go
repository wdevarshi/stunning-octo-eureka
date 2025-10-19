package backend

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/bluesg/transport-analytics/proto"
)

type MockRepository struct {
	CreateLineFn      func(ctx context.Context, name string) (*Line, error)
	ListLinesFn       func(ctx context.Context) ([]Line, error)
	GetLineFn         func(ctx context.Context, id uuid.UUID) (*Line, error)
	UpdateLineFn      func(ctx context.Context, id uuid.UUID, name string) (*Line, error)
	DeleteLineFn      func(ctx context.Context, id uuid.UUID) error
	GetOrCreateLineFn func(ctx context.Context, name string) (*Line, error)

	CreateStationFn      func(ctx context.Context, name string, lineID uuid.UUID, status string) (*StationWithLine, error)
	ListStationsFn       func(ctx context.Context, lineID *uuid.UUID) ([]StationWithLine, error)
	GetStationFn         func(ctx context.Context, id uuid.UUID) (*StationWithLine, error)
	UpdateStationFn      func(ctx context.Context, id uuid.UUID, name, status *string) (*StationWithLine, error)
	DeleteStationFn      func(ctx context.Context, id uuid.UUID) error
	GetOrCreateStationFn func(ctx context.Context, name string, lineID uuid.UUID) (*Station, error)

	CreateIncidentFn             func(ctx context.Context, stationID, lineID uuid.UUID, ts time.Time, durationMinutes int32, incidentType string) (*Incident, error)
	GetIncidentWithDetailsFn     func(ctx context.Context, incidentID uuid.UUID) (*IncidentWithDetails, error)
	GetTopBreakdownsByLineFn     func(ctx context.Context, limit int32) ([]BreakdownCount, error)
	GetTopBreakdownsByStationFn  func(ctx context.Context, limit int32) ([]BreakdownCount, error)
	CalculateMTBFFn              func(ctx context.Context) ([]MTBFResult, error)
	GetRecentDisruptionsFn       func(ctx context.Context, lineName, stationName string, limit int32) ([]IncidentWithDetails, error)
}

func (m *MockRepository) CreateLine(ctx context.Context, name string) (*Line, error) {
	if m.CreateLineFn != nil {
		return m.CreateLineFn(ctx, name)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) ListLines(ctx context.Context) ([]Line, error) {
	if m.ListLinesFn != nil {
		return m.ListLinesFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) GetLine(ctx context.Context, id uuid.UUID) (*Line, error) {
	if m.GetLineFn != nil {
		return m.GetLineFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) UpdateLine(ctx context.Context, id uuid.UUID, name string) (*Line, error) {
	if m.UpdateLineFn != nil {
		return m.UpdateLineFn(ctx, id, name)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) DeleteLine(ctx context.Context, id uuid.UUID) error {
	if m.DeleteLineFn != nil {
		return m.DeleteLineFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *MockRepository) GetOrCreateLine(ctx context.Context, name string) (*Line, error) {
	if m.GetOrCreateLineFn != nil {
		return m.GetOrCreateLineFn(ctx, name)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) CreateStation(ctx context.Context, name string, lineID uuid.UUID, status string) (*StationWithLine, error) {
	if m.CreateStationFn != nil {
		return m.CreateStationFn(ctx, name, lineID, status)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) ListStations(ctx context.Context, lineID *uuid.UUID) ([]StationWithLine, error) {
	if m.ListStationsFn != nil {
		return m.ListStationsFn(ctx, lineID)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) GetStation(ctx context.Context, id uuid.UUID) (*StationWithLine, error) {
	if m.GetStationFn != nil {
		return m.GetStationFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) UpdateStation(ctx context.Context, id uuid.UUID, name, status *string) (*StationWithLine, error) {
	if m.UpdateStationFn != nil {
		return m.UpdateStationFn(ctx, id, name, status)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) DeleteStation(ctx context.Context, id uuid.UUID) error {
	if m.DeleteStationFn != nil {
		return m.DeleteStationFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *MockRepository) GetOrCreateStation(ctx context.Context, name string, lineID uuid.UUID) (*Station, error) {
	if m.GetOrCreateStationFn != nil {
		return m.GetOrCreateStationFn(ctx, name, lineID)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) CreateIncident(ctx context.Context, stationID, lineID uuid.UUID, ts time.Time, durationMinutes int32, incidentType string) (*Incident, error) {
	if m.CreateIncidentFn != nil {
		return m.CreateIncidentFn(ctx, stationID, lineID, ts, durationMinutes, incidentType)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) GetIncidentWithDetails(ctx context.Context, incidentID uuid.UUID) (*IncidentWithDetails, error) {
	if m.GetIncidentWithDetailsFn != nil {
		return m.GetIncidentWithDetailsFn(ctx, incidentID)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) GetTopBreakdownsByLine(ctx context.Context, limit int32) ([]BreakdownCount, error) {
	if m.GetTopBreakdownsByLineFn != nil {
		return m.GetTopBreakdownsByLineFn(ctx, limit)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) GetTopBreakdownsByStation(ctx context.Context, limit int32) ([]BreakdownCount, error) {
	if m.GetTopBreakdownsByStationFn != nil {
		return m.GetTopBreakdownsByStationFn(ctx, limit)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) CalculateMTBF(ctx context.Context) ([]MTBFResult, error) {
	if m.CalculateMTBFFn != nil {
		return m.CalculateMTBFFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) GetRecentDisruptions(ctx context.Context, lineName, stationName string, limit int32) ([]IncidentWithDetails, error) {
	if m.GetRecentDisruptionsFn != nil {
		return m.GetRecentDisruptionsFn(ctx, lineName, stationName, limit)
	}
	return nil, errors.New("not implemented")
}

func setupServiceWithMock() (*Service, *MockRepository) {
	mockRepo := &MockRepository{}
	service := &Service{repo: mockRepo}
	return service, mockRepo
}


func TestCreateLine_Success(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()
	now := time.Now()

	mockRepo.CreateLineFn = func(ctx context.Context, name string) (*Line, error) {
		assert.Equal(t, "Test Line", name)
		return &Line{
			ID:        lineID,
			Name:      name,
			CreatedAt: now,
		}, nil
	}

	req := &pb.CreateLineRequest{Name: "Test Line"}
	resp, err := service.CreateLine(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, lineID.String(), resp.Id)
	assert.Equal(t, "Test Line", resp.Name)
	assert.NotNil(t, resp.CreatedAt)
}

func TestCreateLine_EmptyName(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.CreateLineRequest{Name: ""}
	resp, err := service.CreateLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "name must not be empty")
}

func TestCreateLine_NameWithOnlySpaces(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.CreateLineRequest{Name: "   "}
	resp, err := service.CreateLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestCreateLine_NameTooLong(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	longName := string(make([]byte, 101))
	req := &pb.CreateLineRequest{Name: longName}
	resp, err := service.CreateLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "must not exceed 100 characters")
}

func TestCreateLine_RepositoryError(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	mockRepo.CreateLineFn = func(ctx context.Context, name string) (*Line, error) {
		return nil, errors.New("database error")
	}

	req := &pb.CreateLineRequest{Name: "Test Line"}
	resp, err := service.CreateLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestListLines_Success(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	now := time.Now()
	mockLines := []Line{
		{ID: uuid.New(), Name: "Line 1", CreatedAt: now},
		{ID: uuid.New(), Name: "Line 2", CreatedAt: now},
		{ID: uuid.New(), Name: "Line 3", CreatedAt: now},
	}

	mockRepo.ListLinesFn = func(ctx context.Context) ([]Line, error) {
		return mockLines, nil
	}

	resp, err := service.ListLines(ctx, nil)

	require.NoError(t, err)
	assert.Len(t, resp.Lines, 3)
	assert.Equal(t, "Line 1", resp.Lines[0].Name)
	assert.Equal(t, "Line 2", resp.Lines[1].Name)
	assert.Equal(t, "Line 3", resp.Lines[2].Name)
}

func TestListLines_EmptyResult(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	mockRepo.ListLinesFn = func(ctx context.Context) ([]Line, error) {
		return []Line{}, nil
	}

	resp, err := service.ListLines(ctx, nil)

	require.NoError(t, err)
	assert.Len(t, resp.Lines, 0)
}

func TestListLines_RepositoryError(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	mockRepo.ListLinesFn = func(ctx context.Context) ([]Line, error) {
		return nil, errors.New("database error")
	}

	resp, err := service.ListLines(ctx, nil)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGetLine_Success(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()
	now := time.Now()

	mockRepo.GetLineFn = func(ctx context.Context, id uuid.UUID) (*Line, error) {
		assert.Equal(t, lineID, id)
		return &Line{
			ID:        lineID,
			Name:      "Test Line",
			CreatedAt: now,
		}, nil
	}

	req := &pb.GetLineRequest{Id: lineID.String()}
	resp, err := service.GetLine(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, lineID.String(), resp.Id)
	assert.Equal(t, "Test Line", resp.Name)
}

func TestGetLine_InvalidUUID(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.GetLineRequest{Id: "invalid-uuid"}
	resp, err := service.GetLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "invalid line ID")
}

func TestGetLine_NotFound(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()

	mockRepo.GetLineFn = func(ctx context.Context, id uuid.UUID) (*Line, error) {
		return nil, ErrNotFound
	}

	req := &pb.GetLineRequest{Id: lineID.String()}
	resp, err := service.GetLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "line not found")
}

func TestUpdateLine_Success(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()
	now := time.Now()

	mockRepo.UpdateLineFn = func(ctx context.Context, id uuid.UUID, name string) (*Line, error) {
		assert.Equal(t, lineID, id)
		assert.Equal(t, "Updated Line", name)
		return &Line{
			ID:        lineID,
			Name:      name,
			CreatedAt: now,
		}, nil
	}

	req := &pb.UpdateLineRequest{
		Id:   lineID.String(),
		Name: "Updated Line",
	}
	resp, err := service.UpdateLine(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, lineID.String(), resp.Id)
	assert.Equal(t, "Updated Line", resp.Name)
}

func TestUpdateLine_InvalidUUID(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.UpdateLineRequest{
		Id:   "invalid-uuid",
		Name: "Updated Line",
	}
	resp, err := service.UpdateLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUpdateLine_EmptyName(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.UpdateLineRequest{
		Id:   uuid.New().String(),
		Name: "",
	}
	resp, err := service.UpdateLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUpdateLine_NameTooLong(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	longName := string(make([]byte, 101))
	req := &pb.UpdateLineRequest{
		Id:   uuid.New().String(),
		Name: longName,
	}
	resp, err := service.UpdateLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUpdateLine_NotFound(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()

	mockRepo.UpdateLineFn = func(ctx context.Context, id uuid.UUID, name string) (*Line, error) {
		return nil, ErrNotFound
	}

	req := &pb.UpdateLineRequest{
		Id:   lineID.String(),
		Name: "Updated Line",
	}
	resp, err := service.UpdateLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestDeleteLine_Success(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()

	mockRepo.DeleteLineFn = func(ctx context.Context, id uuid.UUID) error {
		assert.Equal(t, lineID, id)
		return nil
	}

	req := &pb.DeleteLineRequest{Id: lineID.String()}
	resp, err := service.DeleteLine(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDeleteLine_InvalidUUID(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.DeleteLineRequest{Id: "invalid-uuid"}
	resp, err := service.DeleteLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestDeleteLine_NotFound(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()

	mockRepo.DeleteLineFn = func(ctx context.Context, id uuid.UUID) error {
		return ErrNotFound
	}

	req := &pb.DeleteLineRequest{Id: lineID.String()}
	resp, err := service.DeleteLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestDeleteLine_RepositoryError(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()

	mockRepo.DeleteLineFn = func(ctx context.Context, id uuid.UUID) error {
		return errors.New("database error")
	}

	req := &pb.DeleteLineRequest{Id: lineID.String()}
	resp, err := service.DeleteLine(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}


func TestCreateStation_Success(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()
	lineID := uuid.New()
	now := time.Now()

	mockRepo.CreateStationFn = func(ctx context.Context, name string, lID uuid.UUID, status string) (*StationWithLine, error) {
		assert.Equal(t, "Test Station", name)
		assert.Equal(t, lineID, lID)
		assert.Equal(t, "active", status)
		return &StationWithLine{
			ID:        stationID,
			Name:      name,
			LineID:    lID,
			LineName:  "Test Line",
			Status:    status,
			CreatedAt: now,
		}, nil
	}

	req := &pb.CreateStationRequest{
		Name:   "Test Station",
		LineId: lineID.String(),
		Status: "active",
	}
	resp, err := service.CreateStation(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, stationID.String(), resp.Id)
	assert.Equal(t, "Test Station", resp.Name)
	assert.Equal(t, lineID.String(), resp.LineId)
	assert.Equal(t, "Test Line", resp.LineName)
	assert.Equal(t, "active", resp.Status)
}

func TestCreateStation_DefaultStatus(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()
	lineID := uuid.New()
	now := time.Now()

	mockRepo.CreateStationFn = func(ctx context.Context, name string, lID uuid.UUID, status string) (*StationWithLine, error) {
		assert.Equal(t, "active", status) // Should default to active
		return &StationWithLine{
			ID:        stationID,
			Name:      name,
			LineID:    lID,
			LineName:  "Test Line",
			Status:    "active",
			CreatedAt: now,
		}, nil
	}

	req := &pb.CreateStationRequest{
		Name:   "Test Station",
		LineId: lineID.String(),
	}
	resp, err := service.CreateStation(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "active", resp.Status)
}

func TestCreateStation_EmptyName(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.CreateStationRequest{
		Name:   "",
		LineId: uuid.New().String(),
		Status: "active",
	}
	resp, err := service.CreateStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "name must not be empty")
}

func TestCreateStation_NameTooLong(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	longName := string(make([]byte, 101))
	req := &pb.CreateStationRequest{
		Name:   longName,
		LineId: uuid.New().String(),
		Status: "active",
	}
	resp, err := service.CreateStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "must not exceed 100 characters")
}

func TestCreateStation_InvalidLineID(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.CreateStationRequest{
		Name:   "Test Station",
		LineId: "invalid-uuid",
		Status: "active",
	}
	resp, err := service.CreateStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "invalid line ID")
}

func TestCreateStation_InvalidStatus(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.CreateStationRequest{
		Name:   "Test Station",
		LineId: uuid.New().String(),
		Status: "invalid-status",
	}
	resp, err := service.CreateStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "status must be one of")
}

func TestCreateStation_LineNotFound(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()

	mockRepo.CreateStationFn = func(ctx context.Context, name string, lID uuid.UUID, status string) (*StationWithLine, error) {
		return nil, ErrNotFound
	}

	req := &pb.CreateStationRequest{
		Name:   "Test Station",
		LineId: lineID.String(),
		Status: "active",
	}
	resp, err := service.CreateStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "line not found")
}

func TestCreateStation_RepositoryError(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()

	mockRepo.CreateStationFn = func(ctx context.Context, name string, lID uuid.UUID, status string) (*StationWithLine, error) {
		return nil, errors.New("database error")
	}

	req := &pb.CreateStationRequest{
		Name:   "Test Station",
		LineId: lineID.String(),
		Status: "active",
	}
	resp, err := service.CreateStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestListStations_Success(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()
	now := time.Now()
	mockStations := []StationWithLine{
		{ID: uuid.New(), Name: "Station 1", LineID: lineID, LineName: "Test Line", Status: "active", CreatedAt: now},
		{ID: uuid.New(), Name: "Station 2", LineID: lineID, LineName: "Test Line", Status: "active", CreatedAt: now},
		{ID: uuid.New(), Name: "Station 3", LineID: lineID, LineName: "Test Line", Status: "maintenance", CreatedAt: now},
	}

	mockRepo.ListStationsFn = func(ctx context.Context, lID *uuid.UUID) ([]StationWithLine, error) {
		return mockStations, nil
	}

	resp, err := service.ListStations(ctx, &pb.ListStationsRequest{})

	require.NoError(t, err)
	assert.Len(t, resp.Stations, 3)
	assert.Equal(t, "Station 1", resp.Stations[0].Name)
	assert.Equal(t, "active", resp.Stations[0].Status)
	assert.Equal(t, "maintenance", resp.Stations[2].Status)
}

func TestListStations_FilterByLineID(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()
	now := time.Now()
	mockStations := []StationWithLine{
		{ID: uuid.New(), Name: "Station 1", LineID: lineID, LineName: "Test Line", Status: "active", CreatedAt: now},
	}

	mockRepo.ListStationsFn = func(ctx context.Context, lID *uuid.UUID) ([]StationWithLine, error) {
		assert.NotNil(t, lID)
		assert.Equal(t, lineID, *lID)
		return mockStations, nil
	}

	req := &pb.ListStationsRequest{LineId: lineID.String()}
	resp, err := service.ListStations(ctx, req)

	require.NoError(t, err)
	assert.Len(t, resp.Stations, 1)
}

func TestListStations_InvalidLineID(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.ListStationsRequest{LineId: "invalid-uuid"}
	resp, err := service.ListStations(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestListStations_EmptyResult(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	mockRepo.ListStationsFn = func(ctx context.Context, lID *uuid.UUID) ([]StationWithLine, error) {
		return []StationWithLine{}, nil
	}

	resp, err := service.ListStations(ctx, &pb.ListStationsRequest{})

	require.NoError(t, err)
	assert.Len(t, resp.Stations, 0)
}

func TestListStations_RepositoryError(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	mockRepo.ListStationsFn = func(ctx context.Context, lID *uuid.UUID) ([]StationWithLine, error) {
		return nil, errors.New("database error")
	}

	resp, err := service.ListStations(ctx, &pb.ListStationsRequest{})

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGetStation_Success(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()
	lineID := uuid.New()
	now := time.Now()

	mockRepo.GetStationFn = func(ctx context.Context, id uuid.UUID) (*StationWithLine, error) {
		assert.Equal(t, stationID, id)
		return &StationWithLine{
			ID:        stationID,
			Name:      "Test Station",
			LineID:    lineID,
			LineName:  "Test Line",
			Status:    "active",
			CreatedAt: now,
		}, nil
	}

	req := &pb.GetStationRequest{Id: stationID.String()}
	resp, err := service.GetStation(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, stationID.String(), resp.Id)
	assert.Equal(t, "Test Station", resp.Name)
	assert.Equal(t, lineID.String(), resp.LineId)
	assert.Equal(t, "Test Line", resp.LineName)
	assert.Equal(t, "active", resp.Status)
}

func TestGetStation_InvalidUUID(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.GetStationRequest{Id: "invalid-uuid"}
	resp, err := service.GetStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "invalid station ID")
}

func TestGetStation_NotFound(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()

	mockRepo.GetStationFn = func(ctx context.Context, id uuid.UUID) (*StationWithLine, error) {
		return nil, ErrNotFound
	}

	req := &pb.GetStationRequest{Id: stationID.String()}
	resp, err := service.GetStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "station not found")
}

func TestUpdateStation_Success_NameAndStatus(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()
	lineID := uuid.New()
	now := time.Now()

	mockRepo.UpdateStationFn = func(ctx context.Context, id uuid.UUID, name, status *string) (*StationWithLine, error) {
		assert.Equal(t, stationID, id)
		assert.NotNil(t, name)
		assert.Equal(t, "Updated Station", *name)
		assert.NotNil(t, status)
		assert.Equal(t, "maintenance", *status)
		return &StationWithLine{
			ID:        stationID,
			Name:      *name,
			LineID:    lineID,
			LineName:  "Test Line",
			Status:    *status,
			CreatedAt: now,
		}, nil
	}

	req := &pb.UpdateStationRequest{
		Id:     stationID.String(),
		Name:   "Updated Station",
		Status: "maintenance",
	}
	resp, err := service.UpdateStation(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, stationID.String(), resp.Id)
	assert.Equal(t, "Updated Station", resp.Name)
	assert.Equal(t, "maintenance", resp.Status)
}

func TestUpdateStation_OnlyName(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()
	lineID := uuid.New()
	now := time.Now()

	mockRepo.UpdateStationFn = func(ctx context.Context, id uuid.UUID, name, status *string) (*StationWithLine, error) {
		assert.NotNil(t, name)
		assert.Nil(t, status)
		return &StationWithLine{
			ID:        stationID,
			Name:      *name,
			LineID:    lineID,
			LineName:  "Test Line",
			Status:    "active",
			CreatedAt: now,
		}, nil
	}

	req := &pb.UpdateStationRequest{
		Id:   stationID.String(),
		Name: "Updated Station",
	}
	resp, err := service.UpdateStation(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "Updated Station", resp.Name)
}

func TestUpdateStation_OnlyStatus(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()
	lineID := uuid.New()
	now := time.Now()

	mockRepo.UpdateStationFn = func(ctx context.Context, id uuid.UUID, name, status *string) (*StationWithLine, error) {
		assert.Nil(t, name)
		assert.NotNil(t, status)
		return &StationWithLine{
			ID:        stationID,
			Name:      "Test Station",
			LineID:    lineID,
			LineName:  "Test Line",
			Status:    *status,
			CreatedAt: now,
		}, nil
	}

	req := &pb.UpdateStationRequest{
		Id:     stationID.String(),
		Status: "closed",
	}
	resp, err := service.UpdateStation(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "closed", resp.Status)
}

func TestUpdateStation_InvalidUUID(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.UpdateStationRequest{
		Id:   "invalid-uuid",
		Name: "Updated Station",
	}
	resp, err := service.UpdateStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUpdateStation_NameTooLong(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	longName := string(make([]byte, 101))
	req := &pb.UpdateStationRequest{
		Id:   uuid.New().String(),
		Name: longName,
	}
	resp, err := service.UpdateStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUpdateStation_InvalidStatus(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.UpdateStationRequest{
		Id:     uuid.New().String(),
		Status: "invalid-status",
	}
	resp, err := service.UpdateStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUpdateStation_NoFieldsProvided(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.UpdateStationRequest{
		Id: uuid.New().String(),
	}
	resp, err := service.UpdateStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "at least one field")
}

func TestUpdateStation_NotFound(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()

	mockRepo.UpdateStationFn = func(ctx context.Context, id uuid.UUID, name, status *string) (*StationWithLine, error) {
		return nil, ErrNotFound
	}

	req := &pb.UpdateStationRequest{
		Id:   stationID.String(),
		Name: "Updated Station",
	}
	resp, err := service.UpdateStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestUpdateStation_RepositoryError(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()

	mockRepo.UpdateStationFn = func(ctx context.Context, id uuid.UUID, name, status *string) (*StationWithLine, error) {
		return nil, errors.New("database error")
	}

	req := &pb.UpdateStationRequest{
		Id:   stationID.String(),
		Name: "Updated Station",
	}
	resp, err := service.UpdateStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestDeleteStation_Success(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()

	mockRepo.DeleteStationFn = func(ctx context.Context, id uuid.UUID) error {
		assert.Equal(t, stationID, id)
		return nil
	}

	req := &pb.DeleteStationRequest{Id: stationID.String()}
	resp, err := service.DeleteStation(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDeleteStation_InvalidUUID(t *testing.T) {
	service, _ := setupServiceWithMock()
	ctx := context.Background()

	req := &pb.DeleteStationRequest{Id: "invalid-uuid"}
	resp, err := service.DeleteStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestDeleteStation_NotFound(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()

	mockRepo.DeleteStationFn = func(ctx context.Context, id uuid.UUID) error {
		return ErrNotFound
	}

	req := &pb.DeleteStationRequest{Id: stationID.String()}
	resp, err := service.DeleteStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestDeleteStation_RepositoryError(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()

	mockRepo.DeleteStationFn = func(ctx context.Context, id uuid.UUID) error {
		return errors.New("database error")
	}

	req := &pb.DeleteStationRequest{Id: stationID.String()}
	resp, err := service.DeleteStation(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestCreateStation_AllValidStatuses(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	validStatuses := []string{"active", "inactive", "maintenance", "closed"}

	for _, validStatus := range validStatuses {
		t.Run(validStatus, func(t *testing.T) {
			stationID := uuid.New()
			lineID := uuid.New()
			now := time.Now()

			mockRepo.CreateStationFn = func(ctx context.Context, name string, lID uuid.UUID, status string) (*StationWithLine, error) {
				assert.Equal(t, validStatus, status)
				return &StationWithLine{
					ID:        stationID,
					Name:      name,
					LineID:    lID,
					LineName:  "Test Line",
					Status:    status,
					CreatedAt: now,
				}, nil
			}

			req := &pb.CreateStationRequest{
				Name:   "Test Station",
				LineId: lineID.String(),
				Status: validStatus,
			}
			resp, err := service.CreateStation(ctx, req)

			require.NoError(t, err)
			assert.Equal(t, validStatus, resp.Status)
		})
	}
}

func TestUpdateStation_AllValidStatuses(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	validStatuses := []string{"active", "inactive", "maintenance", "closed"}

	for _, validStatus := range validStatuses {
		t.Run(validStatus, func(t *testing.T) {
			stationID := uuid.New()
			lineID := uuid.New()
			now := time.Now()

			mockRepo.UpdateStationFn = func(ctx context.Context, id uuid.UUID, name, status *string) (*StationWithLine, error) {
				assert.NotNil(t, status)
				assert.Equal(t, validStatus, *status)
				return &StationWithLine{
					ID:        stationID,
					Name:      "Test Station",
					LineID:    lineID,
					LineName:  "Test Line",
					Status:    *status,
					CreatedAt: now,
				}, nil
			}

			req := &pb.UpdateStationRequest{
				Id:     stationID.String(),
				Status: validStatus,
			}
			resp, err := service.UpdateStation(ctx, req)

			require.NoError(t, err)
			assert.Equal(t, validStatus, resp.Status)
		})
	}
}
func TestCreateLine_Idempotent(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()
	now := time.Now()
	lineName := "Test Line"

	mockRepo.CreateLineFn = func(ctx context.Context, name string) (*Line, error) {
		return &Line{
			ID:        lineID,
			Name:      lineName,
			CreatedAt: now,
		}, nil
	}

	req := &pb.CreateLineRequest{
		Name: lineName,
	}

	resp1, err1 := service.CreateLine(ctx, req)
	require.NoError(t, err1)
	assert.Equal(t, lineID.String(), resp1.Id)
	assert.Equal(t, lineName, resp1.Name)

	resp2, err2 := service.CreateLine(ctx, req)
	require.NoError(t, err2)
	assert.Equal(t, lineID.String(), resp2.Id)
	assert.Equal(t, lineName, resp2.Name)

	assert.Equal(t, resp1.Id, resp2.Id)
	assert.Equal(t, resp1.Name, resp2.Name)
}

func TestCreateStation_Idempotent(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()
	lineID := uuid.New()
	now := time.Now()
	stationName := "Test Station"

	mockRepo.CreateStationFn = func(ctx context.Context, name string, lID uuid.UUID, status string) (*StationWithLine, error) {
		return &StationWithLine{
			ID:        stationID,
			Name:      stationName,
			LineID:    lineID,
			LineName:  "Test Line",
			Status:    "active",
			CreatedAt: now,
		}, nil
	}

	req := &pb.CreateStationRequest{
		Name:   stationName,
		LineId: lineID.String(),
		Status: "active",
	}

	resp1, err1 := service.CreateStation(ctx, req)
	require.NoError(t, err1)
	assert.Equal(t, stationID.String(), resp1.Id)
	assert.Equal(t, stationName, resp1.Name)

	resp2, err2 := service.CreateStation(ctx, req)
	require.NoError(t, err2)
	assert.Equal(t, stationID.String(), resp2.Id)
	assert.Equal(t, stationName, resp2.Name)

	assert.Equal(t, resp1.Id, resp2.Id)
	assert.Equal(t, resp1.Name, resp2.Name)
	assert.Equal(t, resp1.Status, resp2.Status)
}

func TestUpdateLine_Idempotent(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	lineID := uuid.New()
	now := time.Now()
	newName := "Updated Line"

	mockRepo.UpdateLineFn = func(ctx context.Context, id uuid.UUID, name string) (*Line, error) {
		return &Line{
			ID:        lineID,
			Name:      newName,
			CreatedAt: now,
		}, nil
	}

	req := &pb.UpdateLineRequest{
		Id:   lineID.String(),
		Name: newName,
	}

	resp1, err1 := service.UpdateLine(ctx, req)
	require.NoError(t, err1)
	assert.Equal(t, lineID.String(), resp1.Id)
	assert.Equal(t, newName, resp1.Name)

	resp2, err2 := service.UpdateLine(ctx, req)
	require.NoError(t, err2)
	assert.Equal(t, lineID.String(), resp2.Id)
	assert.Equal(t, newName, resp2.Name)

	assert.Equal(t, resp1.Id, resp2.Id)
	assert.Equal(t, resp1.Name, resp2.Name)
}

func TestUpdateStation_Idempotent(t *testing.T) {
	service, mockRepo := setupServiceWithMock()
	ctx := context.Background()

	stationID := uuid.New()
	lineID := uuid.New()
	now := time.Now()
	newName := "Updated Station"
	newStatus := "maintenance"

	mockRepo.UpdateStationFn = func(ctx context.Context, id uuid.UUID, name, status *string) (*StationWithLine, error) {
		return &StationWithLine{
			ID:        stationID,
			Name:      newName,
			LineID:    lineID,
			LineName:  "Test Line",
			Status:    newStatus,
			CreatedAt: now,
		}, nil
	}

	req := &pb.UpdateStationRequest{
		Id:     stationID.String(),
		Name:   newName,
		Status: newStatus,
	}

	resp1, err1 := service.UpdateStation(ctx, req)
	require.NoError(t, err1)
	assert.Equal(t, stationID.String(), resp1.Id)
	assert.Equal(t, newName, resp1.Name)
	assert.Equal(t, newStatus, resp1.Status)

	resp2, err2 := service.UpdateStation(ctx, req)
	require.NoError(t, err2)
	assert.Equal(t, stationID.String(), resp2.Id)
	assert.Equal(t, newName, resp2.Name)
	assert.Equal(t, newStatus, resp2.Status)

	assert.Equal(t, resp1.Id, resp2.Id)
	assert.Equal(t, resp1.Name, resp2.Name)
	assert.Equal(t, resp1.Status, resp2.Status)
}
