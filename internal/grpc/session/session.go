package session

import (
	"context"
	"time"
	"uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	parkingv1 "parking/gen/go/proto/parking"
	"parking/internal/domain"
	grpcsvc "parking/internal/grpc"
)

// Service defines session business operations.
type Service interface {
	// OpenSession creates a parking session when a car enters the session.
	OpenSession(
		ctx context.Context,
		stateNumber string,
		openedAt time.Time,
		parkingID uuid.UUID,
	) (*domain.Session, error)
	// CloseSession finishes a parking session when a car exits the session.
	CloseSession(
		ctx context.Context,
		stateNumber string,
		closedAt time.Time,
		parkingID uuid.UUID,
	) (*domain.Session, error)
	// CalculatePrice calculates a session price.
	CalculatePrice(ctx context.Context, stateNumber string, parkingID uuid.UUID) (int64, error)
	// GetSessionsByCarNumber finds parking sessions by the car state number.
	GetSessionsByCarNumber(
		ctx context.Context,
		stateNumber string,
		parkingID uuid.UUID,
	) ([]domain.Session, error)
	// GetActiveSessionsByParking finds active parking sessions by the session ID.
	GetActiveSessionsByParking(ctx context.Context, parkingID uuid.UUID) ([]domain.Session, error)
}

type serverAPI struct {
	parkingv1.UnimplementedParkingServer

	service Service
}

// Register registers the Parking service implementation with the gRPC server.
func Register(gRPC *grpc.Server, service Service) {
	parkingv1.RegisterParkingServer(gRPC, &serverAPI{
		service: service,
	})
}

// TODO: add validation for requests

// Enter registers entry to the session.
func (s *serverAPI) Enter(
	ctx context.Context,
	req *parkingv1.EnterRequest,
) (*parkingv1.ParkingSession, error) {
	parkingID, err := uuid.Parse(req.GetParkingId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	session, err := s.service.OpenSession(ctx, req.GetStateNumber(), req.GetEnteredAt().AsTime(), parkingID)
	if err != nil {
		return nil, grpcsvc.MapError(err)
	}

	return toProtoSession(session), nil
}

// Exit registers exits from the session.
func (s *serverAPI) Exit(
	ctx context.Context,
	req *parkingv1.ExitRequest,
) (*parkingv1.ParkingSession, error) {
	parkingID, err := uuid.Parse(req.GetParkingId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	session, err := s.service.CloseSession(ctx, req.GetStateNumber(), req.GetExitedAt().AsTime(), parkingID)
	if err != nil {
		return nil, grpcsvc.MapError(err)
	}

	return toProtoSession(session), nil
}

// CalculatePrice calculates the session fee.
func (s *serverAPI) CalculatePrice(
	ctx context.Context,
	req *parkingv1.CalculatePriceRequest,
) (*parkingv1.CalculatePriceResponse, error) {
	parkingID, err := uuid.Parse(req.GetParkingId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	cost, err := s.service.CalculatePrice(ctx, req.GetStateNumber(), parkingID)
	if err != nil {
		return nil, grpcsvc.MapError(err)
	}

	return &parkingv1.CalculatePriceResponse{
		Cost: cost,
	}, nil
}

// GetCarSessions finds and returns sessions by the car.
func (s *serverAPI) GetCarSessions(
	req *parkingv1.GetCarSessionsRequest,
	stream grpc.ServerStreamingServer[parkingv1.ParkingSession],
) error {
	parkingID, err := uuid.Parse(req.GetParkingId())
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	sessions, err := s.service.GetSessionsByCarNumber(stream.Context(), req.GetStateNumber(), parkingID)
	if err != nil {
		return grpcsvc.MapError(err)
	}

	s.sendSessions(stream, sessions)

	return nil
}

// GetActiveSessions finds and returns active sessions by the session.
func (s *serverAPI) GetActiveSessions(
	req *parkingv1.GetActiveSessionsRequest,
	stream grpc.ServerStreamingServer[parkingv1.ParkingSession],
) error {
	parkingID, err := uuid.Parse(req.GetParkingId())
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	sessions, err := s.service.GetActiveSessionsByParking(stream.Context(), parkingID)
	if err != nil {
		return grpcsvc.MapError(err)
	}

	s.sendSessions(stream, sessions)

	return nil
}

func (s *serverAPI) sendSessions(
	stream grpc.ServerStreamingServer[parkingv1.ParkingSession],
	sessions []domain.Session,
) error {
	for _, session := range sessions {
		if err := stream.Send(toProtoSession(&session)); err != nil {
			return err
		}
	}

	return nil
}
