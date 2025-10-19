package main

import (
	"context"
	"mime"
	"net/http"

	"github.com/bluesg/transport-analytics/backend"
	"github.com/bluesg/transport-analytics/config"
	myapp "github.com/bluesg/transport-analytics/proto"
	"github.com/bluesg/transport-analytics/version"
	"github.com/go-coldbrew/core"
	"github.com/go-coldbrew/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/proto"

	openapi "github.com/bluesg/transport-analytics/third_party/OpenAPI"
)

type cbSvc struct {
	stopper      core.CBStopper
	db           *sqlx.DB
	transportSvc *backend.Service
}

func (s *cbSvc) FailCheck(fail bool) {
}

func (s *cbSvc) Stop() {
	if s.db != nil {
		s.db.Close()
	}
	if s.stopper != nil {
		s.stopper.Stop()
	}
}

func (s *cbSvc) InitHTTP(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	err := myapp.RegisterTransportAnalyticsHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		return err
	}
	// Wrap mux to handle OPTIONS requests
	return nil
}

func (s *cbSvc) HTTPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers for ALL requests
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle CORS preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (s *cbSvc) GRPCGatewayMuxOptions() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithForwardResponseOption(addCORSHeaders),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
			// Add CORS headers even for errors
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
		}),
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			md := metadata.MD{}
			return md
		}),
	}
}

func addCORSHeaders(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	// Set CORS headers for all successful responses
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	return nil
}

func (s *cbSvc) InitGRPC(ctx context.Context, server *grpc.Server) error {
	cfg := config.Get()

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Error(ctx, "Failed to connect to database", "error", err)
		return err
	}
	s.db = db

	s.db.SetMaxOpenConns(25)
	s.db.SetMaxIdleConns(5)

	log.Info(ctx, "Database connection established")

	repo := backend.NewRepository(db)
	s.transportSvc = backend.NewService(repo)

	myapp.RegisterTransportAnalyticsServer(server, s.transportSvc)

	healthgrpc.RegisterHealthServer(server, &healthService{})

	log.Info(ctx, "Transport analytics assessment registered")

	return nil
}

type healthService struct {
	healthgrpc.UnimplementedHealthServer
}

func (h *healthService) Check(ctx context.Context, req *healthgrpc.HealthCheckRequest) (*healthgrpc.HealthCheckResponse, error) {
	return &healthgrpc.HealthCheckResponse{
		Status: healthgrpc.HealthCheckResponse_SERVING,
	}, nil
}

func getOpenAPIHandler() http.Handler {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		log.Error(context.Background(), "msg", "error adding mime type", "err", err)
	}
	return http.FileServer(http.FS(openapi.ContentFS))
}

func main() {
	cfg := config.GetColdBrewConfig()
	if cfg.AppName == "" {
		cfg.AppName = "backend-analytics"
	}
	cfg.ReleaseName = version.GitCommit

	cb := core.New(cfg)
	cb.SetOpenAPIHandler(getOpenAPIHandler())

	err := cb.SetService(&cbSvc{})
	if err != nil {
		panic(err)
	}

	log.Error(context.Background(), cb.Run())
}
